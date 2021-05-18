package connection

import (
	"fmt"
	"io"
	"net"
	"net/hyren/nyrah/minecraft/protocol"

	ProxyApp "net/hyren/nyrah/applications"
)

func copy(wc io.WriteCloser, r io.Reader) {
	defer wc.Close()
	io.Copy(wc, r)
}

func SendToProxy(connection *protocol.Connection, name string) {
	var inetSocketAddress = ProxyApp.GetProxyAddress(name)

	ds, err := net.Dial("tcp", fmt.Sprintf("%s:%d", inetSocketAddress.GetHostAddress(), inetSocketAddress.GetPort()))

	if err != nil {
		connection.Close()

		fmt.Println(err)

		return
	}

	us := connection.Handle

	go copy(ds, us)

	bg := protocol.NewConnection(ds)

	for _, item := range connection.PacketQueue {
		id, err := bg.Write(item)

		if err != nil {
			fmt.Printf("Error in packet queue: %s\n", err)
		} else {
			fmt.Printf("Wroted packet id #%d\n", id)
		}
	}

	go copy(us, ds)
}
