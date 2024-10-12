package server

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
	bw := NewBitWriter(6)
	bw.WriteBits(a.TTL, 32)
	bw.WriteBits(uint32(a.Length), 16)
	b = append(b, bw.Buffer()...)
	b = append(b, a.Data...)
	return b
}
