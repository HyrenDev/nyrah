package main

import (
	"log"

	Proxy "./misc/proxy"
	Config "./misc/utils"
)

func main() {
	log.Println("Starting proxy server")

	Proxy.CreateServer(
		Config.GetServerAddress(),
		Config.GetServerPort(),
	)
}
