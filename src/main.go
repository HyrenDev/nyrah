package main

import (
	NyrahProviders "net/hyren/nyrah/misc/providers"
	ProxyServer "net/hyren/nyrah/misc/server"
)

func main() {
	err := NyrahProviders.PreparePrimaryProviders()

	if err != nil {
		panic(err)
	}

	ProxyServer.StartServer()
}
