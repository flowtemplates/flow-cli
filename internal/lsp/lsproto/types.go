package lsproto

import (
	"encoding/json"
	"fmt"
)

type TraceValue string

// Predefined Language kinds
// Since: 3.18.0
type LanguageKind string

// A literal to identify a text document in the client.
type TextDocumentIdentifier struct {
	// The text document's uri.
	Uri DocumentURI `json:"uri"`
}

// A set of predefined position encoding kinds.
//
// Since: 3.17.0
type PositionEncodingKind string

type ClientInfo struct {
	// The name of the client as defined by the client.
	Name string `json:"name"`

	// The client's version as defined by the client.
	Version *string `json:"version,omitempty"`
}

// An item to transfer a text document from the client to the
// server.
type TextDocumentItem struct {
	// The text document's uri.
	Uri DocumentURI `json:"uri"`

	// The text document's language identifier.
	LanguageId LanguageKind `json:"languageId"`

	// The version number of this document (it will increase after each
	// change, including undo/redo).
	Version int32 `json:"version"`

	// The content of the opened text document.
	Text string `json:"text"`
}

// The parameters sent in an open text document notification
type DidOpenTextDocumentParams struct {
	// The document that was opened.
	TextDocument TextDocumentItem `json:"textDocument"`
}

// Since: 3.16.0
type SemanticTokensParams struct {
	// The text document.
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

// Since: 3.16.0
type SemanticTokens struct {
	// An optional result id. If provided and clients support delta updating
	// the client will include the result id in the next semantic token request.
	// A server can then instead of computing all semantic tokens again simply
	// send a delta.
	ResultId *string `json:"resultId,omitempty"`

	// The actual tokens.
	Data []uint32 `json:"data"`
}

// The initialize parameters
type InitializeParams struct {
	// WorkDoneProgressParams

	// The process Id of the parent process that started
	// the server.
	//
	// Is `null` if the process has not been started by another process.
	// If the parent process is not alive then the server should exit.
	ProcessId *int32 `json:"processId"`

	// Information about the client
	//
	// Since: 3.15.0
	ClientInfo *ClientInfo `json:"clientInfo,omitempty"`

	// The locale the client is currently showing the user interface
	// in. This must not necessarily be the locale of the operating
	// system.
	//
	// Uses IETF language tags as the value's syntax
	// (See https://en.wikipedia.org/wiki/IETF_language_tag)
	//
	// Since: 3.16.0
	Locale *string `json:"locale,omitempty"`

	// The rootPath of the workspace. Is null
	// if no folder is open.
	//
	// Deprecated: in favour of rootUri.
	RootPath *Nullable[string] `json:"rootPath,omitempty"`

	// The rootUri of the workspace. Is null if no
	// folder is open. If both `rootPath` and `rootUri` are set
	// `rootUri` wins.
	//
	// Deprecated: in favour of workspaceFolders.
	RootUri *DocumentURI `json:"rootUri"`

	// The capabilities provided by the client (editor or tool)
	// Capabilities ClientCapabilities `json:"capabilities"`

	// User provided initialization options.
	InitializationOptions *any `json:"initializationOptions,omitempty"`

	// The initial trace setting. If omitted trace is disabled ('off').
	Trace *TraceValue `json:"trace,omitempty"`
}

type ServerInfo struct {
	// The name of the server as defined by the server.
	Name string `json:"name"`

	// The server's version as defined by the server.
	Version *string `json:"version,omitempty"`
}

type WorkDoneProgressOptions struct {
	WorkDoneProgress *bool `json:"workDoneProgress,omitempty"`
}

// Since: 3.16.0
type SemanticTokensLegend struct {
	// The token types a server uses.
	TokenTypes []SemanticTokenType `json:"tokenTypes"`

	// The token modifiers a server uses.
	TokenModifiers []string `json:"tokenModifiers"`
}

// Semantic tokens options to support deltas for full documents
//
// Since: 3.18.0
type SemanticTokensFullDelta struct {
	// The server supports deltas for full documents.
	Delta *bool `json:"delta,omitempty"`
}

type BooleanOrEmptyObject struct {
	Boolean     *bool
	EmptyObject *struct{}
}

type BooleanOrSemanticTokensFullDelta struct {
	Boolean                 *bool
	SemanticTokensFullDelta *SemanticTokensFullDelta
}

// Defines how the host (editor) should sync
// document changes to the language server.
type TextDocumentSyncKind uint32

const (
	// Documents should not be synced at all.
	TextDocumentSyncKindNone TextDocumentSyncKind = 0
	// Documents are synced by always sending the full content
	// of the document.
	TextDocumentSyncKindFull TextDocumentSyncKind = 1
	// Documents are synced by sending the full content on open.
	// After that only incremental updates to the document are
	// send.
	TextDocumentSyncKindIncremental TextDocumentSyncKind = 2
)

func (e *TextDocumentSyncKind) UnmarshalJSON(data []byte) error {
	var v uint32
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case 0, 1, 2:
		*e = TextDocumentSyncKind(v)
		return nil
	default:
		return fmt.Errorf("unknown TextDocumentSyncKind value: %v", v)
	}
}

type SaveOptions struct {
	// The client is supposed to include the content on save.
	IncludeText *bool `json:"includeText,omitempty"`
}

type BooleanOrSaveOptions struct {
	Boolean     *bool
	SaveOptions *SaveOptions
}

type TextDocumentSyncOptions struct {
	// Open and close notifications are sent to the server. If omitted open close notification should not
	// be sent.
	OpenClose *bool `json:"openClose,omitempty"`

	// Change notifications are sent to the server. See TextDocumentSyncKind.None, TextDocumentSyncKind.Full
	// and TextDocumentSyncKind.Incremental. If omitted it defaults to TextDocumentSyncKind.None.
	Change *TextDocumentSyncKind `json:"change,omitempty"`

	// If present will save notifications are sent to the server. If omitted the notification should not be
	// sent.
	WillSave *bool `json:"willSave,omitempty"`

	// If present will save wait until requests are sent to the server. If omitted the request should not be
	// sent.
	WillSaveWaitUntil *bool `json:"willSaveWaitUntil,omitempty"`

	// If present save notifications are sent to the server. If omitted the notification should not be
	// sent.
	Save *BooleanOrSaveOptions `json:"save,omitempty"`
}

// Since: 3.16.0
type SemanticTokensOptions struct {
	WorkDoneProgressOptions

	// The legend used by the server
	Legend SemanticTokensLegend `json:"legend"`

	// Server supports providing semantic tokens for a specific range
	// of a document.
	Range *BooleanOrEmptyObject `json:"range,omitempty"`

	// Server supports providing semantic tokens for a full document.
	Full *BooleanOrSemanticTokensFullDelta `json:"full,omitempty"`
}

type TextDocumentSyncOptionsOrTextDocumentSyncKind struct {
	TextDocumentSyncOptions *TextDocumentSyncOptions
	TextDocumentSyncKind    *TextDocumentSyncKind
}

type ServerCapabilities struct {
	// The position encoding the server picked from the encodings offered
	// by the client via the client capability `general.positionEncodings`.
	//
	// If the client didn't provide any position encodings the only valid
	// value that a server can return is 'utf-16'.
	//
	// If omitted it defaults to 'utf-16'.
	//
	// Since: 3.17.0
	PositionEncoding *PositionEncodingKind `json:"positionEncoding,omitempty"`

	// Defines how text documents are synced. Is either a detailed structure
	// defining each notification or for backwards compatibility the
	// TextDocumentSyncKind number.
	TextDocumentSync *TextDocumentSyncOptionsOrTextDocumentSyncKind `json:"textDocumentSync,omitempty"`

	// The server provides semantic tokens support.
	//
	// Since: 3.16.0
	SemanticTokensProvider *SemanticTokensOptions `json:"semanticTokensProvider,omitempty"`
}

type InitializeResult struct {
	// The capabilities the language server provides.
	Capabilities ServerCapabilities `json:"capabilities"`

	// Information about the server.
	//
	// Since: 3.15.0
	ServerInfo *ServerInfo `json:"serverInfo,omitempty"`
}

// A set of predefined token types. This set is not fixed
// an clients can specify additional token types via the
// corresponding client capabilities.
//
// Since: 3.16.0
type SemanticTokenType string

const (
	SemanticTokenTypesnamespace     SemanticTokenType = "namespace"
	SemanticTokenTypestype          SemanticTokenType = "type"
	SemanticTokenTypesclass         SemanticTokenType = "class"
	SemanticTokenTypesenum          SemanticTokenType = "enum"
	SemanticTokenTypesinterface     SemanticTokenType = "interface"
	SemanticTokenTypesstruct        SemanticTokenType = "struct"
	SemanticTokenTypestypeParameter SemanticTokenType = "typeParameter"
	SemanticTokenTypesparameter     SemanticTokenType = "parameter"
	SemanticTokenTypesvariable      SemanticTokenType = "variable"
	SemanticTokenTypesproperty      SemanticTokenType = "property"
	SemanticTokenTypesenumMember    SemanticTokenType = "enumMember"
	SemanticTokenTypesevent         SemanticTokenType = "event"
	SemanticTokenTypesfunction      SemanticTokenType = "function"
	SemanticTokenTypesmethod        SemanticTokenType = "method"
	SemanticTokenTypesmacro         SemanticTokenType = "macro"
	SemanticTokenTypeskeyword       SemanticTokenType = "keyword"
	SemanticTokenTypesmodifier      SemanticTokenType = "modifier"
	SemanticTokenTypescomment       SemanticTokenType = "comment"
	SemanticTokenTypesstring        SemanticTokenType = "string"
	SemanticTokenTypesnumber        SemanticTokenType = "number"
	SemanticTokenTypesregexp        SemanticTokenType = "regexp"
	SemanticTokenTypesoperator      SemanticTokenType = "operator"
	// Since: 3.17.0
	SemanticTokenTypesdecorator SemanticTokenType = "decorator"
	// Since: 3.18.0
	SemanticTokenTypeslabel SemanticTokenType = "label"
)

func (e *SemanticTokenType) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "class", "comment", "decorator",
		"enum", "enumMember", "event",
		"function", "interface", "keyword",
		"label", "macro", "method",
		"modifier", "namespace", "number",
		"operator", "parameter", "property",
		"regexp", "string", "struct",
		"type", "typeParameter", "variable":
		*e = SemanticTokenType(v)
		return nil
	default:
		return fmt.Errorf("unknown SemanticTokenTypes value: %v", v)
	}
}

// A set of predefined token modifiers. This set is not fixed
// an clients can specify additional token types via the
// corresponding client capabilities.
//
// Since: 3.16.0
type SemanticTokenModifiers string

const (
	SemanticTokenModifiersdeclaration    SemanticTokenModifiers = "declaration"
	SemanticTokenModifiersdefinition     SemanticTokenModifiers = "definition"
	SemanticTokenModifiersreadonly       SemanticTokenModifiers = "readonly"
	SemanticTokenModifiersstatic         SemanticTokenModifiers = "static"
	SemanticTokenModifiersdeprecated     SemanticTokenModifiers = "deprecated"
	SemanticTokenModifiersabstract       SemanticTokenModifiers = "abstract"
	SemanticTokenModifiersasync          SemanticTokenModifiers = "async"
	SemanticTokenModifiersmodification   SemanticTokenModifiers = "modification"
	SemanticTokenModifiersdocumentation  SemanticTokenModifiers = "documentation"
	SemanticTokenModifiersdefaultLibrary SemanticTokenModifiers = "defaultLibrary"
)

func (e *SemanticTokenModifiers) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "abstract", "async", "declaration",
		"defaultLibrary", "definition", "deprecated",
		"documentation", "modification", "readonly",
		"static":
		*e = SemanticTokenModifiers(v)
		return nil
	default:
		return fmt.Errorf("unknown SemanticTokenModifiers value: %v", v)
	}
}
