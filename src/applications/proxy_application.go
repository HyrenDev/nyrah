package applications

import (
	"errors"
	"fmt"
	"log"
	"math/rand"

	Databases "../databases"
	ProxyStatus "./status"
)

func GetProxyAddress(key string) string {
	db := Databases.StartPostgres()

	row, err := db.Query("SELECT \"address\", \"port\" FROM \"apps\" WHERE \"name\"='" + key + "'")

	if err != nil {
		panic(err)
	}

	var address string
	var port int

	for row.Next() {
		row.Scan(&address, &port)
	}

	row.Close()
	db.Close()

	return fmt.Sprintf("%s:%d", address, port)
}

func GetRandomProxy(proxies []string) (string, error) {
	rand.Shuffle(len(proxies), func(i, j int) {
		proxies[i], proxies[j] = proxies[j], proxies[i]
	})

	for _, proxy := range proxies {
		log.Println("Getting status from ", proxy)

		var address = GetProxyAddress(proxy)

		var online = ProxyStatus.IsProxyOnline(address)

		if online {
			return proxy, nil
		}
	}

	return "", errors.New("couldn't find an proxy online")
}
