package server

import "fmt"

type BitWriter struct {
	buffer []byte
	pos    int // Bit position within the buffer
}

// NewBitWriter creates a new BitWriter for a buffer of a given size (in bytes)
func NewBitWriter(size int) *BitWriter {
	return &BitWriter{
		buffer: make([]byte, size),
		pos:    0,
	}
}

// WriteBits writes `value` into the buffer using `bitCount` bits
func (bw *BitWriter) WriteBits(value uint32, bitCount int) error {
	if bitCount > 32 {
		return fmt.Errorf("bitCount exceeds 32 bits")
	}

	for bitCount > 0 {
		bytePos := bw.pos / 8
		bitPos := bw.pos % 8
		bitsInCurrentByte := 8 - bitPos

		// Calculate how many bits we can write in the current byte
		bitsToWrite := bitCount
		if bitsToWrite > bitsInCurrentByte {
			bitsToWrite = bitsInCurrentByte
		}

		// Shift value to get the most significant bits we are writing into the right place
		shiftedValue := (value >> uint(bitCount-bitsToWrite)) & ((1 << uint(bitsToWrite)) - 1)

		// Write the bits into the buffer
		bw.buffer[bytePos] |= byte(shiftedValue << (8 - bitPos - bitsToWrite))

		// Update the bit position and decrease bitCount
		bw.pos += bitsToWrite
		bitCount -= bitsToWrite
	}

	return nil
}

func (bw *BitWriter) Buffer() []byte {
	return bw.buffer
}

func (bw *BitWriter) Write1Bit(value bool) {
	var n uint32
	// Set to 1 if the value is false, otherwise 0 is the default
	if !value {
		n = 1
	}
	err := bw.WriteBits(n, 1)
	if err != nil {
		fmt.Printf("error while writing 1 bit to the bitwriter buffer: %s", err)
	}
}

func (bw *BitWriter) Write3Bit(value uint8) {
	err := bw.WriteBits(uint32(value), 3)
	if err != nil {
		fmt.Printf("error while writing 3 bits to the bitwriter buffer: %s", err)
	}
}

func (bw *BitWriter) Write4Bit(value uint8) {
	err := bw.WriteBits(uint32(value), 4)
	if err != nil {
		fmt.Printf("error while writing 4 bit to the bitwriter buffer: %s", err)
	}
}

func (bw *BitWriter) Write16Bit(value uint16) {
	err := bw.WriteBits(uint32(value), 16)
	if err != nil {
		fmt.Printf("error while writing 16 bit to the bitwriter buffer: %s", err)
	}
}
