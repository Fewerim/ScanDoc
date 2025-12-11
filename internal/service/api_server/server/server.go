package server

import (
	"context"
	"errors"
	"net/http"
	hand "proWeb/internal/service/api_server/handlers"

	"github.com/gorilla/mux"
)

type Server struct {
	httpHandler *hand.HTTPHandler
	server      *http.Server
}

func NewServer(h *hand.HTTPHandler) *Server {
	return &Server{
		httpHandler: h,
	}
}

func StartServer(s *Server) error {
	router := mux.NewRouter()

	router.Path("/process").Methods("POST").HandlerFunc(s.httpHandler.ProcessHandler)
	router.Path("/batch").Methods("POST").HandlerFunc(s.httpHandler.BatchHandler)
	router.Path("/health").Methods("GET").HandlerFunc(s.httpHandler.HealthHandler)

	s.server = &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	if err := s.server.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		} else {
			return err
		}
	}
	return nil
}

func (s *Server) Stop() error {
	if s.server == nil {
		return nil
	}

	return s.server.Shutdown(context.Background())
}
