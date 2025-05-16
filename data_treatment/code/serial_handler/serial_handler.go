package serial_handler

import (
	"io"
	"strings"
)

func ReceiveDataFromPort(s io.ReadCloser) (string, error) {
	buffer := make([]byte, 1)
	var builder strings.Builder
	braceCount := 0
	reading := false

	for {
		n, err := s.Read(buffer)
		if err != nil {
			return "", err
		}
		if n == 0 {
			continue
		}
		b := buffer[0]
		if b == '{' {
			braceCount++
			reading = true
		}
		if reading {
			builder.WriteByte(b)
		}
		if b == '}' {
			braceCount--
			if braceCount == 0 {
				return builder.String(), nil
			}
		}
	}
}
