package health

import (
	"github.com/go-chi/chi/v5"

	"github.com/thang1834/go-goss/internal/middleware"
)

func RegisterHTTPEndPoints(router *chi.Mux, uc UseCase) *Handler {
	h := NewHandler(uc)

	router.Route("/api/health", func(router chi.Router) {
		router.Use(middleware.Json)

		router.Get("/", h.Health)
		router.Get("/readiness", h.Readiness)
	})

	return h
}
