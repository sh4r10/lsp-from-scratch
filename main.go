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
		diagnostics := state.OpenDocument(request.Params.TextDocument.URI, request.Params.TextDocument.Text)
		response := lsp.PublishDiagnosticsNotification{
			Notification: lsp.Notification{
				RPC:    "2.0",
				Method: "textDocument/publishDiagnostics",
			},
			Params: lsp.PublishDiagnosticsParams{
				URI:         request.Params.TextDocument.URI,
				Diagnostics: diagnostics,
			},
		}
		writeResponse(writer, response)
	case "textDocument/didChange":
		var request lsp.TextDocumentDidChangeNotification
		if err := json.Unmarshal(content, &request); err != nil {
			logger.Printf("Couldn't parse textDocument/didChange: %s", err)
			return
		}
		logger.Printf("Changed  file: %s", request.Params.TextDocument.URI)
		for _, change := range request.Params.ContentChanges {
			diagnostics := state.UpdateDocument(request.Params.TextDocument.URI, change.Text)
			response := lsp.PublishDiagnosticsNotification{
				Notification: lsp.Notification{
					RPC:    "2.0",
					Method: "textDocument/publishDiagnostics",
				},
				Params: lsp.PublishDiagnosticsParams{
					URI:         request.Params.TextDocument.URI,
					Diagnostics: diagnostics,
				},
			}
			writeResponse(writer, response)
		}
	case "textDocument/hover":
		var request lsp.TextDocumentHoverRequest
		if err := json.Unmarshal(content, &request); err != nil {
			logger.Printf("Couldn't parse textDocument/hover: %s", err)
			return
		}
		response := state.Hover(&request.ID, request.Params.TextDocument.URI, request.Params.Position)
		writeResponse(writer, response)
	case "textDocument/definition":
		var request lsp.TextDocumentDefinitionRequest
		if err := json.Unmarshal(content, &request); err != nil {
			logger.Printf("Couldn't parse textDocument/definition: %s", err)
			return
		}
		response := state.Definition(&request.ID, request.Params.TextDocument.URI, request.Params.Position)
		writeResponse(writer, response)

	case "textDocument/codeAction":
		var request lsp.TextDocumentCodeActionsRequest
		if err := json.Unmarshal(content, &request); err != nil {
			logger.Printf("Couldn't parse textDocument/codeAction: %s", err)
			return
		}

		response := state.CodeAction(&request.ID, request.Params.TextDocument.URI)
		writeResponse(writer, response)

	case "textDocument/completion":
		var request lsp.TextDocumentCompletionRequest
		if err := json.Unmarshal(content, &request); err != nil {
			logger.Printf("Couldn't parse textDocument/codeAction: %s", err)
			return
		}

		response := state.Completion(&request.ID, request.Params.TextDocument.URI)
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
