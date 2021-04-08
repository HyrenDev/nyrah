package proxy

import (
	"gominet"
	"log"

	PacketHandler "../packets"
)

func CreateServer(address string, port int) {
	server := gominet.NewServer(address, port, PacketHandler.HandlePackets)

	if server == nil {
		panic("Failed to create minecraft server")
	}

	err := server.ListenAndServe()

	if err != nil {
		panic(err)
	}

	log.Println("Started minecraft server on ", address, " with port ", port, ">")
}