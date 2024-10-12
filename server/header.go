package server

/*
DNSHeader structure is as 12 bytes structure
*/

type DNSHeader struct {
	PacketIdentifier    uint16
	QueryResponse       bool  // 1 bit
	OperationCode       uint8 // 4 bits
	AuthoritativeAnswer bool  // 1 bit
	TruncatedMessage    bool  // 1 bit
	RecursionDesired    bool  // 1 bit
	RecursionAvailable  bool  // 1 bit
	Reserved            uint8 // 3 bits
	ResponseCode        uint8 // 4 bits
	QuestionCount       uint16
	AnswerCount         uint16
	AuthorityCount      uint16
	AdditionalCount     uint16
}

func NewDNSHeader(id uint16, qr bool, opcode uint8, aa, tc, rd, ra bool, z, rcode uint8, qdcount, ancount, nscount, arcount uint16) *DNSHeader {
	return &DNSHeader{
		PacketIdentifier:    id,
		QueryResponse:       qr,
		OperationCode:       opcode,
		AuthoritativeAnswer: aa,
		TruncatedMessage:    tc,
		RecursionDesired:    rd,
		RecursionAvailable:  ra,
		Reserved:            z,
		ResponseCode:        rcode,
		QuestionCount:       qdcount,
		AnswerCount:         ancount,
		AuthorityCount:      nscount,
		AdditionalCount:     arcount,
	}
}

// Encode in binary format, refer to DNSHeader struct for schema
func (h *DNSHeader) Encode() []byte {
	bw := NewBitWriter(12)
	bw.Write16Bit(h.PacketIdentifier)
	bw.Write1Bit(h.QueryResponse)
	bw.Write4Bit(h.OperationCode)
	bw.Write1Bit(h.AuthoritativeAnswer)
	bw.Write1Bit(h.TruncatedMessage)
	bw.Write1Bit(h.RecursionDesired)
	bw.Write1Bit(h.RecursionAvailable)
	bw.Write3Bit(h.Reserved)
	bw.Write4Bit(h.ResponseCode)
	bw.Write16Bit(h.QuestionCount)
	bw.Write16Bit(h.AnswerCount)
	bw.Write16Bit(h.AuthorityCount)
	bw.Write16Bit(h.AdditionalCount)
	return bw.Buffer()
}

func DecodeDNSHeader(data []byte) *DNSHeader {
	br := NewBitReader(data)
	return &DNSHeader{
		PacketIdentifier:    uint16(br.ReadBits(16)),
		QueryResponse:       br.ReadBits(1) != 0,
		OperationCode:       uint8(br.ReadBits(4)),
		AuthoritativeAnswer: br.ReadBits(1) != 0,
		TruncatedMessage:    br.ReadBits(1) != 0,
		RecursionDesired:    br.ReadBits(1) != 0,
		RecursionAvailable:  br.ReadBits(1) != 0,
		Reserved:            uint8(br.ReadBits(3)),
		ResponseCode:        uint8(br.ReadBits(4)),
		QuestionCount:       uint16(br.ReadBits(16)),
		AnswerCount:         uint16(br.ReadBits(16)),
		AuthorityCount:      uint16(br.ReadBits(16)),
		AdditionalCount:     uint16(br.ReadBits(16)),
	}
}
