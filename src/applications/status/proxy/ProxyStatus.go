package proxy

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"net"

	NyrahProvider "net/hyren/nyrah/misc/providers"
)

func GetApplicationOnlinePlayers(application string) (int, error) {
	var bytes, err = redis.Bytes(
		NyrahProvider.REDIS_MAIN.Provide().Do("GET", fmt.Sprintf("applications:%s", application)),
	)

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
	var bytes, err = redis.Bytes(
		NyrahProvider.REDIS_MAIN.Provide().Do("GET", fmt.Sprintf("applications:%s", application)),
	)

	if err != nil {
		return "", err
	}

	var data map[string]interface{}

	err = json.Unmarshal(bytes, &data)

	if err != nil {
		return "", err
	}

	address := data["inet_socket_address"].(map[string]interface{})

	return fmt.Sprintf(
		"%s:%d",
		address["host"].(string),
		int(address["port"].(float64)),
	), nil
}

func IsProxyOnline(server string) bool {
	_, err := net.Dial("tcp", server)

	if err != nil {
		return false
	}

	return true
}
