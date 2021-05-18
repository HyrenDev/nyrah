package providers

import (
	"errors"
	PostgreSQL "net/hyren/nyrah/providers/databases/postgresql"
	Redis "net/hyren/nyrah/providers/databases/redis"
)

var (
	primaryProvidersPrepared = false

	POSTGRESQL_MAIN = new(PostgreSQL.PostgreSQLDatabaseProvider)
	REDIS_MAIN      = new(Redis.RedisDatabaseProvider)
)

func PreparePrimaryProviders() error {
	if primaryProvidersPrepared == true {
		return errors.New("Application already prepared")
	}

	POSTGRESQL_MAIN.Prepare()
	REDIS_MAIN.Prepare()

	primaryProvidersPrepared = true

	return nil
}