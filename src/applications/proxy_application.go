package applications

import (
	"errors"
	"fmt"
	"github.com/patrickmn/go-cache"
	_ "github.com/patrickmn/go-cache"
	"log"
	"net/hyren/nyrah/applications/status"
	"time"

	Databases "net/hyren/nyrah/databases"
)

var (
	CACHE = cache.New(cache.NoExpiration, 10*time.Second)
)

type InetSocketAddress struct {
	host string
	port int
}

func (inetSocketAddress InetSocketAddress) GetHostAddress() string {
	return inetSocketAddress.host
}

func (inetSocketAddress InetSocketAddress) GetPort() int {
	return inetSocketAddress.port
}

func GetProxyAddress(key string) InetSocketAddress {
	inetSocketAddress, found := CACHE.Get(fmt.Sprintf("%s_inet_socket_address", key))

	if !found {
		fmt.Println("Fetching ip address from", key, "in database...")

		db := Databases.StartMariaDB()

		row, err := db.Query(fmt.Sprintf(
			"SELECT `address`, `port` FROM `applications` WHERE `name`='%s'",
			key,
		))

		defer db.Close()

		if err != nil {
			log.Println(err)

			defer row.Close()
		} else {
			var address string
			var port int

			if row.Next() {
				row.Scan(&address, &port)

				defer row.Close()
			}

			inetSocketAddress = InetSocketAddress {
				host: address,
				port: port,
			}

			CACHE.Set(fmt.Sprintf("%s_inet_socket_address", key), inetSocketAddress, 5*time.Minute)
		}
	}

	return inetSocketAddress.(InetSocketAddress)
}

func FetchAvailableProxiesNames() ([]string, error) {
	availableProxiesNames, found := CACHE.Get("available_proxies_name")

	if !found {
		db := Databases.StartMariaDB()

		rows, err := db.Query("SELECT `name` FROM `applications` WHERE `application_type`='PROXY';")

		defer db.Close()

		if err != nil {
			return make([]string, 0), err
		}

		var proxies []string

		for rows.Next() {
			var name string

			err := rows.Scan(&name)

			if err != nil {
				return make([]string, 0), err
			}

			proxies = append(proxies, name)
		}

		defer rows.Close()

		availableProxiesNames = proxies

		CACHE.Set("available_proxies_name", availableProxiesNames, 5*time.Minute)
	}

	return availableProxiesNames.([]string), nil
}

func GetRandomProxy(proxies []string) (string, error) {
	proxyApplicationName, err := status.GetBalancedProxyApplicationName(proxies)

	if proxyApplicationName == "" {
		return "", errors.New("Cannot find an proxy to send the player")
	}

	fmt.Println("Getting status from", proxyApplicationName)

	if err != nil {
		return "", err
	}

	return proxyApplicationName, nil
}
