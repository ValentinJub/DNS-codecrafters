package server

import (
	"bytes"
	"strings"
)

type DNSQuestion struct {
	Name       string // ex: google.com - encoded as a sequence of labels
	RecordType uint16 // https://www.rfc-editor.org/rfc/rfc1035#section-3.2.2
	Class      uint16 // The class in practice is always set to 1
}

func NewDNSQuestion(name string, recordType, class uint16) *DNSQuestion {
	return &DNSQuestion{Name: name, RecordType: recordType, Class: class}
}

func (q *DNSQuestion) Encode() []byte {
	b := encodeName(q.Name)
	bw := NewBitWriter(4)
	bw.Write16Bit(q.RecordType)
	bw.Write16Bit(q.Class)
	b = append(b, bw.buffer...)
	return b
}

func encodeName(name string) []byte {
	labels := strings.Split(name, ".")
	buff := new(bytes.Buffer)
	for _, label := range labels {
		lname := []byte(label)
		buff.Write([]byte{byte(len(lname))})
		buff.Write(lname)
	}
	buff.Write([]byte{'\x00'}) // padding
	return buff.Bytes()
}
