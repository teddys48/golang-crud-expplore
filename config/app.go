package config

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/gookit/slog"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"github.com/teddys48/kmpro/app/auth"
	"github.com/teddys48/kmpro/app/test"
	"github.com/teddys48/kmpro/app/users"
	"github.com/teddys48/kmpro/middleware"
	"github.com/teddys48/kmpro/route"
	"gorm.io/gorm"
)

type AppConfig struct {
	DB       *gorm.DB
	App      *http.Client
	Log      *slog.Logger
	Validate *validator.Validate
	Config   *viper.Viper
	Redis    *redis.Client
	Route    *chi.Mux
}

func App(config *AppConfig) {
	// config.Route.Use(func(h http.Handler) http.Handler {
	// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		// w.Header().Set("Content-Type", "application/json")
	// 		r.Header.Add("Content-Type", "application/json")
	// 	})
	// })
	config.Route.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// w.Header().Set("Content-Type", "application/json")
			if r.Method == "OPTIONS" {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Methods", "*")
				w.Header().Set("Access-Control-Allow-Headers", "*")
				w.WriteHeader(200)
				h.ServeHTTP(w, r)
			} else {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Methods", "*")
				w.Header().Set("Access-Control-Allow-Headers", "*")
				h.ServeHTTP(w, r)
			}
		})
	})
	authMiddleware := middleware.NewAuthMiddleware(config.Log, config.Config, config.Redis, config.Route)
	config.Route.Use(chiMiddleware.Logger)

	// config.Route.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	// w.Header().Set("Content-Type", "application/json")
	// 	// json.NewEncoder(w).Encode("Welcome!")
	// 	helper.ReturnResponse(w, "Welcome!")
	// })

	testRepo := test.NewTestRepository(config.Log)
	testUsecase := test.NewTestUsecase(config.DB, config.Log, config.Validate, config.Config, config.Redis, testRepo)
	testHandler := test.NewTestHandler(testUsecase)

	authRepo := auth.NewAuthRepository(config.Config)
	authUsecase := auth.NewAuthUseCase(config.DB, config.Validate, authRepo, config.Config, config.Redis)
	authHandler := auth.NewAuthHandler(authUsecase)

	userRepo := users.Newrepository(config.Config)
	userUsecase := users.NewUseCase(config.DB, config.Validate, userRepo, config.Config, config.Redis)
	userHandler := users.NewHandler(userUsecase)

	routeConfig := route.RouteConfig{
		AuthMiddleware: authMiddleware,
		Route:          config.Route,
		TestHandler:    testHandler,
		AuthHandler:    authHandler,
		UsersHandler:   userHandler,
	}

	routeConfig.Setup()

	// config.Route.HandleFunc("*", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("what are you looking for"))
	// })

	http.Handle("/", config.Route)
	http.Handle("*", routeConfig.Route)
}
