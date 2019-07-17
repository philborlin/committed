package syncable

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/oliveagle/jsonpath"
	"github.com/philborlin/committed/types"
	"github.com/spf13/viper"
)

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

// SQLSyncable struct
type SQLSyncable struct {
	config   *sqlConfig
	insert   *sqlInsert
	database *types.SQLDB
	DB       *sql.DB
}

type sqlInsert struct {
	stmt     *sql.Stmt
	jsonPath []string
}

func sqlParser(v *viper.Viper, databases map[string]types.Database) (TopicSyncable, error) {
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
	return newSQLSyncable(config, databases)
}

// NewSQLSyncable creates a new syncable
func newSQLSyncable(sqlConfig *sqlConfig, databases map[string]types.Database) (*SQLSyncable, error) {
	database := databases[sqlConfig.sqlDB]
	if database == nil {
		return nil, fmt.Errorf("Database %s is not setup", sqlConfig.sqlDB)
	}
	if database.Type() != "sql" {
		return nil, fmt.Errorf("Database %s is not a sql database", sqlConfig.sqlDB)
	}
	sqlDB := database.(*types.SQLDB)

	return &SQLSyncable{config: sqlConfig, database: sqlDB}, nil
}

// Sync syncs implements Syncable
func (s *SQLSyncable) Sync(ctx context.Context, bytes []byte) error {
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

func createDDL(config *sqlConfig) string {
	var ddl strings.Builder
	ddl.WriteString(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (", config.table))
	for i, column := range config.mappings {
		ddl.WriteString(fmt.Sprintf("%s %s", column.column, column.sqlType))
		if i < len(config.mappings)-1 {
			ddl.WriteString(",")
		}
	}
	if config.primaryKey != "" {
		ddl.WriteString(fmt.Sprintf(",PRIMARY KEY (%s)", config.primaryKey))
	}
	for _, index := range config.indexes {
		ddl.WriteString(fmt.Sprintf(",INDEX %s (%s)", index.indexName, index.columnNames))
	}
	ddl.WriteString(");")

	return ddl.String()
}

func unwrapMappings(db *sql.DB, table string, mappings []sqlMapping) (*sqlInsert, error) {
	sql := createSQL(table, mappings)

	stmt, err := db.Prepare(sql)
	if err != nil {
		log.Fatalf("Error Preparing sql [%s]: %v", sql, err)
	}

	var jsonPaths []string
	for _, mapping := range mappings {
		jsonPaths = append(jsonPaths, mapping.jsonPath)
	}

	return &sqlInsert{stmt, jsonPaths}, nil
}

func createSQL(table string, sqlMappings []sqlMapping) string {
	var sql strings.Builder

	fmt.Fprintf(&sql, "INSERT INTO %s(", table)
	for i, item := range sqlMappings {
		if i == 0 {
			fmt.Fprintf(&sql, "%s", item.column)
		} else {
			fmt.Fprintf(&sql, ",%s", item.column)
		}
	}
	fmt.Fprint(&sql, ") VALUES (")
	for i := range sqlMappings {
		if i == 0 {
			fmt.Fprintf(&sql, "$%d", i+1)
		} else {
			fmt.Fprintf(&sql, ",$%d", i+1)
		}
	}
	fmt.Fprint(&sql, ")")

	return sql.String()
}

// Init implements Syncable
func (s *SQLSyncable) Init() error {
	if err := s.database.Init(); err != nil {
		return err
	}
	s.DB = s.database.DB

	_, err := s.DB.Exec(createDDL(s.config))
	if err != nil {
		return err
	}

	insert, err := unwrapMappings(s.DB, s.config.table, s.config.mappings)
	if err != nil {
		return err
	}
	s.insert = insert

	return nil
}

// Close implements Syncable
func (s *SQLSyncable) Close() error {
	return s.DB.Close()
}

// topics implements TopicSyncable
func (s *SQLSyncable) topics() []string {
	return []string{s.config.topic}
}
