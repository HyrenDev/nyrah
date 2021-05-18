package applications

import (
	"errors"
	"fmt"
	"net/hyren/nyrah/cache/local"
	"net/hyren/nyrah/misc/providers"
	"time"

	ProxyStatus "net/hyren/nyrah/applications/implementations"
	"net/hyren/nyrah/misc/io"
)

func GetProxyAddress(key string) io.InetSocketAddress {
	inetSocketAddress, found := local.CACHE.Get(fmt.Sprintf("%s_inet_socket_address", key))

	if !found {
		fmt.Println("Fetching ip address from", key, "in database...")

		row, err := providers.MARIA_DB_MAIN.Provide().Query(fmt.Sprintf(
			"SELECT `address`, `port` FROM `applications` WHERE `name`='%s'",
			key,
		))

		if err != nil {
			fmt.Println(err)

			defer row.Close()
		} else {
			var address string
			var port int

			if row.Next() {
				row.Scan(&address, &port)

				defer row.Close()
			}

			inetSocketAddress = io.InetSocketAddress {
				Host: address,
				Port: port,
			}

			local.CACHE.Set(fmt.Sprintf("%s_inet_socket_address", key), inetSocketAddress, 5*time.Minute)
		}
	}

	return inetSocketAddress.(io.InetSocketAddress)
}

func FetchAvailableProxiesNames() ([]string, error) {
	availableProxiesNames, found := local.CACHE.Get("available_proxies_name")

	if !found {
		rows, err := providers.MARIA_DB_MAIN.Provide().Query("SELECT `name` FROM `applications` WHERE `application_type`='PROXY';")

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

		local.CACHE.Set("available_proxies_name", availableProxiesNames, 5*time.Minute)
	}

	return availableProxiesNames.([]string), nil
}

func GetRandomProxy(proxies []string) (string, error) {
	proxyApplicationName, err := ProxyStatus.GetBalancedProxyApplicationName(proxies)

	if proxyApplicationName == "" {
		return "", errors.New("Cannot find an proxy to send the player")
	}

	fmt.Println("Getting status from", proxyApplicationName)

	if err != nil {
		return "", err
	}

	return proxyApplicationName, nil
}
