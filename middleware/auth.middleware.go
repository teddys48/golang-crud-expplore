package middleware

import (
	"fmt"
	"net/http"

	"github.com/gookit/slog"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func NewAuthMiddleware(log *slog.Logger, config *viper.Viper, redis *redis.Client, next http.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("lewat dong")
			next.ServeHTTP(w, r)
		})
	}
}
