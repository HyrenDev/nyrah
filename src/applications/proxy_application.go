package applications

import (
	Databases "../databases"
	"./status"
	"errors"
	"fmt"
	"log"
)

func GetProxyAddress(key string) string {
	db := Databases.StartPostgres()

	row, err := db.Query("SELECT \"address\", \"port\" FROM \"applications\" WHERE \"name\"='" + key + "'")

	if err != nil {
		log.Println(err)
	}

	var address string
	var port int

	for row.Next() {
		_ = row.Scan(&address, &port)
	}

	_ = row.Close()
	_ = db.Close()

	return fmt.Sprintf("%s:%d", address, port)
}

func GetRandomProxy(proxies []string) (string, error) {
	proxyApplicationName, err := status.GetBalancedProxyApplicationName(proxies)

	log.Println("Getting status from ", proxyApplicationName)

	if err != nil {
		log.Println(err)

		return "", err
	}

	return proxyApplicationName, errors.New("couldn't find an proxy online")
}
