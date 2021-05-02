package proxy

import (
	"log"
	"net/hyren/nyrah/minecraft"

	PacketHandler "net/hyren/nyrah/misc/packets"
)

func CreateServer(address string, port int) {
	server := minecraft.NewServer(address, port, PacketHandler.HandlePackets)

	if server == nil {
		log.Println("Failed to create minecraft server")

		return
	}

	err := server.ListenAndServe()

	if err != nil {
		log.Println(err)

		return
	}

	log.Println("Started minecraft server on ", address, " with port ", port)
}
