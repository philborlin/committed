package sql

import (
	"context"
	"database/sql"
	gosql "database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/oliveagle/jsonpath"
	"github.com/philborlin/committed/internal/node/syncable"
	"github.com/philborlin/committed/internal/node/types"
	"github.com/spf13/viper"
)

type parser struct{}

func (p *parser) Parse(v *viper.Viper, dbs map[string]syncable.Database) (syncable.Syncable, error) {
	return sqlParser(v, dbs)
}

func init() {
	syncable.RegisterParser("sql", &parser{})
}

type index struct {
	indexName   string
	columnNames string // comma separated list of columns
}

type sqlMapping struct {
	jsonPath string
	column   string
	sqlType  string
	// TODO Add a concept of an optional mapping that doesn't error if it is missing
}

type sqlConfig struct {
	sqlDB      string
	topic      string
	table      string
	mappings   []sqlMapping
	indexes    []index
	primaryKey string // comma separated list of columns
}

// Syncable struct
type Syncable struct {
	config   *sqlConfig
	insert   *sqlInsert
	database *DB
	DB       *gosql.DB
}

type sqlInsert struct {
	stmt     *sql.Stmt
	jsonPath []string
}

func sqlParser(v *viper.Viper, databases map[string]syncable.Database) (syncable.Syncable, error) {
	topic := v.GetString("sql.topic")
	sqlDB := v.GetString("sql.db")
	table := v.GetString("sql.table")
	primaryKey := v.GetString("sql.primaryKey")

	var mappings []sqlMapping
	for _, item := range v.Get("sql.mappings").([]interface{}) {
		m := item.(map[string]interface{})
		mapping := sqlMapping{
			jsonPath: m["jsonPath"].(string),
			column:   m["column"].(string),
			sqlType:  m["type"].(string),
		}
		mappings = append(mappings, mapping)
	}

	var indexes []index
	for _, item := range v.Get("sql.indexes").([]interface{}) {
		m := item.(map[string]interface{})
		i := index{
			indexName:   m["name"].(string),
			columnNames: m["index"].(string),
		}
		indexes = append(indexes, i)
	}

	config := &sqlConfig{
		sqlDB:      sqlDB,
		topic:      topic,
		table:      table,
		mappings:   mappings,
		indexes:    indexes,
		primaryKey: primaryKey,
	}
	return newSyncable(config, databases)
}

// NewSyncable creates a new syncable
// TODO Move zero back into this package
func newSyncable(sqlConfig *sqlConfig, databases map[string]syncable.Database) (syncable.Syncable, error) {
	database := databases[sqlConfig.sqlDB]
	if database == nil {
		return &syncable.ZeroSyncable{}, fmt.Errorf("Database %s is not setup", sqlConfig.sqlDB)
	}
	if database.Type() != "sql" {
		return &syncable.ZeroSyncable{}, fmt.Errorf("Database %s is not a sql database", sqlConfig.sqlDB)
	}
	sqlDB := database.(*DB)

	return &Syncable{config: sqlConfig, database: sqlDB}, nil
}

// Sync syncs implements Syncable
func (s *Syncable) Sync(ctx context.Context, entry *types.AcceptedProposal) error {
	bytes := []byte(entry.Data)
	var jsonData interface{}
	json.Marshal(string(bytes))
	err := json.Unmarshal(bytes, &jsonData)
	if err != nil {
		log.Printf("Error Unmarshalling json: %v", err)
		return err
	}

	var values []interface{}
	for _, path := range s.insert.jsonPath {
		res, err := jsonpath.JsonPathLookup(jsonData, path)
		if err != nil {
			log.Printf("Error while parsing [%v] in [%v]: %v\n", path, jsonData, err)
			return err
		}
		values = append(values, res)
	}

	tx, err := s.DB.BeginTx(context.Background(), &sql.TxOptions{Isolation: 0, ReadOnly: false})

	if err != nil {
		log.Printf("Error while creating transaction: %v", err)
		return err
	}
	_, err = tx.Stmt(s.insert.stmt).ExecContext(ctx, values...)
	if err != nil {
		log.Printf("Error while executing statement: %v", err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Error while executing commit: %v", err)
		return err
	}

	return nil
}

func unwrapMappings(db *DB, config *sqlConfig) (*sqlInsert, error) {
	sql := db.dialect.CreateSQL(config.table, config.mappings)

	stmt, err := db.DB.Prepare(sql)
	if err != nil {
		log.Fatalf("Error Preparing sql [%s]: %v", sql, err)
	}

	var jsonPaths []string
	for _, mapping := range config.mappings {
		jsonPaths = append(jsonPaths, mapping.jsonPath)
	}

	return &sqlInsert{stmt, jsonPaths}, nil
}

// Init implements Syncable
func (s *Syncable) Init(ctx context.Context) error {
	return s.init(false)
}

func (s *Syncable) init(ignoreCreateDDLError bool) error {
	if err := s.database.Init(); err != nil {
		return err
	}
	s.DB = s.database.DB

	_, err := s.DB.Exec(s.database.dialect.CreateDDL(s.config))
	if err != nil && !ignoreCreateDDLError {
		return err
	}

	insert, err := unwrapMappings(s.database, s.config)
	if err != nil {
		return err
	}
	s.insert = insert

	return nil
}

// Close implements Syncable
func (s *Syncable) Close() error {
	return s.DB.Close()
}

// Topics implements Syncable
func (s *Syncable) Topics() []string {
	return []string{s.config.topic}
}
