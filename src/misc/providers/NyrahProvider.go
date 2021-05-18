package providers

import (
	"errors"
	MariaDB "net/hyren/nyrah/providers/databases/mariadb"
	Redis "net/hyren/nyrah/providers/databases/redis"
)

var (
	primaryProvidersPrepared = false

	MARIA_DB_MAIN = new(MariaDB.MariaDBDatabaseProvider)
	REDIS_MAIN = new(Redis.RedisDatabaseProvider)
)

func PreparePrimaryProviders() error {
	if primaryProvidersPrepared == true {
		return errors.New("Application already prepared")
	}

	MARIA_DB_MAIN.Prepare()
	REDIS_MAIN.Prepare()

	primaryProvidersPrepared = true

	return nil
}