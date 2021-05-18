package implementations

import (
	"errors"
	"fmt"
	"net/hyren/nyrah/applications/status"
	"net/hyren/nyrah/minecraft"
	PacketHandler "net/hyren/nyrah/misc/packets"

	ProxyStatus "net/hyren/nyrah/applications/status/implementations"
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