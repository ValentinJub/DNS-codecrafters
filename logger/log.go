package logger

import (
	"encoding/hex"
	"fmt"
)

const (
	Reset = "\033[0m"
	Pink  = "\033[35m"
)

// Log incoming or outgoing request in hexadecimal format, direction takes 0 for incoming and 1 for outgoing
func LogRequest(data []byte, direction int) {
	x := hex.EncodeToString(data)
	xclean := ""
	for i, char := range x {
		if i%32 == 0 && i != 0 {
			xclean += "\n"
		} else if i%2 == 0 && i != 0 {
			xclean += " "
		}
		xclean += string(char)
	}
	dir := "Outgoing"
	if direction == 0 {
		dir = "Incoming"
	}
	fmt.Printf("%s%sData (hex):\n%s%s\n\n", Pink, dir, xclean, Reset)
}
