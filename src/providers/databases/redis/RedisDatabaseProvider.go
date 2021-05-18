package redis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"net/hyren/nyrah/environment"
	DatabaseProviders "net/hyren/nyrah/providers/databases"
	"time"
)

type RedisDatabaseProvider struct {
	DatabaseProviders.IDatabaseProvider

	pool *redis.Pool
}

func (redisDatabaseProvider RedisDatabaseProvider) Prepare() {
	//
}

func (redisDatabaseProvider RedisDatabaseProvider) Provide() redis.Conn {
	var main = environment.Get("databases").(map[string]interface{})["redis"].(map[string]interface{})["main"].(map[string]interface{})

	var host = main["host"].(string)
	var port = int(main["port"].(float64))
	var password = main["password"].(string)

	var pool = &redis.Pool {
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			connection, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", host, port))

			if err != nil || connection == nil {
				return connection, err
			}

			if _, err := connection.Do("AUTH", password); err != nil {
				_ = connection.Close()

				return nil, err
			}

			if _, err := connection.Do("SELECT", "0"); err != nil {
				_ = connection.Close()

				return nil, err
			}

			return connection, err
		},
		TestOnBorrow: func(connection redis.Conn, time time.Time) error {
			_, err := connection.Do("PING")

			return err
		},
	}

	connection := pool.Get()

	defer connection.Close()

	return connection
}