package server

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/dns-server-starter-go/logger"
)

type Server struct {
	address           string
	port              string
	forward           string
	localForwardAddr  *net.UDPAddr
	remoteForwardAddr *net.UDPAddr
	connForward       *net.UDPConn
	cache             map[string]DNSAnswer
}

func NewServer(address, port string) *Server {
	return &Server{address: address, port: port, cache: make(map[string]DNSAnswer)}
}

func (s *Server) NewForwarder(destination string) {
	s.forward = destination
}

// Initialise a UDP Address we can listen from
func (s *Server) InitUDPEndpoint() (*net.UDPAddr, error) {
	return net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%s", s.address, s.port))
}

// Initialise a UDP Address we can listen from, it'll be use to listen to response to forwarded requests
func (s *Server) InitForwarderUDPEndpoint() (*net.UDPAddr, error) {
	return net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%s", s.address, "2054"))
}

func (s *Server) InitForwarderDestinationEndpoint() (*net.UDPAddr, error) {
	return net.ResolveUDPAddr("udp", s.forward)
}

// Listening on UDP Address and handles UDP connections and responses
func (s *Server) ListenUDP(udpAddress *net.UDPAddr) error {
	udpConn, err := net.ListenUDP("udp", udpAddress)
	if err != nil {
		return err
	}
	defer udpConn.Close()
	err = s.handleUDPEndpoint(*udpConn)
	return err
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

		clientData := buf[:size]
		logger.LogIOData([]byte(clientData), 0, false)

		decoder := NewDecoder(clientData)

		header, err := decoder.DecodeHeader() // Header should be 12 byte long
		if err != nil {
			return err
		}
		questions := decoder.DecodeQuestions()

		if s.forward != "" { // Forward each question to the forwarder address, the forwarder ony supports one question at a time
			answers, err := s.forwardQuestions(header, questions)
			if err != nil {
				return err
			}
			resp := makeResponse(header, questions, answers)
			logger.LogIOData(resp, 1, false)
			_, err = udpConn.WriteToUDP(resp, source)
			if err != nil {
				return err
			}
		} else { // Forge a semi-hardcoded answer, the IPs returned all point to 8.8.8.8
			response := mockResponse(header, questions)
			logger.LogIOData(response, 1, false)
			_, err = udpConn.WriteToUDP(response, source)
			if err != nil {
				return err
			}
		}

	}
}

func makeResponse(header DNSHeader, questions []DNSQuestion, answers []DNSAnswer) []byte {
	qs := []byte{}
	for _, q := range questions {
		qs = append(qs, q.Encode()...)
	}
	as := []byte{}
	for _, a := range answers {
		as = append(as, a.Encode()...)
	}
	header.SetResponse()
	header.AnswerCount = uint16(len(answers))
	resp := append(header.Encode(), qs...)
	resp = append(resp, as...)
	return resp
}

func (s *Server) dialUpResolver() error {
	local, err := s.InitForwarderUDPEndpoint()
	if err != nil {
		fmt.Println("Can't setup listener port for forwarder")
		return err
	}
	remote, err := s.InitForwarderDestinationEndpoint()
	if err != nil {
		fmt.Println("Can't connect to the forwarder")
		return err
	}
	// Establish connection with forwarding server
	conn, err := net.DialUDP("udp", local, remote)
	if err != nil {
		return err
	}
	s.localForwardAddr = local
	s.remoteForwardAddr = remote
	s.connForward = conn
	return nil
}

// True if we're still connecter to Resolver
func (s *Server) isForwarderConnAlive() bool {
	return s.connForward != nil
}

// Forward question to Resolver if the address isn't in the Server Cache
func (s *Server) forwardQuestions(header DNSHeader, questions []DNSQuestion) ([]DNSAnswer, error) {
	if !s.isForwarderConnAlive() {
		err := s.dialUpResolver()
		if err != nil {
			return []DNSAnswer{}, err
		}
	}
	answers := make([]DNSAnswer, 0)
	for _, q := range questions {
		if answer, found := s.cache[q.Name]; found {
			answers = append(answers, answer)
			continue
		}
		request := createRequest(header, q)
		_, err := s.connForward.Write(request)
		if err != nil {
			return []DNSAnswer{}, err
		}
		logger.LogIOData(request, 1, true)

		buff := make([]byte, 512)
		n, _, err := s.connForward.ReadFromUDP(buff)
		if err != nil {
			return []DNSAnswer{}, err
		}

		resolverData := buff[:n]
		logger.LogIOData([]byte(resolverData), 0, true)
		// Since the answer is made up of the request + the answer, we caluclate the request lenght and start reading the answer from this offset
		decoder := NewDecoder(resolverData)
		answer := decoder.DecodeAnswer(len(request))
		s.logToCache(q.Name, answer)
		answers = append(answers, answer)
	}
	return answers, nil
}

// Memoize known addresses, reducing the need to query the resolver
func (s *Server) logToCache(question string, answer DNSAnswer) {
	s.cache[question] = answer
}

// For one question only
func createRequest(header DNSHeader, question DNSQuestion) []byte {
	header.SetQuery()
	header.QuestionCount = 1
	return append(header.Encode(), question.Encode()...)
}

// Mock Response, implemented when no Resolver is provided
func mockResponse(header DNSHeader, questions []DNSQuestion) []byte {
	h := NewDNSHeader(
		header.PacketIdentifier,
		true,
		header.OperationCode,
		false,
		true,
		header.RecursionDesired,
		true,
		0,
		4,
		uint16(len(questions)), // QuestionCount
		uint16(len(questions)), // AnswerCount
		0,
		0,
	)
	qs := make([]byte, 0)
	for _, q := range questions {
		qs = append(qs, q.Encode()...)
	}
	resp := append(h.Encode(), qs...)
	a := NewDNSAnswer(
		questions,
		60,
		4,
		[]int8{8, 8, 8, 8},
	)
	return append(resp, a.Encode()...)
}
