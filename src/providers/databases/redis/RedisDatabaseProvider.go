package redis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"net/hyren/nyrah/environment"
	DatabaseProviders "net/hyren/nyrah/providers/databases"
	"time"
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

	redisDatabaseProvider.pool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", host, port))

			if err != nil {
				log.Println(nil)
			}

			if _, err := c.Do("AUTH", password); err != nil {
				_ = c.Close()

				return nil, err
			}

			if _, err := c.Do("SELECT", 0); err != nil {
				_ = c.Close()

				return nil, err
			}

			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")

			return err
		},
	}
}

func (redisDatabaseProvider RedisDatabaseProvider) Provide() redis.Conn {
	redisConnection := redisDatabaseProvider.pool.Get()

	defer redisConnection.Close()

	return redisConnection
}