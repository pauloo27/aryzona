package server

import (
	"github.com/pauloo27/aryzona/internal/server/handler"
	"github.com/go-chi/chi/v5"
)

func route(r *chi.Mux) {
	r.Get("/soccer/banner-{t1}-{t2}.png", handler.RenderBanner)
	r.Get("/healthz", handler.Health)
}
