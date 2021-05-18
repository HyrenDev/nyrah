package applications

import (
	"errors"
	"fmt"
	"net/hyren/nyrah/cache/local"
	"net/hyren/nyrah/misc/providers"
	"time"

	ProxyStatus "net/hyren/nyrah/applications/implementations"
)

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
