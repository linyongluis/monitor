package plugins

import (
	"bytes"
	"strings"
)

// 删除输出的\x00和多余的空格
func TrimOutput(buffer bytes.Buffer) string {
	return strings.TrimSpace(string(bytes.TrimRight(buffer.Bytes(), "\x00")))
}
