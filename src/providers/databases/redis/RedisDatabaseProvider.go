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

	redisDatabaseProvider.pool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			fmt.Printf("Connecting to redis database (%s:%d)...\n", host, port)

			c, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", host, port))

			if err != nil {
				fmt.Println(err)
			}

			if _, err := c.Do("AUTH", password); err != nil {
				fmt.Println("Autho:", err)

				_ = c.Close()

				return nil, err
			}

			if _, err := c.Do("SELECT", 0); err != nil {
				fmt.Println("Select:", err)

				_ = c.Close()

				return nil, err
			}

			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")

			fmt.Println("Ping:", err)

			return err
		},
	}

	response, err := redisDatabaseProvider.pool.Get().Do("PING")

	if err != nil {
		fmt.Println("Test:", err)
	} else {
		fmt.Println("Response:", response)
	}
}

func (redisDatabaseProvider RedisDatabaseProvider) Provide() *redis.Pool {
	return redisDatabaseProvider.pool
}