package mariadb

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/hyren/nyrah/environment"

	DatabaseProviders "net/hyren/nyrah/providers/databases"
)

var pool = make(chan *sql.DB, 10)

type MariaDBDatabaseProvider struct {
	DatabaseProviders.IDatabaseProvider
}

func (databaseProvider MariaDBDatabaseProvider) Prepare() {
	log.Println("Preparing mysql connection...")

	var postgres = environment.Get("databases").(map[string]interface{})["maria_db"].(map[string]interface{})

	var host = postgres["host"].(string)
	var port = int(postgres["port"].(float64))
	var user = postgres["user"].(string)
	var password = postgres["password"].(string)
	var database = postgres["database"].(string)

	db, err := sql.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s",
		user, password, host, port, database,
	))

	if err != nil {
		panic(err)
	}

	err = db.Ping()

	if err != nil {
		panic(err)
	}

	log.Println("Connection established successfully!")

	pool <- db
}

func (databaseProvider MariaDBDatabaseProvider) Provide() *sql.DB {
	return <- pool
}
