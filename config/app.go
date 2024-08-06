package config

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gookit/slog"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type AppConfig struct {
	DB       *gorm.DB
	App      *http.Client
	Log      *slog.Logger
	Validate *validator.Validate
	Config   *viper.Viper
	Redis    *redis.Client
	Route    *mux.Router
}

func App(config *AppConfig) {
	config.Route.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
		})
	})

	config.Route.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome!"))
	}).Methods("GET")

	config.Route.HandleFunc("*", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("what are you looking for"))
	})

	http.Handle("/", config.Route)
}
