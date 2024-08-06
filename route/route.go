package route

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/teddys48/kmpro/app/test"
)

type RouteConfig struct {
	Route          *mux.Router
	AuthMiddleware http.Handler
	TestHandler    test.TestHandler
}
