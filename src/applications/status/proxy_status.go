package status

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"math/rand"
	"net"
	"sort"

	Databases "../../databases"
)

type ApplicationStatus struct {
	name          string
	onlinePlayers int
}

type ApplicationStatusSorter struct {
	applicationStatus []ApplicationStatus
}

func (s ApplicationStatusSorter) Len() int {
	return len(s.applicationStatus)
}

func (s ApplicationStatusSorter) Less(i, j int) bool {
	return s.applicationStatus[i].onlinePlayers < s.applicationStatus[j].onlinePlayers
}

func (s ApplicationStatusSorter) Swap(i, j int) {
	s.applicationStatus[i], s.applicationStatus[j] = s.applicationStatus[j], s.applicationStatus[i]
}

func sortApplicationStatus(applicationStatus []ApplicationStatus) {
	sorter := ApplicationStatusSorter{
		applicationStatus: applicationStatus,
	}

	sort.Sort(sorter)
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

	for i := 0; i < len(indexes); i++ {
		var name = proxies[indexes[i]]

		onlinePlayers, _ := GetApplicationOnlinePlayers(name)

		applicationsStatus[i] = ApplicationStatus{
			name:          name,
			onlinePlayers: onlinePlayers,
		}
	}

	rand.Shuffle(len(proxies), func(i, j int) {
		proxies[i], proxies[j] = proxies[j], proxies[i]
	})

	if len(applicationsStatus) > 1 {
		sortApplicationStatus(applicationsStatus)
	}

	return applicationsStatus[0].name, nil
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
