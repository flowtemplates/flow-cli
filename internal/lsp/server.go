package lsp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"

	"github.com/flowtemplates/flow-cli/internal/lsp/lsproto"
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

	// initializeParams *lsproto.InitializeParams

	// logger         *project.Logger
	// projectService *project.Service
	// converters     *converters
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

		s.logger.DebugContext(ctx, "request recieved", req)

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
