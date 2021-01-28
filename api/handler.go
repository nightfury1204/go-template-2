package api

import (
	"log"
	"net/http"

	"bitbucket.org/evaly/go-boilerplate/api/middleware"
	"github.com/go-chi/chi"
)

func brandsRouter(ctrl *BrandsController) http.Handler {
	h := chi.NewRouter()
	h.Group(func(r chi.Router) {
		r.Get("/", ctrl.ListBrand)
	})

	return h
}

func workerRouter(ctrl *WorkerController) http.Handler {
	h := chi.NewRouter()
	h.With(middleware.AppKeyChecker("api-key-from-env")).Group(func(r chi.Router) {
		r.Route("/brands", func(r chi.Router) {
			r.Post("/", ctrl.AddNewBrand)
		})
	})
	return h
}

func healthRouter(ctrl *SystemController) http.Handler {
	log.Println("healthRouter")
	h := chi.NewRouter()
	h.Group(func(r chi.Router) {
		r.Get("/api", ctrl.apiCheck)
		r.Get("/worker", ctrl.workerCheck)
	})
	return h
}
