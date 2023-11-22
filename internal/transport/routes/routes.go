package routes

import (
	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var Router *chi.Mux

func init() {
	Router = chi.NewRouter()

	Router.Use(middleware.Logger)

}
