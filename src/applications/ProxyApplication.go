package applications

import (
	"errors"
	"fmt"
	"log"
	"net/hyren/nyrah/applications/status"
	ProxyStatus "net/hyren/nyrah/applications/status/proxy"
	"net/hyren/nyrah/cache/local"
	"net/hyren/nyrah/misc/io"
	"net/hyren/nyrah/misc/providers"
	"time"
)

func GetProxyAddress(key string) io.InetSocketAddress {
	inetSocketAddress, found := local.CACHE.Get(fmt.Sprintf("%s_inet_socket_address", key))

	if !found {
		log.Println("Fetching ip address from", key, "in database...")

		connection := providers.POSTGRESQL_MAIN.Provide()

		defer connection.Close()

		row, err := connection.Query(fmt.Sprintf(
			`SELECT "address", "port" FROM "applications" WHERE "name"='%s'`,
			key,
		))

		if err != nil {
			panic(err)
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

func GetRandomProxy(proxies ...string) (string, error) {
	var err error

	if len(proxies) == 0 {
		proxies, err = FetchAvailableProxiesNames()

		if err != nil {
			return "", err
		}
	}

	proxyApplicationName, err := GetBalancedProxyApplicationName(proxies)

	if proxyApplicationName == "" {
		return "", errors.New("Cannot find an proxy to send the player")
	}

	log.Println("Getting status from", proxyApplicationName)

	if err != nil {
		return "", err
	}

	return proxyApplicationName, nil
}

func FetchAvailableProxiesNames() ([]string, error) {
	availableProxiesNames, found := local.CACHE.Get("available_proxies_name")

	if !found {
		connection := providers.POSTGRESQL_MAIN.Provide()

		defer connection.Close()

		rows, err := connection.Query(fmt.Sprintf(
			`SELECT "name" FROM "applications" WHERE "application_type"='PROXY';`,
		))

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

func GetBalancedProxyApplicationName(proxies []string) (string, error) {
	var indexes = make([]int, 0)

	for index, proxy := range proxies {
		proxyAddress, err := ProxyStatus.GetApplicationAddress(proxy)

		if err == nil {
			online := ProxyStatus.IsProxyOnline(
				proxyAddress,
			)

			if online {
				indexes = append(indexes, index)
			}
		}
	}

	applicationsStatus := make([]status.ApplicationStatus, len(indexes))

	var totalPlayers = 0

	for i := 0; i < len(indexes); i++ {
		var name = proxies[indexes[i]]

		onlinePlayers, _ := ProxyStatus.GetApplicationOnlinePlayers(name)

		applicationsStatus[i] = status.ApplicationStatus{
			Name:          name,
			OnlinePlayers: onlinePlayers,
		}

		totalPlayers += onlinePlayers
	}

	if len(applicationsStatus) <= 0 {
		return "", errors.New("don't have online proxies")
	}

	var index = totalPlayers % len(applicationsStatus)

	return applicationsStatus[index].Name, nil
}