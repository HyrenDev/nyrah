package implementations

import (
	"errors"
	"fmt"
	"net/hyren/nyrah/applications/status"
	"net/hyren/nyrah/cache/local"
	"net/hyren/nyrah/minecraft"
	"net/hyren/nyrah/misc/io"
	"net/hyren/nyrah/misc/providers"
	"time"

	ProxyStatus "net/hyren/nyrah/applications/status/implementations"
	PacketHandler "net/hyren/nyrah/misc/packets"
)

func CreateServer(address string, port int) {
	server := minecraft.NewServer(address, port, PacketHandler.HandlePackets)

	if server == nil {
		fmt.Println("Failed to create minecraft server")

		return
	}

	err := server.ListenAndServe()

	if err != nil {
		fmt.Println(err)

		return
	}

	fmt.Println("Started minecraft server on ", address, " with port ", port)
}

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