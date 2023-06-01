package server

import (
	"fmt"
	"net/http"

	"github.com/pauloo27/aryzona/internal/config"
	"github.com/pauloo27/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

/* #nosec G114 */
func StartHTTPServer() {
	if config.Config.HTTPServerPort == 0 {
		logger.Warn("HTTP server disabled")
		return
	}

	logger.Infof("Starting HTTP server at port %d...", config.Config.HTTPServerPort)
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	route(r)

	err := http.ListenAndServe(fmt.Sprintf(":%d", config.Config.HTTPServerPort), r)
	if err != nil {
		logger.Fatal(err)
	}
}