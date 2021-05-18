package mariadb

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/hyren/nyrah/environment"

	DatabaseProviders "net/hyren/nyrah/providers/databases"
)

type MariaDBDatabaseProvider struct {
	DatabaseProviders.IDatabaseProvider

	connection *sql.DB
}

func (databaseProvider MariaDBDatabaseProvider) Prepare() {
	var postgres = environment.Get("databases").(map[string]interface{})["maria_db"].(map[string]interface{})

	var host = postgres["host"].(string)
	var port = int(postgres["port"].(float64))
	var user = postgres["user"].(string)
	var password = postgres["password"].(string)
	var database = postgres["database"].(string)

	var err error

	databaseProvider.connection, err = sql.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s",
		user, password, host, port, database,
	))

	if err != nil {
		fmt.Println(err)
	}

	err = databaseProvider.connection.Ping()

	if err != nil {
		fmt.Println(err)
	}
}

func (databaseProvider MariaDBDatabaseProvider) Provide() *sql.DB {
	return databaseProvider.connection
}
