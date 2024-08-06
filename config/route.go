package config

import "github.com/gorilla/mux"

func NewRoute() *mux.Router {
	r := mux.NewRouter()

	return r
}
