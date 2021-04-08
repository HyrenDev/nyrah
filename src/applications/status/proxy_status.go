package status

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"net"

	Databases "../../databases"
)

func IsProxyOnline(server string) bool {
	_, err := net.Dial("tcp", server)

	if err != nil {
		return false
	}

	return true
}

func GetProxyPlayerCount(proxy string) (interface{}, error) {
	redisConnection := Databases.StartRedis().Get()

	var proxyApplicationStatus, err = redis.Values(
		redisConnection.Do("GET", fmt.Sprintf("applications:%s", proxy)),
	)

	if err != nil {
		return nil, err
	}

	return proxyApplicationStatus, nil
}
