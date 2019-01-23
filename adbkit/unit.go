package adbkit

import (
	"fmt"
	"strings"
)

// 对命令进行编码
func EncodeCommend(command string) []byte {
	prefix := strings.ToUpper("0000" + fmt.Sprintf("%X", len(command)))
	length := len(prefix)
	prefix = prefix[length-4 : length]
	fmt.Println(prefix + command)
	return []byte(prefix + command)
}
