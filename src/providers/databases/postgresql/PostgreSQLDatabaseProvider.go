package postgresql

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/hyren/nyrah/environment"
	DatabaseProviders "net/hyren/nyrah/providers/databases"
	"sync"
)

var (
	once     sync.Once
	host     string
	port     int
	user     string
	password string
	database string
	schema   string
)

type PostgreSQLDatabaseProvider struct {
	DatabaseProviders.IDatabaseProvider
}

func (databaseProvider PostgreSQLDatabaseProvider) Prepare() {
	var postgres = environment.Get("databases").(map[string]interface{})["postgresql"].(map[string]interface{})

	host = postgres["host"].(string)
	port = int(postgres["port"].(float64))
	user = postgres["user"].(string)
	password = postgres["password"].(string)
	database = postgres["database"].(string)
	schema = postgres["schema"].(string)
}

func (databaseProvider PostgreSQLDatabaseProvider) Provide() *sql.DB {
	connection, err := sql.Open("postgres", fmt.Sprintf(
		`host=%s port=%d user=%s password=%s dbname=%s sslmode=disable search_path=%s`,
		host, port, user, password, database, schema,
	))

	if err != nil {
		panic(err)
	}

	log.Println("PostgreSQL connection established successfully!")

	return connection
}
