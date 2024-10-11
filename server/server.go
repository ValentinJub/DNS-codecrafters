package server

import (
	"fmt"
	"net"
)

type Server struct {
	address string
	port    string
}

func NewServer(address, port string) *Server {
	return &Server{address: address, port: port}
}

// Initialise a UDP Address we can listen from
func (s *Server) InitUDPEndpoint() (*net.UDPAddr, error) {
	return net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%s", s.address, s.port))
}

// Listening on UDP Address and handles UDP connections and responses
func (s *Server) ListenUDP(udpAddress *net.UDPAddr) error {
	udpConn, err := net.ListenUDP("udp", udpAddress)
	if err != nil {
		return err
	}
	defer udpConn.Close()
	s.handleUDPEndpoint(*udpConn)
	return nil
}

// Read from the UDP Endpoint and send response
func (s *Server) handleUDPEndpoint(udpConn net.UDPConn) error {
	buf := make([]byte, 512)

	for {
		size, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			return err
		}

		receivedData := string(buf[:size])
		fmt.Printf("Received %d bytes from %s: %s\n", size, source, receivedData)

		// Create an empty response
		response := []byte{}

		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			return err
		}
	}
}
