package databases

import (
	"fmt"

	"database/sql"
	_ "github.com/lib/pq"

	Env "../environment"
)

func StartPostgres() *sql.DB {
	var data = Env.ReadFile()

	var databases = data["databases"].(map[string]interface{})
	var global = databases["global"].(map[string]interface{})
	var postgres = global["postgres"].(map[string]interface{})

	var host = postgres["host"].(string)
	var port = int(postgres["port"].(float64))
	var user = postgres["user"].(string)
	var password = postgres["password"].(string)
	var database = postgres["database"].(string)
	var schema = postgres["schema"].(string)

	var infos = fmt.Sprintf(
	"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable search_path=%s",
		host,port, user, password, database, schema,
	)

	db, err := sql.Open("postgres", infos)

	if err != nil {
		panic(err)
	}

	err = db.Ping()

	if err != nil {
		panic(err)
	}

	return db
}
