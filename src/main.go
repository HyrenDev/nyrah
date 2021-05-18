package main

import (
	"log"
	"net/hyren/nyrah/minecraft"
	"os"

	PacketHandler "net/hyren/nyrah/misc/packets"
	NyrahProviders "net/hyren/nyrah/misc/providers"
	Config "net/hyren/nyrah/misc/utils"
)

func main() {
	err := NyrahProviders.PreparePrimaryProviders()

	if err != nil {
		log.Println(err)

		os.Exit(0)
	} else {
		log.Println("Starting proxy server...")

		server := minecraft.NewServer(Config.GetServerAddress(), Config.GetServerPort(), PacketHandler.HandlePackets)

		if server == nil {
			panic("Failed to create minecraft server")
		} else {
			err = server.ListenAndServe()

			if err != nil {
				panic(err)
			}

			log.Printf(
				"Started minecraft server on %s with port %d\n",
				Config.GetServerAddress(),
				Config.GetServerPort(),
			)
		}
	}
}
