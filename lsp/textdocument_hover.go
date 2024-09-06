package lsp

type TextDocumentHoverRequest struct {
	Request
	Params HoverParams `json:"params"`
}
type HoverParams struct {
	TextDocumentPositionParams
}

type HoverResult struct {
	Contents string `json:"contents"`
}

type TextDocumentHoverResponse struct {
	Response
	Result HoverResult `json:"result"`
}
