package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/teddys48/kmpro/config"
	"github.com/teddys48/kmpro/helper"
)

func StartNonTLSServer(port string, srv *http.Server) {
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Printf("Error: %v\n", err)
	}
	// if err != nil {
	// 	panic(err)
	// } else {
	// 	slog.Infof("Starting server on port %v", port)
	// 	fmt.Printf("Starting server on port %v", port)
	// }
}

func StartTLSServer(port string, srv *http.Server) {
	if err := srv.ListenAndServeTLS("server.crt", "server.key"); err != nil && err != http.ErrServerClosed {
		fmt.Printf("Error: %v\n", err)
	}
	// if err != nil {
	// 	panic(err)
	// }
}

func main() {
	// go StartNonTLSServer()
	// godotenv.Load()
	// mux := http.NewServeMux()
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

	srv := &http.Server{
		Addr:    ":8080",
		Handler: nil,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// log.Info("Starting apps...")
	port := viperConfig.GetString("web.port")
	if viperConfig.GetString("app.env") == "development" {
		StartNonTLSServer(port, srv)
	} else {
		StartTLSServer(port, srv)
	}

	<-stop
	fmt.Println("Shutting down server...")

	// Buat context dengan timeout untuk graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Tutup server dengan graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Error during shutdown: %v\n", err)
	}

	fmt.Println("Server stopped gracefully")
	// err := http.ListenAndServeTLS(":"+port, "server.crt", "server.key", nil)
	// if err != nil {
	// 	panic(err)
	// }
}
