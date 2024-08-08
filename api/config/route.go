package config

import (
	"github.com/go-chi/chi/v5"
)

func NewRoute() *chi.Mux {
	r := chi.NewRouter()

	return r
}
