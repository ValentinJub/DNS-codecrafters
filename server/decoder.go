package server

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// type RequestHeader = request.RequestHeader

type Decoder struct {
	data   []byte //this is the raw request bytes
	max    int    //max index, going over will panic
	offset int    //to keep track of where we are when decoding the data
}

func NewDecoder(d []byte) *Decoder {
	return &Decoder{data: d, max: len(d)}
}

// Returns the header of the request, allowing to figure out what type of request we have received
func (d *Decoder) DecodeHeader() (DNSHeader, error) {
	if len(d.data) < 12 {
		return DNSHeader{}, fmt.Errorf("error: data length is under 12 bytes, cannot decode header")
	}
	br := NewBitReader(d.data)
	d.offset = 12
	return DNSHeader{
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
	}, nil
}

// To be used to decode a request, it doesn't stop until the
func (d *Decoder) DecodeQuestions() (questions []DNSQuestion) {
	for d.offset < len(d.data) {
		questions = append(questions, d.DecodeQuestion())
	}
	return
}

func (d *Decoder) DecodeQuestion() DNSQuestion {
	label, inc := decodeLabel(d.data, d.offset, false)
	d.offset += inc + 1
	// fmt.Printf("Decoded label: %s\n", label)
	typ := d.readInt16()
	class := d.readInt16()
	// fmt.Printf("The offset is: %d and the len of the data is: %d \n", d.offset, len(d.data))
	return *NewDNSQuestion(label, uint16(typ), uint16(class))
}

func (d *Decoder) DecodeAnswer(offset int) DNSAnswer {
	d.offset = offset
	question := d.DecodeQuestion()
	ttl := d.readInt32()
	len := d.readInt16()
	ip := []int8{
		d.readInt8(),
		d.readInt8(),
		d.readInt8(),
		d.readInt8(),
	}
	// fmt.Printf("Decoded IP: %v\n", ip)
	return *NewDNSAnswer([]DNSQuestion{question}, uint32(ttl), uint16(len), ip)
}

// Decode a label at starting position `offset` in `data`, `noref` = true means we don't look for pointers, useful to search for a label reference
func decodeLabel(data []byte, offset int, noref bool) (label string, inc int) {
	// The label parts are prefixed with their length
	numOfCharToRead := int(data[offset])
	// fmt.Printf("Num of chars to read: %d\n", numOfCharToRead)
	offset++
	for {
		char := data[offset]
		// fmt.Printf("Char read:\t%c\n", char)
		if !noref {
			br := NewBitReader(data[offset:])
			v := br.ReadBits(2)
			if v == 3 { // This is a pointer
				index := br.ReadBits(14)
				labelParts, _ := decodeLabel(data, int(index), true)
				return label + "." + labelParts, inc + 2
			}
		}
		if char == '\x00' { // It's the end of the label
			return label, inc + 1
		}
		if numOfCharToRead > 0 {
			label += string(char)
			numOfCharToRead--
		} else { // It's another label part, register its length
			numOfCharToRead = int(char)
			label += "."
		}
		offset++
		inc++
		if offset >= len(data) {
			fmt.Println("error, about to read more data than there is in decode label")
			return label, inc
		}
	}
}

func (d *Decoder) readInt32() int32 {
	if d.offset+4 > d.max {
		fmt.Printf("Error: trying to read from pos: %d to %d when the max is %d", d.offset, d.offset+4, d.max)
		return 0
	}
	n := buffToInt32(d.data[d.offset : d.offset+4])
	d.offset += 4
	return n
}

func (d *Decoder) readInt16() int16 {
	if d.offset+2 > d.max {
		fmt.Printf("Error: trying to read from pos: %d to %d when the max is %d", d.offset, d.offset+2, d.max)
		return 0
	}
	n := buffToInt16(d.data[d.offset : d.offset+2])
	d.offset += 2
	return n
}

func (d *Decoder) readInt8() int8 {
	if d.offset+1 > d.max {
		fmt.Printf("Error: trying to read from pos: %d to %d when the max is %d", d.offset, d.offset+1, d.max)
		return 0
	}
	n := buffToInt8(d.data[d.offset : d.offset+1])
	d.offset += 1
	return n
}

func buffToInt32(buff []byte) int32 {
	var n int32
	if err := binary.Read(bytes.NewReader(buff), binary.BigEndian, &n); err != nil {
		fmt.Println(err)
		return 0
	}
	return n
}

func buffToInt16(buff []byte) int16 {
	var n int16
	if err := binary.Read(bytes.NewReader(buff), binary.BigEndian, &n); err != nil {
		fmt.Println(err)
		return 0
	}
	return n
}

func buffToInt8(buff []byte) int8 {
	var n int8
	if err := binary.Read(bytes.NewReader(buff), binary.BigEndian, &n); err != nil {
		fmt.Println(err)
		return 0
	}
	return n
}
