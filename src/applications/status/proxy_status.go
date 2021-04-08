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
			proxies = OriginalRemoveIndex(proxies, index)
		}
	}

	if len(proxies) > 1 {
		sort.Slice(proxies, func(index1 int, index2 int) bool {
			onlinePlayers1, _ := GetApplicationOnlinePlayers(proxies[index1])
			onlinePlayers2, _ := GetApplicationOnlinePlayers(proxies[index2])

			return onlinePlayers2 > onlinePlayers1
		})
	}

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

func OriginalRemoveIndex(arr []string, pos int) []string {
	newArray := make([]string, len(arr)-1)
	k := 0
	for i := 0; i < (len(arr) - 1); {
		if i != pos {
			newArray[i] = arr[k]
			k++
		} else {
			k++
		}
		i++
	}

	return newArray
}
