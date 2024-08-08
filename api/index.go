package handler

import (
	"fmt"
	"net/http"

	"github.com/joho/godotenv"
	. "github.com/tbxark/g4vercel"
	"github.com/teddys48/kmpro/config"
	"github.com/teddys48/kmpro/helper"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	server := New()
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

	// mux.ServeHTTP(w, r)

	// log.Info("Starting apps...")
	// port := os.Getenv("appPort")
	// err := http.ListenAndServe(":"+port, nil)
	// if err != nil {
	// 	panic(err)
	// }

	// server.Use(Recovery(func(err interface{}, c *Context) {
	// 	if httpError, ok := err.(HttpError); ok {
	// 		c.JSON(httpError.Status, H{
	// 			"message": httpError.Error(),
	// 		})
	// 	} else {
	// 		message := fmt.Sprintf("%s", err)
	// 		c.JSON(500, H{
	// 			"message": message,
	// 		})
	// 	}
	// }))
	// server.GET("/", func(context *Context) {
	// 	context.JSON(200, H{
	// 		"message": "OK",
	// 	})
	// })

	server.Handle(w, r)
}
