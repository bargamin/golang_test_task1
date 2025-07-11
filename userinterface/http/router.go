package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"

	"golang_test_task1/userinterface/http/handlers"
)

func NewRouter(formHandler *handlers.FormHandler) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Routes
	r.Route("/", func(router chi.Router) {
		router.Get("/", formHandler.Form)
		router.Post("/", formHandler.ProcessForm)
	})

	return r
}
