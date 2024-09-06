package lsp

type TextDocumentDidChangeNotification struct {
	Notification
	Params DidChangeTextDocumentParams `json:"params" `
}

type DidChangeTextDocumentParams struct {
	TextDocument   VersionedTextDocumentIdentifier  `json:"textDocument"`
	ContentChanges []TextDocumentContentChangeEvent `json:"contentChanges"`
}

type TextDocumentContentChangeEvent struct {
	// the updated text of the whole document, as docSync is set to 1
	// can be used to get incremental changes for a more advanced lsp
	Text string `json:"text"`
}
