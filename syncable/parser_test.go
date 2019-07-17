package syncable

import (
	"bytes"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/philborlin/committed/types"
	"github.com/spf13/viper"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Syncable Parser", func() {
	Describe("Simple toml", func() {
		var (
			data []byte
			err  error
			dbs  map[string]types.Database
		)

		JustBeforeEach(func() {
			data, err = ioutil.ReadFile("./simple.toml")
			Expect(err).To(BeNil())
			dbs, err = databases()
			Expect(err).To(BeNil())
		})

		It("should parse with SQL toml", func() {
			name, parsed, err := Parse("toml", bytes.NewReader(data), dbs)
			Expect(err).To(BeNil())
			Expect(name).To(Equal("foo"))

			// wrapper := parsed.(*syncableWrapper)
			// sqlSyncable := wrapper.Syncable.(*SQLSyncable)
			syncable := parsed.(*SQLSyncable)

			actual := syncable.config
			expected := simpleConfig()

			Expect(actual).To(Equal(expected))
		})
	})
})

// func TestParseWithSQLToml(t *testing.T) {
// 	dat, err := ioutil.ReadFile("./simple.toml")
// 	if err != nil {
// 		t.Fatalf("Failed with error %v", err)
// 	}

// 	_, parsed, err := Parse("toml", bytes.NewReader(dat), databases())
// 	if err != nil {

// 	}
// 	fmt.Printf("%T\n", parsed)
// 	wrapper := parsed.(*syncableWrapper)
// 	sqlSyncable := wrapper.Syncable.(*SQLSyncable)

// 	actual := sqlSyncable.config
// 	expected := simpleConfig()

// 	if !reflect.DeepEqual(actual, expected) {
// 		t.Fatalf("Expected %v but was %v", expected, actual)
// 	}
// }

func TestSQLParser(t *testing.T) {
	dat, err := ioutil.ReadFile("./simple.toml")
	if err != nil {
		t.Fatalf("Failed with error %v", err)
	}

	var v = viper.New()

	v.SetConfigType("toml")
	v.ReadConfig(bytes.NewBuffer(dat))

	dbs, err := databases()
	if err != nil {
		t.Fatalf("Failed with error %v", err)
	}
	topicSyncable, err := sqlParser(v, dbs)
	if err != nil {
		t.Fatalf("Failed with error %v", err)
	}

	actual := topicSyncable.(*SQLSyncable).config

	expected := simpleConfig()

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Expected %v but was %v", expected, actual)
	}
}

func simpleConfig() *sqlConfig {
	m1 := sqlMapping{jsonPath: "$.Key", column: "pk", sqlType: "TEXT"}
	m2 := sqlMapping{jsonPath: "$.One", column: "one", sqlType: "TEXT"}
	m := []sqlMapping{m1, m2}

	i1 := index{indexName: "firstIndex", columnNames: "one"}
	i := []index{i1}
	return &sqlConfig{sqlDB: "testdb", topic: "test1", table: "foo", mappings: m, indexes: i, primaryKey: "pk"}
}

func databases() (map[string]types.Database, error) {
	sqlDB := types.NewSQLDB("ramsql", "memory://foo")
	err := sqlDB.Init()
	if err != nil {
		return nil, err
	}
	m := make(map[string]types.Database)
	m["testdb"] = sqlDB
	return m, nil
}
