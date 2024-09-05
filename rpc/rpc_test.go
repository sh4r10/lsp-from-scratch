package rpc_test

import (
	"lsp-from-scratch/rpc"
	"testing"
)

type ExampleMessage struct {
	Name string
}

func TestEncoding(t *testing.T) {
	expected := "Content-Length: 17\r\n\r\n{\"Name\":\"shariq\"}"
	actual := rpc.EncodeMessage(ExampleMessage{Name: "shariq"})
	if expected != actual {
		t.Fatalf("\nExpected: \r\n%s \n\nGot: \r\n%s", expected, actual)
	}
}

func TestDecode(t *testing.T) {
	incoming := "Content-Length: 19\r\n\r\n{\"method\":\"shariq\"}"
	expectedLength := 19
	method, content, err := rpc.DecodeMessage([]byte(incoming))
	contentLength := len(content)
	if err != nil {
		t.Fatal(err)
	}
	if expectedLength != contentLength {
		t.Fatalf("\nExpected: %d \nGot: %d", expectedLength, contentLength)
	}

	if method != "shariq" {
		t.Fatalf("\nExpected: %s \nGot: %s", "shariq", method)
	}
}
