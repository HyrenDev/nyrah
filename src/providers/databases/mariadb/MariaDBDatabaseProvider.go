package mariadb

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/hyren/nyrah/environment"

	DatabaseProviders "net/hyren/nyrah/providers/databases"
)

var pool *sql.DB

type MariaDBDatabaseProvider struct {
	DatabaseProviders.IDatabaseProvider
}

func (databaseProvider MariaDBDatabaseProvider) Prepare() {
	var postgres = environment.Get("databases").(map[string]interface{})["maria_db"].(map[string]interface{})

	var host = postgres["host"].(string)
	var port = int(postgres["port"].(float64))
	var user = postgres["user"].(string)
	var password = postgres["password"].(string)
	var database = postgres["database"].(string)

	log.Printf("Connecting to MySQL database (%s:%d)...\n", host, port)

	var err error

	pool, err = sql.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s",
		user, password, host, port, database,
	))

	if err != nil {
		panic(err)
	}

	err = pool.Ping()

	if err != nil {
		panic(err)
	}

	pool.SetMaxOpenConns(10)
	pool.SetConnMaxIdleTime(5000)

	log.Println("MySQL connection established successfully!")
}

func (databaseProvider MariaDBDatabaseProvider) Provide() *sql.DB {
	return pool
}
