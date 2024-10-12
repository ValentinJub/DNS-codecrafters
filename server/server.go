package server

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/dns-server-starter-go/logger"
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
	// Conventionally, DNS packets are sent using UDP transport and are limited to 512 bytes
	buf := make([]byte, 512)

	for {
		size, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			return err
		}

		receivedData := buf[:size]
		logger.LogIOData([]byte(receivedData), 0)
		header := DecodeDNSHeader(receivedData)

		response := createResponse(*header)
		logger.LogIOData(response, 1)
		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			return err
		}
	}
}

func createResponse(header DNSHeader) []byte {
	if header.RecursionDesired {
		fmt.Println("Recurssion is true")
	}
	h := NewDNSHeader(header.PacketIdentifier, true, header.OperationCode, false, true, header.RecursionDesired, true, 0, 4, 1, 1, 0, 0)
	q := NewDNSQuestion("codecrafters.io", 1, 1)
	a := NewDNSAnswer(*q, 60, 4, []byte{8, 8, 8, 8})
	resp := append(h.Encode(), q.Encode()...)
	return append(resp, a.Encode()...)
}
