package status

import (
	"fmt"
	"log"
	"net"

	redis "../../databases"
)

func IsProxyOnline(server string) bool {
	_, err := net.Dial("tcp", server)

	if err != nil {
		return false
	}

	return true
}

func GetProxyPlayerCount(proxy string) (int, error) {
	redisConnection := redis.StartRedis().Get()

	var proxyApplicationStatus, err = redisConnection.Do("GET", fmt.Sprintf("applications:%s", proxy))

	if err != nil {
		return 0, err
	}

	log.Println(proxyApplicationStatus)

	return 0, nil
}
