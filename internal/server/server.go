package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/pauloo27/aryzona/internal/config"
)

func StartHTTPServer() error {
	slog.Info("Starting HTTP server", "port", config.Config.HTTPServerPort)
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	route(r)

	server := &http.Server{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      r,
		Addr:         fmt.Sprintf(":%d", config.Config.HTTPServerPort),
	}

	return server.ListenAndServe()
}
