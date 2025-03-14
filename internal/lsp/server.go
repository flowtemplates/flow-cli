package lsp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/url"
	"os"

	"github.com/flowtemplates/flow-cli/internal/lsp/lsproto"
	"github.com/flowtemplates/flow-go/lexer"
	"github.com/flowtemplates/flow-go/token"
)

type Server struct {
	r      *lsproto.BaseReader
	w      *lsproto.BaseWriter
	stderr io.Writer

	logger *slog.Logger
	// requestTime   time.Time

	// cwd                string
	// newLine            core.NewLineKind
	// fs                 vfs.FS
	// defaultLibraryPath string

	initializeParams *lsproto.InitializeParams

	// logger         *project.Logger
	// projectService *project.Service
	// converters     *converters
}

func ptrTo[T any](value T) *T {
	return &value
}

type ServerOptions struct {
	In     io.Reader
	Out    io.Writer
	Err    io.Writer
	Logger *slog.Logger
}

func NewServer(opts *ServerOptions) *Server {
	return &Server{
		r:      lsproto.NewBaseReader(opts.In),
		w:      lsproto.NewBaseWriter(opts.Out),
		stderr: opts.Err,
		logger: opts.Logger,
	}
}

var (
	// These slices must match the order of the indices in the above const block.
	semanticTypeLegend = []lsproto.SemanticTokenType{
		lsproto.SemanticTokenTypeskeyword,
		lsproto.SemanticTokenTypesvariable,
	}
	semanticModifierLegend = []string{}
)

func openFileFromURI(uri string) ([]byte, error) {
	// Parse the URI
	parsedURL, err := url.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URI: %w", err)
	}

	// Extract the file path from the URI
	filePath := parsedURL.Path
	if parsedURL.Scheme != "file" {
		return nil, fmt.Errorf("unsupported scheme: %s", parsedURL.Scheme)
	}

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Read and print the contents of the file
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return content, nil
}

func tokensToIntEncoding(tokens []token.Token) []uint32 {
	var (
		encoded           []uint32
		prevLine, prevCol uint32
	)

	for _, symbol := range tokens {
		var semanticType uint32

		if symbol.IsOneOfMany(token.IF, token.CASE) {
			semanticType = 0
		} else if symbol.IsOneOfMany(token.COMM_TEXT, token.RCOMM, token.LCOMM) {
			semanticType = 1
		} else {
			continue
		}

		// Some fairly painful encoding.
		newLine := uint32(symbol.Pos.Line - 1)
		var newCol uint32
		newCol = uint32(symbol.Pos.Column - 1)
		if prevLine == newLine {
			newCol -= prevCol
		}

		encoded = append(encoded, newLine-prevLine, newCol, uint32(len(symbol.Val)), semanticType, 0)
		prevLine = newLine
		prevCol = uint32(symbol.Pos.Column - 1)
	}

	return encoded
}

func (s *Server) Run(ctx context.Context) error {
	for {
		req, err := s.read()
		if err != nil {
			// if errors.Is(err, lsproto.ErrInvalidRequest) {
			// if err := s.sendError(nil, err); err != nil {
			// 	return err
			// }
			// continue
			// }
			return err
		}

		s.logger.DebugContext(ctx, "request recieved", "body", req)

		switch req.Method {
		case lsproto.MethodInitialize:
			if s.initializeParams != nil {
				return errors.New("initiialize error")
			}

			if err := s.sendResult(req.ID, lsproto.InitializeResult{
				Capabilities: lsproto.ServerCapabilities{
					TextDocumentSync: &lsproto.TextDocumentSyncOptionsOrTextDocumentSyncKind{
						TextDocumentSyncOptions: &lsproto.TextDocumentSyncOptions{
							OpenClose: ptrTo(true),
							Change:    ptrTo(lsproto.TextDocumentSyncKindIncremental),
							Save: &lsproto.BooleanOrSaveOptions{
								SaveOptions: &lsproto.SaveOptions{
									IncludeText: ptrTo(true),
								},
							},
						},
					},
					SemanticTokensProvider: &lsproto.SemanticTokensOptions{
						WorkDoneProgressOptions: lsproto.WorkDoneProgressOptions{
							WorkDoneProgress: ptrTo(false),
						},
						Legend: lsproto.SemanticTokensLegend{
							TokenTypes:     semanticTypeLegend,
							TokenModifiers: semanticModifierLegend,
						},
						Range: &lsproto.BooleanOrEmptyObject{
							Boolean: ptrTo(false),
						},
						Full: &lsproto.BooleanOrSemanticTokensFullDelta{
							Boolean: ptrTo(true),
						},
					},
				},
				ServerInfo: &lsproto.ServerInfo{
					Name:    "flow-langserver",
					Version: ptrTo("1.0"),
				},
			}); err != nil {
				return err
			}
		case lsproto.MethodInitialized:
			s.logger.DebugContext(ctx, "initialized")
		case lsproto.MethodTextDocumentDidOpen:
		case lsproto.MethodTextDocumentSemanticTokensFull:
			var body lsproto.SemanticTokensParams
			if err := json.Unmarshal(req.Params, &body); err != nil {
				return fmt.Errorf("failed to parse: %w", err)
			}

			source, err := openFileFromURI(string(body.TextDocument.Uri))
			if err != nil {
				return err
			}

			tokens := lexer.TokensFromString(string(source))
			s.logger.DebugContext(ctx, "tokens", tokens)

			encoded := tokensToIntEncoding(tokens)

			if err := s.sendResult(req.ID, lsproto.SemanticTokens{
				Data: encoded,
			}); err != nil {
				return err
			}
		default:
			s.logger.ErrorContext(ctx, "unknown method", req.Method)
		}

		// if s.initializeParams == nil {
		// 	if req.Method == lsproto.MethodInitialize {
		// 		if err := s.handleInitialize(req); err != nil {
		// 			return err
		// 		}
		// 	} else {
		// 		if err := s.sendError(req.ID, lsproto.ErrServerNotInitialized); err != nil {
		// 			return err
		// 		}
		// 	}
		// 	continue
		// }

		// if err := s.handleMessage(req); err != nil {
		// 	return err
		// }
	}
}

func (s *Server) sendResult(id *lsproto.ID, result any) error {
	return s.sendResponse(&lsproto.ResponseMessage{
		ID:     id,
		Result: result,
	})
}

func (s *Server) sendResponse(resp *lsproto.ResponseMessage) error {
	data, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	s.logger.Debug("response", "body", resp)
	return s.w.Write(data)
}

func (s *Server) read() (*lsproto.RequestMessage, error) {
	data, err := s.r.Read()
	if err != nil {
		return nil, err
	}

	req := &lsproto.RequestMessage{}
	if err := json.Unmarshal(data, req); err != nil {
		return nil, fmt.Errorf("%w: %w", lsproto.ErrInvalidRequest, err)
	}

	return req, nil
}
