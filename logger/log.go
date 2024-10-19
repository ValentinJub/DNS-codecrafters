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
func LogIOData(data []byte, direction int, resolver bool) {
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
	opt := "to Client"
	if direction == 0 {
		dir = "Incoming"
		opt = "from Client"
	}
	if resolver {
		if direction == 0 {
			dir = "Incoming"
			opt = "from Resolver"
		} else {
			opt = "to Resolver"
		}
	}
	fmt.Printf("%s%s data %s (hex):\n%s%s\n\n", Pink, dir, opt, xclean, Reset)
}
