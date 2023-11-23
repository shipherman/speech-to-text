package routes

import (
	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var Router *chi.Mux

func init() {
	Router = chi.NewRouter()

	Router.Use(middleware.Logger)
	Router.Get("/", h.HandleRoot)
	Router.Post("/api/user/register", h.HandleRegister)
	Router.Post("/api/user/login", h.HandleLogin)
	Router.Post("/api/transcribe", h.Transcribe)
}
