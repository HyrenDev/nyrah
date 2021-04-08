package status

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"net"
	"sort"

	Databases "../../databases"
)

func GetBalancedProxyApplicationName(proxies []string) (string, error) {
	log.Println(len(proxies))

	for index, proxy := range proxies {
		var online = IsProxyOnline(proxy)

		if !online {
			proxies = append(proxies[:index], proxies[index+1:]...)
		}
	}

	sort.Slice(proxies, func(index1 int, index2 int) bool {
		onlinePlayers1, _ := GetApplicationOnlinePlayers(proxies[index1])
		onlinePlayers2, _ := GetApplicationOnlinePlayers(proxies[index2])

		return onlinePlayers2 > onlinePlayers1
	})

	log.Println(len(proxies))

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

func IsProxyOnline(server string) bool {
	_, err := net.Dial("tcp", server)

	if err != nil {
		return false
	}

	return true
}
