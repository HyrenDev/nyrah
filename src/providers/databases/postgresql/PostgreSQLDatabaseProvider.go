package postgresql

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/hyren/nyrah/environment"
	"time"

	DatabaseProviders "net/hyren/nyrah/providers/databases"
)

var connection *sql.DB

type PostgreSQLDatabaseProvider struct {
	DatabaseProviders.IDatabaseProvider
}

func (databaseProvider PostgreSQLDatabaseProvider) Prepare() {
	var postgres = environment.Get("databases").(map[string]interface{})["postgresql"].(map[string]interface{})

	var host = postgres["host"].(string)
	var port = int(postgres["port"].(float64))
	var user = postgres["user"].(string)
	var password = postgres["password"].(string)
	var database = postgres["database"].(string)
	var schema = postgres["schema"].(string)

	log.Printf("Connecting to PostgreSQL database (%s:%d)...\n", host, port)

	var err error

	connection, err = sql.Open("postgres", fmt.Sprintf(
		`host=%s port=%d user=%s password=%s dbname=%s sslmode=disable search_path=%s`,
		host, port, user, password, database, schema,
	))

	if err != nil {
		panic(err)
	}

	err = connection.Ping()

	if err != nil {
		panic(err)
	}

	connection.SetMaxOpenConns(10)
	connection.SetConnMaxIdleTime(5000)
	connection.SetMaxIdleConns(0)
	connection.SetConnMaxLifetime(999999*time.Hour)

	log.Println("PostgreSQL connection established successfully!")
}

func (databaseProvider PostgreSQLDatabaseProvider) Provide() *sql.DB {
	err := connection.Ping()

	if err != nil {
		databaseProvider.Prepare()
	}

	return connection
}
