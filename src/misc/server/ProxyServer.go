package server

import (
	"log"

	Minecraft "net/hyren/nyrah/minecraft"
	PacketHandler "net/hyren/nyrah/misc/packets"
	Config "net/hyren/nyrah/misc/utils"
)

func StartServer() {
	log.Println("Starting proxy server...")

	server := Minecraft.NewServer(
		Config.GetServerAddress(),
		Config.GetServerPort(),
		PacketHandler.HandlePackets,
	)

	err := server.ListenAndServe()

	if err != nil {
		panic(err)
	}
}