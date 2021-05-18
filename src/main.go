package main

import (
	"fmt"
	"net/hyren/nyrah/minecraft"
	"os"

	PacketHandler "net/hyren/nyrah/misc/packets"
	NyrahProviders "net/hyren/nyrah/misc/providers"
	Config "net/hyren/nyrah/misc/utils"
)

func main() {
	err := NyrahProviders.PreparePrimaryProviders()

	if err != nil {
		fmt.Println(err)

		os.Exit(0)
	} else {
		fmt.Println("Starting proxy server")

		server := minecraft.NewServer(Config.GetServerAddress(), Config.GetServerPort(), PacketHandler.HandlePackets)

		if server == nil {
			fmt.Println("Failed to create minecraft server")

			return
		}

		err := server.ListenAndServe()

		if err != nil {
			fmt.Println(err)

			return
		}

		fmt.Println(
			"Started minecraft server on",
			Config.GetServerAddress(),
			" with port",
			Config.GetServerPort(),
		)
	}
}
