package server

type DNSAnswer struct {
	Questions []DNSQuestion
	TTL       uint32
	Length    uint16
	IP        []int8
}

func NewDNSAnswer(questions []DNSQuestion, ttl uint32, len uint16, data []int8) *DNSAnswer {
	return &DNSAnswer{
		Questions: questions,
		TTL:       ttl,
		Length:    len,
		IP:        data,
	}
}

func (a *DNSAnswer) Encode() []byte {
	b := make([]byte, 0)
	for _, q := range a.Questions {
		b = append(b, q.Encode()...)
		bw := NewBitWriter(6)
		bw.WriteBits(a.TTL, 32)
		bw.WriteBits(uint32(a.Length), 16)
		b = append(b, bw.Buffer()...)
		ip := []byte{byte(a.IP[0]), byte(a.IP[1]), byte(a.IP[2]), byte(a.IP[3])}
		b = append(b, ip...)
	}
	return b
}
