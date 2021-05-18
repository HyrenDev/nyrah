package connector

import (
	"fmt"
	"io"
	"net"
	"net/hyren/nyrah/applications"
	"net/hyren/nyrah/minecraft/protocol"
)

func ConnectToProxy(connection *protocol.Connection, proxy string) {
	var inetSocketAddress = applications.GetProxyAddress(proxy)

	ds, err := net.Dial("tcp", fmt.Sprintf("%s:%d", inetSocketAddress.Host, inetSocketAddress.Port))

	if err != nil {
		connection.Close()

		fmt.Println(err)

		return
	}

	us := connection.Handle

	go func(wc io.WriteCloser, r io.Reader) {
		defer wc.Close()
		io.Copy(wc, r)
	}(ds, us)

	bg := protocol.NewConnection(ds)

	for _, item := range connection.PacketQueue {
		id, err := bg.Write(item)

		if err != nil {
			fmt.Printf("Error in packet queue: %s\n", err)
		} else {
			fmt.Printf("Wroted packet id #%d\n", id)
		}
	}

	go func(wc io.WriteCloser, r io.Reader) {
		defer wc.Close()
		io.Copy(wc, r)
	}(us, ds)
}
