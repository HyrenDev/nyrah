package redis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"net/hyren/nyrah/environment"
	"net/hyren/nyrah/providers/databases"
	"time"
)

type RedisDatabaseProvider struct {
	databases.IDatabaseProvider

	pool *redis.Pool
}

func (redisDatabaseProvider RedisDatabaseProvider) Prepare() {
	var _databases = environment.Get("databases").(map[string]interface{})
	var main = _databases["redis"].(map[string]interface{})["main"].(map[string]interface{})

	var host = main["host"].(string)
	var port = int(main["port"].(float64))
	var password = main["password"].(string)
	var database = 0

	redisServer := fmt.Sprintf("%s:%d", host, port)

	redisDatabaseProvider.pool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			connection, err := redis.Dial("tcp", redisServer)

			if err != nil || connection == nil {
				return connection, err
			}

			if _, err := connection.Do("AUTH", password); err != nil {
				_ = connection.Close()

				return nil, err
			}

			if _, err := connection.Do("SELECT", database); err != nil {
				_ = connection.Close()

				return nil, err
			}

			return connection, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")

			return err
		},
	}
}

func (redisDatabaseProvider RedisDatabaseProvider) Provide() redis.Conn {
	return redisDatabaseProvider.pool.Get()
}