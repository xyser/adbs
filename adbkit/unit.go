package adbkit

import (
	"fmt"
	"strings"
)

func EncodeCommend(command string) []byte {
	prefix := strings.ToUpper("0000" + fmt.Sprintf("%X", len(command)))
	length := len(prefix)
	prefix = prefix[length-4 : length]
	return []byte(prefix + command)
}
