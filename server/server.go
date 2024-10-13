package server

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/dns-server-starter-go/logger"
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
	// Conventionally, DNS packets are sent using UDP transport and are limited to 512 bytes
	buf := make([]byte, 512)

	for {
		size, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			return err
		}

		receivedData := buf[:size]
		logger.LogIOData([]byte(receivedData), 0)
		header := DecodeDNSHeader(receivedData) // Header should be 12 byte long
		question := DecodeDNSQuestion(receivedData[12:])

		response := createResponse(*header, *question)
		logger.LogIOData(response, 1)
		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			return err
		}
	}
}

func createResponse(header DNSHeader, question DNSQuestion) []byte {
	if header.RecursionDesired {
		fmt.Println("Recurssion is true")
	}
	h := NewDNSHeader(header.PacketIdentifier, true, header.OperationCode, false, true, header.RecursionDesired, true, 0, 4, 1, 1, 0, 0)
	a := NewDNSAnswer(question, 60, 4, []byte{8, 8, 8, 8})
	resp := append(h.Encode(), question.Encode()...)
	return append(resp, a.Encode()...)
}
