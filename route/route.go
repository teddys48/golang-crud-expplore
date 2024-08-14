package route

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/teddys48/kmpro/app/auth"
	"github.com/teddys48/kmpro/app/menu"
	"github.com/teddys48/kmpro/app/role"
	"github.com/teddys48/kmpro/app/test"
	"github.com/teddys48/kmpro/app/users"
)

type RouteConfig struct {
	AuthMiddleware func(http.Handler) http.Handler
	Route          *chi.Mux
	TestHandler    test.TestHandler
	AuthHandler    auth.AuthHandler
	UsersHandler   users.Handler
	MenuHandler    menu.Handler
	RoleHandler    role.Handler
}

func (c *RouteConfig) Setup() {
	c.AuthRoute()
	c.GuestRoute()
}

func (c *RouteConfig) AuthRoute() {
	c.Route.With(c.AuthMiddleware).Get("/api/auth/refresh-token", c.AuthHandler.RefreshToken)

	userRoute := c.Route.With(c.AuthMiddleware)
	userRoute.Get("/api/users", c.UsersHandler.All)
	userRoute.Get("/api/users/find", c.UsersHandler.Find)
	userRoute.Post("/api/users/create", c.UsersHandler.Create)
	userRoute.Post("/api/users/update", c.UsersHandler.Update)
	userRoute.Get("/api/users/delete", c.UsersHandler.Delete)

	menuRoute := c.Route.With(c.AuthMiddleware)
	menuRoute.Get("/api/menu", c.MenuHandler.All)
	menuRoute.Get("/api/menu/find", c.MenuHandler.Find)
	menuRoute.Post("/api/menu/create", c.MenuHandler.Create)
	menuRoute.Post("/api/menu/update", c.MenuHandler.Update)
	menuRoute.Get("/api/menu/delete", c.MenuHandler.Delete)

	roleRoute := c.Route.With(c.AuthMiddleware)
	roleRoute.Get("/api/role", c.RoleHandler.All)
	roleRoute.Get("/api/role/find", c.RoleHandler.Find)
	roleRoute.Post("/api/role/create", c.RoleHandler.Create)
	roleRoute.Post("/api/role/update", c.RoleHandler.Update)
	roleRoute.Get("/api/role/delete", c.RoleHandler.Delete)
}

func (c *RouteConfig) GuestRoute() {
	c.Route.Route("/api/auth", func(r chi.Router) {
		r.Post("/login", c.AuthHandler.Login)
	})
}
