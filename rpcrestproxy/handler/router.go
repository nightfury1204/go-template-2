package handler

import (
	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

// GetRouter returns all handler router
func GetRouter(brandHandler *BrandHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)

	// enforce cors policy later
	r.Use(cors.AllowAll().Handler)

	r.Mount("/api/v1/brand", brandHandler.GetRouter())

	return r
}
