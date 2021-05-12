package applications

import (
	"fmt"
	"log"
	"net/hyren/nyrah/applications/status"

	Databases "net/hyren/nyrah/databases"
)

func GetProxyAddress(key string) string {
	db := Databases.StartMariaDB()

	row, err := db.Query(fmt.Sprintf(
		"SELECT `address`, `port` FROM `applications` WHERE `name`='%s'",
		key,
	))

	defer db.Close()

	if err != nil {
		log.Println(err)
	}

	var address string
	var port int

	if row.Next() {
		row.Scan(&address, &port)
	}

	defer row.Close()

	return fmt.Sprintf("%s:%d", address, port)
}

func GetRandomProxy(proxies []string) (string, error) {
	proxyApplicationName, err := status.GetBalancedProxyApplicationName(proxies)

	log.Println("Getting status from", proxyApplicationName)

	if err != nil {
		return "", err
	}

	return proxyApplicationName, nil
}
