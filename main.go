package main

import (
	"fmt"
	"net/http"

	"github.com/gookit/slog"
	"github.com/teddys48/kmpro/config"
	"github.com/teddys48/kmpro/helper"
)

func StartNonTLSServer(port string) {
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	} else {
		slog.Infof("Starting server on port %v", port)
		fmt.Printf("Starting server on port %v", port)
	}
}

func StartTLSServer(port string) {
	err := http.ListenAndServeTLS(":"+port, "server.crt", "server.key", nil)
	if err != nil {
		panic(err)
	}
}

func main() {
	// go StartNonTLSServer()
	// godotenv.Load()
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
		helper.ReturnResponse(w, "Welcomee!")
	})

	// log.Info("Starting apps...")
	port := viperConfig.GetString("web.port")
	if viperConfig.GetString("app.env") == "development" {
		StartNonTLSServer(port)
	} else {
		StartTLSServer(port)
	}
	// err := http.ListenAndServeTLS(":"+port, "server.crt", "server.key", nil)
	// if err != nil {
	// 	panic(err)
	// }
}
