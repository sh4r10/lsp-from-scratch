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
	incoming := "Content-Length: 17\r\n\r\n{\"Name\":\"shariq\"}"
	expected := 17
	actual, err := rpc.DecodeMessage([]byte(incoming))
	if err != nil {
		t.Fatal(err)
	}
	if expected != actual {
		t.Fatalf("\nExpected: %d \nGot: %d", expected, actual)
	}
}
