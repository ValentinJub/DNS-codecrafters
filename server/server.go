package server

import (
	"fmt"
	"net"
)

type DNServer struct {
	address string
	port    string
}

func NewDNServer(address, port string) *DNServer {
	return &DNServer{address: address, port: port}
}

// Initialise a UDP Address we can listen from
func (s *DNServer) InitUDPEndpoint() (*net.UDPAddr, error) {
	return net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%s", s.address, s.port))
}

// Listening on UDP Address and handles UDP connections and responses
func (s *DNServer) ListenUDP(udpAddress *net.UDPAddr) error {
	udpConn, err := net.ListenUDP("udp", udpAddress)
	if err != nil {
		return err
	}
	defer udpConn.Close()
	s.handleUDPEndpoint(*udpConn)
	return nil
}

// Read from the UDP Endpoint and send response
func (s *DNServer) handleUDPEndpoint(udpConn net.UDPConn) error {
	buf := make([]byte, 512)

	for {
		size, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			return err
		}

		receivedData := string(buf[:size])
		fmt.Printf("Received %d bytes from %s: %s\n", size, source, receivedData)

		response := createResponse()

		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			return err
		}
	}
}

func createResponse() []byte {
	h := NewDNSHeader(1234, false, 0, true, true, true, true, 0, 0, 0, 0, 0, 0)
	return h.Encode()
}
