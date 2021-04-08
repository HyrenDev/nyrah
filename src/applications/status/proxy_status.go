package status

import (
	"log"
	"net"
)

func IsProxyOnline(server string) bool {
	_, err := net.Dial("tcp", server)

	if err != nil {
		return false
	}

	return true
}

func getProxyOnlineCount(server string) (int, error) {
	connection, err := net.Dial("tcp", server)

	if err != nil {
		return 0, err
	}

	receivedBuf := make([]byte, 1024)

	n, err := connection.Read(receivedBuf)

	log.Println(n)

	return 0, nil
}

