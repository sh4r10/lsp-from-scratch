package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"lsp-from-scratch/analyzer"
	"lsp-from-scratch/lsp"
	"lsp-from-scratch/rpc"
	"os"
)

func main() {
	logger := getLogger("/mnt/lts/projects/lsp-from-scratch/log.txt")
	logger.Println("Starting to log now")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)

	state := analyzer.NewState()

	for scanner.Scan() {
		message := scanner.Bytes()
		method, content, err := rpc.DecodeMessage(message)
		if err != nil {
			logger.Printf("Error: %s", err)
			continue
		}
		writer := os.Stdout
		handleMessage(logger, writer, state, method, content)
	}
}

func handleMessage(logger *log.Logger, writer io.Writer, state analyzer.State, method string, content []byte) {
	logger.Printf("Received msg with method: %s", method)
	switch method {
	case "initialize":
		var request lsp.InitializeRequest
		if err := json.Unmarshal(content, &request); err != nil {
			logger.Printf("Hey, I couldn't parse this: %s", err)
		}
		logger.Printf("Connected to: %s %s",
			request.Params.ClientInfo.Name,
			request.Params.ClientInfo.Version)
		// our response
		res := lsp.NewInitializeResponse(request.ID)
		writeResponse(writer, res)
		logger.Println("Sent the initilize response")
	case "textDocument/didOpen":
		var request lsp.TextDocumentDidOpenNotification
		if err := json.Unmarshal(content, &request); err != nil {
			logger.Printf("Couldn't parse textDocument/didOpen: %s", err)
			return
		}
		logger.Printf("Opened file: %s", request.Params.TextDocument.URI)
		state.OpenDocument(request.Params.TextDocument.URI, request.Params.TextDocument.Text)
	case "textDocument/didChange":
		var request lsp.TextDocumentDidChangeNotification
		if err := json.Unmarshal(content, &request); err != nil {
			logger.Printf("Couldn't parse textDocument/didChange: %s", err)
			return
		}
		logger.Printf("Changed  file: %s", request.Params.TextDocument.URI)
		for _, change := range request.Params.ContentChanges {
			state.UpdateDocument(request.Params.TextDocument.URI, change.Text)
		}
	case "textDocument/hover":
		var request lsp.TextDocumentHoverRequest
		if err := json.Unmarshal(content, &request); err != nil {
			logger.Printf("Couldn't parse textDocument/hover: %s", err)
			return
		}

		// Create a response and write it back
		response := state.Hover(&request.ID, request.Params.TextDocument.URI, request.Params.Position)
		writeResponse(writer, response)
	}
}

func writeResponse(writer io.Writer, msg any) {
	enc := rpc.EncodeMessage(msg)
	writer.Write([]byte(enc))
}

func getLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		panic("Hey, give me a good file!")
	}
	return log.New(logfile, "[custom-lsp]", log.Ldate|log.Ltime|log.Lshortfile)
}
