package analyzer

import (
	"fmt"
	"lsp-from-scratch/lsp"
	"strings"
)

type State struct {
	// map of filenames to their content
	Documents map[string]string
}

func NewState() State {
	return State{Documents: map[string]string{}}
}

func getDiagnostics(text string) []lsp.Diagnostic {
	list := []lsp.Diagnostic{}
	for row, line := range strings.Split(text, "\n") {
		if strings.Contains(line, "VS Code") {
			index := strings.Index(line, "VS Code")
			list = append(list, lsp.Diagnostic{
				Range:    LineRange(row, index, index+len("VS Code")),
				Severity: 1,
				Source:   "trust me bro",
				Message:  "Please avoid profanities, such as VS Code",
			})
		}
		if strings.Contains(line, "Neovim") {
			index := strings.Index(line, "Neovim")
			list = append(list, lsp.Diagnostic{
				Range:    LineRange(row, index, index+len("Neovim")),
				Severity: 4,
				Source:   "trust me bro",
				Message:  "Exquisite choice, if I may say so <3",
			})
		}
	}
	return list
}

func (s *State) OpenDocument(uri, text string) []lsp.Diagnostic {
	s.Documents[uri] = text
	return getDiagnostics(text)
}

func (s *State) UpdateDocument(uri, text string) []lsp.Diagnostic {
	s.Documents[uri] = text
	return getDiagnostics(text)
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

func (s *State) Definition(id *int, uri string, position lsp.Position) lsp.TextDocumentDefinitionResponse {
	// in an actual lsp this would look up the definition
	// for this case, the definition is always one line above the current line
	response := lsp.TextDocumentDefinitionResponse{
		Response: lsp.Response{
			ID:  id,
			RPC: "2.0",
		},
		Result: lsp.Location{
			URI: uri,
			Range: lsp.Range{
				Start: lsp.Position{
					Line:      position.Line - 1,
					Character: 0,
				},
				End: lsp.Position{
					Line:      position.Line - 1,
					Character: 0,
				},
			},
		},
	}
	return response
}

func (s *State) CodeAction(id *int, uri string) lsp.TextDocumentCodeActionsResponse {
	text := s.Documents[uri]
	actions := []lsp.CodeAction{}
	for row, line := range strings.Split(text, "\n") {
		idx := strings.Index(line, "VS Code")
		if idx >= 0 {
			replaceChange := map[string][]lsp.TextEdit{}
			replaceChange[uri] = []lsp.TextEdit{
				{
					Range:   LineRange(row, idx, idx+len("VS Code")),
					NewText: "NeoVIM",
				},
			}

			actions = append(actions, lsp.CodeAction{
				Title: "Replace VS Code with the GOAT editor",
				Edit:  &lsp.WorkspaceEdit{Changes: replaceChange},
			})

			censorChange := map[string][]lsp.TextEdit{}
			censorChange[uri] = []lsp.TextEdit{
				{
					Range:   LineRange(row, idx, idx+len("VS Code")),
					NewText: "VS C**e",
				},
			}

			actions = append(actions, lsp.CodeAction{
				Title: "Censor VS Code",
				Edit:  &lsp.WorkspaceEdit{Changes: censorChange},
			})
		}
	}

	response := lsp.TextDocumentCodeActionsResponse{
		Response: lsp.Response{
			ID:  id,
			RPC: "2.0",
		},
		Result: actions,
	}

	return response
}

func (s *State) Completion(id *int, uri string) lsp.TextDocumentCompletionResponse {
	items := []lsp.CompletionItem{
		{
			Label:         "LSP From Scratch",
			Detail:        "Very cool project",
			Documentation: "A custom implementation of the LSP Protocol written in Go",
		},
		{
			Label:         "Neovim over VSCode",
			Detail:        "Begone Microsoft",
			Documentation: "Wayy faster",
		},
	}
	// use static analysis tools to get good text-completions
	response := lsp.TextDocumentCompletionResponse{
		Response: lsp.Response{
			ID:  id,
			RPC: "2.0",
		},
		Result: items,
	}

	return response
}

func LineRange(line, start, end int) lsp.Range {
	return lsp.Range{
		Start: lsp.Position{
			Line:      line,
			Character: start,
		},
		End: lsp.Position{
			Line:      line,
			Character: end,
		},
	}
}
