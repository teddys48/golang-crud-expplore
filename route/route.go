package route

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/teddys48/kmpro/app/auth"
	"github.com/teddys48/kmpro/app/test"
)

type RouteConfig struct {
	AuthMiddleware func(http.Handler) http.Handler
	Route          *chi.Mux
	TestHandler    test.TestHandler
	AuthHandler    auth.AuthHandler
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
	// c.Route.Route("/api/", func(r chi.Router) {
	// 	r.Use(c.AuthMiddleware)
	// 	r.Post("/auth/refresh-token",)
	// })
	// a := c.Route.With(c.AuthMiddleware)
	// a.Get("/test", c.TestHandler.TestHandler)
	c.Route.With(c.AuthMiddleware).Post("/api/auth/refresh-token", c.AuthHandler.RefreshToken)
}

func (c *RouteConfig) GuestRoute() {
	c.Route.Route("/api/auth", func(r chi.Router) {
		r.Post("/login", c.AuthHandler.Login)
	})
}
