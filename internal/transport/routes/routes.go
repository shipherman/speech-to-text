package routes

import (
	"github.com/shipherman/speech-to-text/internal/handlers"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var Router *chi.Mux

func init() {
	// Router = chi.NewRouter()

	Router.Use(middleware.Logger)
	Router.Post("/stt", handlers.HandleSTT)
	Router.Post("/api/user/register", handlers.HandleRegister)
	Router.Post("/api/user/login", handlers.HandleLogin)
}
