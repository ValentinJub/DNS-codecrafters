package server

import (
	"bytes"
	"encoding/binary"
	"fmt"
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
	b := encodeLabel(q.Name)
	bw := NewBitWriter(4)
	bw.Write16Bit(q.RecordType)
	bw.Write16Bit(q.Class)
	b = append(b, bw.buffer...)
	return b
}

func DecodeDNSQuestion(data []byte) *DNSQuestion {
	name, offset := decodeLabel(data)
	typ := buffToInt16(data[offset+1 : offset+3])
	class := buffToInt16(data[offset+3 : offset+5])
	return NewDNSQuestion(name, uint16(typ), uint16(class))
}

func decodeLabel(data []byte) (string, int) {
	sep := int(data[0])
	str := ""
	for x, c := range data[1:] {
		if sep == 0 {
			if c == '\x00' {
				return str, x + 1
			} else {
				sep = int(c)
				str += "."
			}
		} else {
			str += string(c)
			sep--
		}
	}
	return "", -1
}

func buffToInt16(buff []byte) int16 {
	var n int16
	if err := binary.Read(bytes.NewReader(buff), binary.BigEndian, &n); err != nil {
		fmt.Println(err)
		return 0
	}
	return n
}

func encodeLabel(name string) []byte {
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
