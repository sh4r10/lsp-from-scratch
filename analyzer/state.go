package analyzer

import (
	"fmt"
	"lsp-from-scratch/lsp"
)

type State struct {
	// map of filenames to their content
	Documents map[string]string
}

func NewState() State {
	return State{Documents: map[string]string{}}
}

func (s *State) OpenDocument(uri, text string) {
	s.Documents[uri] = text
}

func (s *State) UpdateDocument(uri, text string) {
	s.Documents[uri] = text
}

func (s *State) Hover(id *int, uri string, position lsp.Position) lsp.TextDocumentHoverResponse {
	// in an actual lsp this would look up the type through some kind of type
	// analyzer tools / code.
	document := s.Documents[uri]
	response := lsp.TextDocumentHoverResponse{
		Response: lsp.Response{
			ID:  id,
			RPC: "2.0",
		},
		Result: lsp.HoverResult{
			Contents: fmt.Sprintf("File %s, Chars: %d\nCurrent Line: %d, Current Char: %d",
				uri,
				len(document),
				position.Line,
				position.Character,
			),
		},
	}
	return response
}
