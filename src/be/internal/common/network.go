package common

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"strings"
)

// ReadHttpFromTcp read http request or response from tcp connection
func ReadHttpFromTcp(c net.Conn) ([]byte, error) {
	var req bytes.Buffer
	reader := bufio.NewReader(c)
	contentLength := 0
	for {
		buf, err := reader.ReadBytes('\n')
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, err
		}
		if _, err := req.Write(buf); err != nil {
			return nil, err
		}

		if strings.HasPrefix(string(buf), "Content-Length") {
			_, err := fmt.Sscanf(string(buf), "Content-Length: %d", &contentLength)
			if err != nil {
				return nil, err
			}
		}
		if string(buf) == "\r\n" {
			break
		}
	}
	buf := make([]byte, contentLength)
	_, err := io.ReadFull(reader, buf)
	if err != nil {
		return nil, err
	}
	if _, err := req.Write(buf); err != nil {
		return nil, err
	}

	return req.Bytes(), nil
}
