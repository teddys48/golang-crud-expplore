package route

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/teddys48/kmpro/app/test"
)

type RouteConfig struct {
	AuthMiddleware func(http.Handler) http.Handler
	Route          *chi.Mux
	TestHandler    test.TestHandler
}

func (c *RouteConfig) Setup() {
	c.AuthRoute()
	c.GuestRoute()
}

func (c *RouteConfig) AuthRoute() {
	// prefix := c.Route.Path("/asas").Subrouter().Use()
	// c.Route.Route("/api", func(r chi.Router) {
	// 	r.Use(c.AuthMiddleware)
	// 	r.Get("/test", c.TestHandler.TestHandler)
	// })
	a := c.Route.With(c.AuthMiddleware)
	a.Get("/test", c.TestHandler.TestHandler)
}

func (c *RouteConfig) GuestRoute() {
	// c.Route.Get("/", func(w http.ResponseWriter, r *http.Request) {
	// 	helper.ReturnResponse(w, "Welcome!")
	// })
}
