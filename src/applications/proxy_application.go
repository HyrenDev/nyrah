package applications

import (
	Databases "../databases"
	ProxyStatus "./status"
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
		row.Scan(&address, &port)
	}

	row.Close()
	db.Close()

	return fmt.Sprintf("%s:%d", address, port)
}

func GetRandomProxy(proxies []string) (string, error) {
	//rand.Shuffle(len(proxies), func(i, j int) {
	//	proxies[i], proxies[j] = proxies[j], proxies[i]
	//})
	//
	for _, proxy := range proxies {
		log.Println("Getting status from ", proxy)

		var address = GetProxyAddress(proxy)

		var online = ProxyStatus.IsProxyOnline(address)

		if online {
			var applicationsStatus, err = ProxyStatus.GetProxyPlayerCount(proxy)

			if err != nil {
				log.Println(err)
			} else {
				log.Println(applicationsStatus)

				return proxy, nil
			}
		}
	}

	return "", errors.New("couldn't find an proxy online")
}
