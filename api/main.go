package main

import (
	"net/http"

	"github.com/teddys48/kmpro/config"
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

	// log.Info("Starting apps...")
	err := http.ListenAndServe(":7000", nil)
	if err != nil {
		panic(err)
	}
}
