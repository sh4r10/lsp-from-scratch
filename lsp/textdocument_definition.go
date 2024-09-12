package lsp

type TextDocumentDefinitionRequest struct {
	Request
	Params DefinitionParams `json:"params"`
}
type DefinitionParams struct {
	TextDocumentPositionParams
}

type TextDocumentDefinitionResponse struct {
	Response
	Result Location `json:"result"`
}
