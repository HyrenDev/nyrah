package main

import (
	"fmt"
	"net/hyren/nyrah/misc/providers"
	"os"

	Proxy "net/hyren/nyrah/applications/implementations"
	Config "net/hyren/nyrah/misc/utils"
)

func main() {
	err := providers.PreparePrimaryProviders()

	if err != nil {
		fmt.Println(err)

		os.Exit(0)
	} else {
		fmt.Println("Starting proxy server")

		Proxy.CreateServer(
			Config.GetServerAddress(),
			Config.GetServerPort(),
		)
	}
}
