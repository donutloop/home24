package server

import (
	"github.com/donutloop/home24/internal/handler"
	"github.com/donutloop/home24/pkg/webscrapper"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/http/httptest"
)

// Server is the HTTP server.
type Server struct {
	mux        *chi.Mux
	logger     *logrus.Logger
	testserver *httptest.Server
	TestURL    string
}

// NewServer creates a new Server.
func New() *Server {
	s := &Server{
		logger: logrus.New(),
		mux:    chi.NewRouter(),
	}
	return s
}

func (s *Server) InitHandlers() {
	ctx := handler.Context{
		LogError: s.logger.Error,
		LogInfo:  s.logger.Info,
	}

	websiteDataService := handler.WebsiteDataService{
		Webscrapper: webscrapper.New(new(http.Client)),
	}
	domainHandler := handler.WrapGetWebsiteDataHandler(websiteDataService.Extract, ctx)

	s.mux.Method(http.MethodPost, "/websitestats", domainHandler)
}

// Start starts the server.
func (s *Server) Start(addr string, test bool) error {

	if test {
		s.testserver = httptest.NewServer(s.mux)
		s.TestURL = s.testserver.URL
		return nil
	}

	return http.ListenAndServe(addr, s.mux)
}

func (s *Server) Stop(test bool) error {
	if test {
		s.testserver.Close()
		return nil
	}
	return nil
}
