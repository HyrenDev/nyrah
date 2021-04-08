package status

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"net"
	"sort"

	Databases "../../databases"
)

func GetBalancedProxyApplicationName(proxies []string) (string, error) {
	for index, proxy := range proxies {
		proxyAddress, err := GetApplicationAddress(proxy)

		if err == nil {
			online := IsProxyOnline(
				proxyAddress,
			)

			if !online {
				proxies = OriginalRemoveIndex(proxies, index)
			}
		} else {
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

	var bytes, err = redis.Bytes(
		redisConnection.Do("GET", fmt.Sprintf("applications:%s", application)),
	)

	if err != nil {
		return 0, err
	}

	var data map[string]interface{}

	err = json.Unmarshal(bytes, &data)

	if err != nil {
		return 0, err
	}
	return int(data["onlinePlayers"].(float64)), nil
}

func GetApplicationAddress(application string) (string, error) {
	redisConnection := Databases.StartRedis().Get()

	var bytes, err = redis.Bytes(
		redisConnection.Do("GET", fmt.Sprintf("applications:%s", application)),
	)

	if err != nil {
		return "", err
	}

	var data map[string]interface{}

	err = json.Unmarshal(bytes, &data)

	if err != nil {
		return "", err
	}

	return data["address"].(string), nil
}

func IsProxyOnline(server string) bool {
	_, err := net.Dial("tcp", server)

	if err != nil {
		return false
	}

	return true
}

func OriginalRemoveIndex(arr []string, pos int) []string {
	newArray := make([]string, len(arr))

	for i := 0; i < len(arr); i++ {
		if i != pos {
			newArray[i] = arr[i]
		} else {
			log.Println("É igual")
		}
	}

	return newArray
}
