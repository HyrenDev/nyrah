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
	var indexes = make([]int, 0)

	for index, proxy := range proxies {
		proxyAddress, err := GetApplicationAddress(proxy)

		if err == nil {
			online := IsProxyOnline(
				proxyAddress,
			)

			if online {
				indexes = append(indexes, index)
			}
		}
	}

	newArray := make([]string, len(indexes))

	for i := 0; i < len(indexes); i++ {
		newArray[i] = proxies[indexes[i]]
	}

	if len(newArray) > 1 {
		sort.Slice(newArray, func(index1 int, index2 int) bool {
			onlinePlayers1, _ := GetApplicationOnlinePlayers(newArray[index1])
			onlinePlayers2, _ := GetApplicationOnlinePlayers(newArray[index2])

			return onlinePlayers2 > onlinePlayers1
		})
	}

	return newArray[0], nil
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
		log.Println("Falhou:", server)

		return false
	}

	log.Println("Address:", server)

	return true
}
