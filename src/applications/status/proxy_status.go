package status

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"sort"

	Databases "../../databases"
)

func GetBalancedProxyApplicationName(proxies []string) (string, error) {
	sort.Slice(proxies, func(index1, index2 int) bool {
		onlinePlayers1, _ := GetApplicationOnlinePlayers(proxies[index1])

		onlinePlayers2, _ := GetApplicationOnlinePlayers(proxies[index2])

		return onlinePlayers2 > onlinePlayers1
	})

	return proxies[0], nil
}

func GetApplicationOnlinePlayers(application string) (int, error) {
	redisConnection := Databases.StartRedis().Get()

	var onlinePlayers, err = redis.Int(
		redisConnection.Do("HGET", fmt.Sprintf("applications:%s", application), "onlinePlayers"),
	)

	if err != nil {
		return 0, err
	}

	return onlinePlayers, nil
}
