package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/teddys48/kmpro/config"
	"github.com/teddys48/kmpro/helper"
)

func StartNonTLSServer() {
	mux := new(http.ServeMux)
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Redirecting to https://localhost/")
		http.Redirect(w, r, "https://localhost/", http.StatusTemporaryRedirect)
	}))

	http.ListenAndServe(":8070", mux)
}

func main() {
	// go StartNonTLSServer()
	godotenv.Load()
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
	err := http.ListenAndServeTLS(":"+port, "server.crt", "server.key", nil)
	if err != nil {
		panic(err)
	}
}
