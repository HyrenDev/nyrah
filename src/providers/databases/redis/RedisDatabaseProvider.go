package redis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"net/hyren/nyrah/environment"
	"time"

	DatabaseProviders "net/hyren/nyrah/providers/databases"
)

var pool *redis.Pool

type RedisDatabaseProvider struct {
	DatabaseProviders.IDatabaseProvider
}

func (redisDatabaseProvider RedisDatabaseProvider) Prepare() {
	var main = environment.Get("databases").(map[string]interface{})["redis"].(map[string]interface{})["main"].(map[string]interface{})

	var host = main["host"].(string)
	var port = int(main["port"].(float64))
	var password = main["password"].(string)

	log.Printf("Connecting to redis database (%s:%d)...\n", host, port)

	pool = &redis.Pool {
		Wait:        true,
		MaxIdle:     3,
		MaxActive:   8,
		IdleTimeout: 2000,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", host, port))

			if err != nil {
				panic(err)
			}

			if _, err = c.Do("AUTH", password); err != nil {
				return nil, err
			}

			if _, err = c.Do("SELECT", 0); err != nil {
				return nil, err
			}

			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")

			return err
		},
	}

	log.Println("Redis connection established successfully!")
}

func (redisDatabaseProvider RedisDatabaseProvider) Provide() *redis.Pool {
	return pool
}