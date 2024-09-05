package rpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

func EncodeMessage(msg any) string {
	// serializes the message, similar to json.stringify
	content, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(content), content)
}

func DecodeMessage(msg []byte) (int, error) {
	header, content, found := bytes.Cut(msg, []byte{'\r', '\n', '\r', '\n'})
	if !found {
		return 0, errors.New("Did not find the separator \\r\\n")
	}
	// get everthing in the header after the content-length section
	contentLengthBytes := header[len("Content-Length: "):]
	// ascii to integer the rest, should be a number in this case
	// because we only have one header
	contentLength, err := strconv.Atoi(string(contentLengthBytes))
	if err != nil {
		return 0, err
	}
	_ = content
	return contentLength, nil
}
