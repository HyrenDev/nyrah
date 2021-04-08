package status

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"math/big"
	"net"
	"sort"

	Databases "../../databases"
)

type ApplicationsStatus struct {
	applicationsStatus []ApplicationStatus
}

type ApplicationStatus struct {
	applicationName string
	applicationType string
	address         string
	onlineSince     big.Int
	heapSize        big.Int
	heapMaxSize     big.Int
	heapFreeSize    big.Int
	onlinePlayers   int
}

func (applicationsStatus *ApplicationsStatus) Append(applicationStatus ApplicationStatus) {
	applicationsStatus.applicationsStatus = append(applicationsStatus.applicationsStatus, applicationStatus)
}

func NewApplicationsStatus() ApplicationsStatus {
	return ApplicationsStatus{
		[]ApplicationStatus{},
	}
}

func (applicationsStatus *ApplicationsStatus) GetSortedApplicationsStatus() []ApplicationStatus {
	var _applicationsStatus = applicationsStatus.applicationsStatus

	sort.Slice(_applicationsStatus, func(index1, index2 int) bool {
		return _applicationsStatus[index1].onlinePlayers > _applicationsStatus[index2].onlinePlayers
	})

	return applicationsStatus.applicationsStatus
}

func GetBalancedProxyApplicationName(proxies []string) (string, error) {
	var newApplications = NewApplicationsStatus()

	for _, proxy := range proxies {
		applicationStatus, err := GetApplicationStatus(proxy)

		if err != nil {
			log.Println("asd")

			continue
		} else {
			log.Println("dale")

			newApplications.Append(applicationStatus)
		}
	}

	return newApplications.applicationsStatus[0].applicationName, nil
}

func IsProxyOnline(server string) bool {
	_, err := net.Dial("tcp", server)

	if err != nil {
		return false
	}

	return true
}

func GetApplicationStatus(application string) (ApplicationStatus, error) {
	redisConnection := Databases.StartRedis().Get()

	var serializedProxyApplicationStatus, err = redis.Bytes(
		redisConnection.Do("GET", fmt.Sprintf("applications:%s", application)),
	)

	if err != nil {
		return ApplicationStatus{}, err
	}

	var proxyApplicationStatus ApplicationStatus

	err = json.Unmarshal(serializedProxyApplicationStatus, &proxyApplicationStatus)

	if err != nil {
		return ApplicationStatus{}, err
	}

	log.Println(proxyApplicationStatus)

	return proxyApplicationStatus, nil
}
