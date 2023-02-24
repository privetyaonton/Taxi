package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RipperAcskt/innotaxi/config"
	"go.uber.org/zap"
)

type Server struct {
	httpServer *http.Server
	Log        *zap.Logger
}

func (s *Server) Run(handler http.Handler, cfg *config.Config) error {
	s.httpServer = &http.Server{
		Addr:    cfg.SERVER_HOST,
		Handler: handler,
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) ShutDown() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	s.Log.Info("Shuttig down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("shut down failed: %w", err)
	}

	s.Log.Info("Server exiting.")
	return nil
}
