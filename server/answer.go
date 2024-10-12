package server

import "encoding/binary"

type DNSAnswer struct {
	DNSQuestion
	TTL    uint32
	Length uint16
	Data   []byte
}

func NewDNSAnswer(question DNSQuestion, ttl uint32, len uint16, data []byte) *DNSAnswer {
	return &DNSAnswer{
		DNSQuestion: question,
		TTL:         ttl,
		Length:      len,
		Data:        data,
	}
}

func (a *DNSAnswer) Encode() []byte {
	b := a.DNSQuestion.Encode()
	binary.BigEndian.AppendUint32(b, a.TTL)
	binary.BigEndian.AppendUint16(b, a.Length)
	b = append(b, a.Data...)
	return b
}
