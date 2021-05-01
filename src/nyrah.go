package main

import (
	Proxy "./misc/proxy"
	Env "./misc/utils"

	"log"
)

func main() {
	log.Println("Starting proxy server")

	Proxy.CreateServer(
		Env.GetServerAddress(),
		Env.GetServerPort(),
	)
}
