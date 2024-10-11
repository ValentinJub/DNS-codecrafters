package main

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/dns-server-starter-go/server"
)

const (
	SERVER_ADDRESS = "127.0.0.1"
	PORT           = "2053"
)

func main() {
	server := server.NewServer(SERVER_ADDRESS, PORT)
	udpAddr, err := server.InitUDPEndpoint()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while initialising UDP endpoint: %s", err)
		return
	}
	server.ListenUDP(udpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while listening on UDP endpoint: %s", err)
		return
	}
}
