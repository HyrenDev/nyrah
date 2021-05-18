package redis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"net/hyren/nyrah/environment"
	"time"

	DatabaseProviders "net/hyren/nyrah/providers/databases"
)

type RedisDatabaseProvider struct {
	DatabaseProviders.IDatabaseProvider

	pool *redis.Pool
}

func (redisDatabaseProvider RedisDatabaseProvider) Prepare() {
	var main = environment.Get("databases").(map[string]interface{})["redis"].(map[string]interface{})["main"].(map[string]interface{})

	var host = main["host"].(string)
	var port = int(main["port"].(float64))
	var password = main["password"].(string)

	fmt.Println("Connect to", host, ":", port, " with", password)

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

			if _, err := connection.Do("SELECT", 0); err != nil {
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