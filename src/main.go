package main

import (
	Proxy "net/hyren/nyrah/misc/proxy"
	Config "net/hyren/nyrah/misc/utils"

	"log"
)

func main() {
	log.Println("Starting proxy server")

	Proxy.CreateServer(
		Config.GetServerAddress(),
		Config.GetServerPort(),
	)
}
