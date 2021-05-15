package status

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"net"

	Databases "net/hyren/nyrah/databases"
)

type ApplicationStatus struct {
	name          string
	onlinePlayers int
}

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

	applicationsStatus := make([]ApplicationStatus, len(indexes))

	var totalPlayers = 0

	for i := 0; i < len(indexes); i++ {
		var name = proxies[indexes[i]]

		onlinePlayers, _ := GetApplicationOnlinePlayers(name)

		applicationsStatus[i] = ApplicationStatus{
			name:          name,
			onlinePlayers: onlinePlayers,
		}

		totalPlayers += onlinePlayers
	}

	if len(applicationsStatus) <= 0 {
		return "", errors.New("don't have online proxies")
	}

	var index = totalPlayers % len(applicationsStatus)

	return applicationsStatus[index].name, nil
}

func GetApplicationOnlinePlayers(application string) (int, error) {
	redisConnection := Databases.StartRedis().Get()

	var bytes, err = redis.Bytes(
		redisConnection.Do("GET", fmt.Sprintf("applications:%s", application)),
	)

	defer redisConnection.Close()

	if err != nil {
		return 0, err
	}

	var data map[string]interface{}

	err = json.Unmarshal(bytes, &data)

	if err != nil {
		return 0, err
	}
	return int(data["online_players"].(float64)), nil
}

func GetApplicationAddress(application string) (string, error) {
	redisConnection := Databases.StartRedis().Get()

	var bytes, err = redis.Bytes(
		redisConnection.Do("GET", fmt.Sprintf("applications:%s", application)),
	)

	defer redisConnection.Close()

	if err != nil {
		return "", err
	}

	var data map[string]interface{}

	err = json.Unmarshal(bytes, &data)

	if err != nil {
		return "", err
	}

	address := data["address"].(map[string]interface{})

	return fmt.Sprintf(
		"%s:%d",
		address["address"].(string),
		address["port"].(int),
	), nil
}

func IsProxyOnline(server string) bool {
	_, err := net.Dial("tcp", server)

	if err != nil {
		return false
	}

	return true
}
