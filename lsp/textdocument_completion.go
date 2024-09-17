package lsp

type TextDocumentCompletionRequest struct {
	Request
	Params CompletionParams `json:"params"`
}
type CompletionParams struct {
	TextDocumentPositionParams
}

type CompletionItem struct {
	Label         string `json:"label"`
	Detail        string `json:"detail"`
	Documentation string `json:"documentation"`
}

type TextDocumentCompletionResponse struct {
	Response
	Result []CompletionItem `json:"result"`
}
