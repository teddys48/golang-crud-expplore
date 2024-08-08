package main

import (
	"net/http"

	"github.com/teddys48/kmpro/config"
	"github.com/teddys48/kmpro/helper"
)

func main() {
	viperConfig := config.NewViper()
	config.NewLogger()
	db := config.NewDatabase(viperConfig)
	validate := config.NewValidator(viperConfig)
	redis := config.NewRedisConfig(viperConfig)
	route := config.NewRoute()

	// mux := http.NewServeMux()

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

	// log.Info("Starting apps...")
	port := viperConfig.GetString("web.port")
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}
