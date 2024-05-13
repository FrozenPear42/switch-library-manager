package nut

import (
	"fmt"
	"github.com/FrozenPear42/switch-library-manager/data"
	"go.uber.org/zap"
	"net/http"
)

type ProgressReporter interface {
	ReportProgress(filePath string, downloaded, total int64)
}

type Server struct {
	logger         *zap.SugaredLogger
	serverHost     string
	serverPort     int
	libraryManager data.LibraryManager
	reporter       ProgressReporter
	//users      []struct {
	//	Username string
	//	Password string
	//}
}

func NewServer(host string, port int, libraryManager data.LibraryManager, reporter ProgressReporter) *Server {
	return &Server{
		logger:         zap.S(),
		serverHost:     host,
		serverPort:     port,
		libraryManager: libraryManager,
		reporter:       reporter,
	}
}

func (s *Server) Listen() error {
	router := NewRouter(s.libraryManager, s.reporter)

	server := http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.serverHost, s.serverPort),
		Handler: router,
	}
	err := server.ListenAndServe()
	if err != nil {
		return err
	}
	s.logger.Infof("started NUT server at %s:%d", s.serverHost, s.serverPort)

	return nil
}
