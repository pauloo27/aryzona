package server

import (
	"fmt"
	"net/http"

	"github.com/Pauloo27/aryzona/internal/config"
	"github.com/Pauloo27/logger"
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

	r.Get("/soccer/banner-{t1}-{t2}.png", renderBanner)

	err := http.ListenAndServe(fmt.Sprintf(":%d", config.Config.HTTPServerPort), r)
	if err != nil {
		logger.Fatal(err)
	}
}
