package lsproto

type Method string

const (
	// Since: 3.16.0
	MethodTextDocumentSemanticTokensFull Method = "textDocument/semanticTokens/full"
	// The initialize request is sent from the client to the server.
	// It is sent once as the request after starting up the server.
	// The requests parameter is of type InitializeParams
	// the response if of type InitializeResult of a Thenable that
	// resolves to such.
	MethodInitialize Method = "initialize"
	// The initialized notification is sent from the client to the
	// server after the client is fully initialized and the server
	// is allowed to send requests from the server to the client.
	MethodInitialized Method = "initialized"
	// A shutdown request is sent from the client to the server.
	// It is sent once when the client decides to shutdown the
	// server. The only notification that is sent after a shutdown request
	// is the exit event.
	MethodShutdown Method = "shutdown"
	// The exit event is sent from the client to the server to
	// ask the server to exit its process.
	MethodExit Method = "exit"
	// The document open notification is sent from the client to the server to signal
	// newly opened text documents. The document's truth is now managed by the client
	// and the server must not try to read the document's truth using the document's
	// uri. Open in this sense means it is managed by the client. It doesn't necessarily
	// mean that its content is presented in an editor. An open notification must not
	// be sent more than once without a corresponding close notification send before.
	// This means open and close notification must be balanced and the max open count
	// is one.
	MethodTextDocumentDidOpen Method = "textDocument/didOpen"
)
