package handler

import (
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/teddys48/kmpro/config"
	"github.com/teddys48/kmpro/helper"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	mux := http.NewServeMux()
	godotenv.Load()
	viperConfig := config.NewViper()
	config.NewLogger()
	db := config.NewDatabase(viperConfig)
	validate := config.NewValidator(viperConfig)
	redis := config.NewRedisConfig(viperConfig)
	route := config.NewRoute()

	config.App(&config.AppConfig{
		DB: db,
		// Log:      log,
		Validate: validate,
		Config:   viperConfig,
		Redis:    redis,
		Route:    route,
	})

	route.Get("/", func(w http.ResponseWriter, r *http.Request) {
		helper.ReturnResponse(w, "Welcome!")
	})

	mux.Handle("/", route)

	// mux.ServeHTTP(w, r)

	// log.Info("Starting apps...")
	port := os.Getenv("appPort")
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}
