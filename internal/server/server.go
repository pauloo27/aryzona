package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/pauloo27/aryzona/internal/config"
	"github.com/pauloo27/logger"
)

func StartHTTPServer() error {
	logger.Infof("Starting HTTP server at port %d...", config.Config.HTTPServerPort)
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
		Addr: fmt.Sprintf(":%d", config.Config.HTTPServerPort),
	}

	return server.ListenAndServe()
}
