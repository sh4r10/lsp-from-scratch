package lsp

type TextDocumentCodeActionsRequest struct {
	Request
	Params TextDocumentCodeActionsParams `json:"params"`
}

type TextDocumentCodeActionsParams struct {
	TextDocument TextDocumentIdentifier         `json:"textDocument"`
	Range        Range                          `json:"range"`
	Context      TextDocumentCodeActionsContext `json:"context"`
}

type TextDocumentCodeActionsContext struct{}

type TextDocumentCodeActionsResponse struct {
	Response
	Result []CodeAction `json:"result"`
}

type CodeAction struct {
	Title   string         `json:"title"`
	Edit    *WorkspaceEdit `json:"edit,omitempty"`
	Command *Command       `json:"command"`
}

type Command struct {
	Title     string        `json:"title"`
	Command   string        `json:"command"`
	Arguments []interface{} `json:"arguments,omitempty"`
}
