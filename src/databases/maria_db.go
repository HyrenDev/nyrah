package databases

import (
	"fmt"
	"log"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	Env "net/hyren/nyrah/environment"
)

func StartMariaDB() *sql.DB {
	var data = Env.ReadFile()

	var databases = data["databases"].(map[string]interface{})
	var postgres = databases["maria_db"].(map[string]interface{})

	var host = postgres["host"].(string)
	var port = int(postgres["port"].(float64))
	var user = postgres["user"].(string)
	var password = postgres["password"].(string)
	var database = postgres["database"].(string)

	var infos = fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s",
		user, password, host, port, database,
	)

	db, err := sql.Open("mysql", infos)

	if err != nil {
		log.Println(err)
	}

	err = db.Ping()

	if err != nil {
		log.Println(err)
	}

	return db
}
