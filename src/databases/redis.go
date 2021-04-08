package databases

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"time"

	Env "../environment"
)

func StartRedis() *redis.Pool {
	var data = Env.ReadFile()

	var databases = data["databases"].(map[string]interface{})
	var main = databases["redis"].(map[string]interface{})["main"].(map[string]interface{})

	var host = main["host"].(string)
	var port = int(main["port"].(float64))
	var password = main["password"].(string)
	var database = 0

	redisServer := fmt.Sprintf("%s:%d", host, port)

	var pool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redisServer)

			if err != nil {
				panic(nil)
			}

			if _, err := c.Do("AUTH", password); err != nil {
				_ = c.Close()

				return nil, err
			}

			if _, err := c.Do("SELECT", database); err != nil {
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

	return pool
}
