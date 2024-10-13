package main

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/dns-server-starter-go/server"
)

const (
	DNServer_ADDRESS = "127.0.0.1"
	PORT             = "2053"
)

func main() {
	DNServer := server.NewServer(DNServer_ADDRESS, PORT)
	udpAddr, err := DNServer.InitUDPEndpoint()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while initialising UDP endpoint: %s", err)
		return
	}
	err = DNServer.ListenUDP(udpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while listening on UDP endpoint: %s", err)
		return
	}
}
