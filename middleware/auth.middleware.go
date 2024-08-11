package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gookit/slog"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"github.com/teddys48/kmpro/helper"
)

type ClaimsToken struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type ClaimsData struct {
	UserID string `json:"user_id"`
}

type contextKey string

var userKey contextKey = "user_id"

func NewAuthMiddleware(log *slog.Logger, config *viper.Viper, redis *redis.Client, next http.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			response := &helper.WebResponse[interface{}]{}

			getToken := r.Header.Get("Authorization")

			key, err := helper.PublicKey()
			if err != nil {
				response = helper.Response("500", err.Error(), nil)

				log.Warnf("[RESPONSE MIDDLEWARE-AUTH]%+v", err.Error())

				helper.ReturnResponse(w, response)
			}

			if getToken == "" {
				response = helper.Response("401", "token not found", nil)

				log.Warnf("[RESPONSE MIDDLEWARE-AUTH] %+v", "token not found")

				helper.ReturnResponse(w, response)
			}

			token := getToken[7:]

			tok, err := jwt.ParseWithClaims(token, &ClaimsToken{}, func(jwtToken *jwt.Token) (interface{}, error) {
				if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
					log.Warnf("[RESPONSE MIDDLEWARE-AUTH] - failed parse token: %+v", fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"]))
					return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
				}

				return key, nil
			})

			if err != nil {
				response = helper.Response("402", err.Error(), nil)

				log.Warnf("[RESPONSE MIDDLEWARE-AUTH] %+v", err.Error())

				helper.ReturnResponse(w, response)
			}

			claims, ok := tok.Claims.(*ClaimsToken)
			if !ok || !tok.Valid {
				response = helper.Response("403", "invalid token", nil)

				log.Warnf("[RESPONSE MIDDLEWARE-AUTH] %+v", "invalid token")

				helper.ReturnResponse(w, response)
			}

			decryptKey := config.GetString("encrypt.key")

			userID, err := helper.Decrypt(claims.UserID, []byte(decryptKey))
			if err != nil {
				response = helper.Response("500", err.Error(), nil)

				log.Warnf("[RESPONSE MIDDLEWARE-AUTH] %+v", err.Error())

				helper.ReturnResponse(w, response)
			}

			fmt.Println("lewat dong")
			ctx := context.WithValue(r.Context(), userKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func NewRefreshTokenMiddleware(log *slog.Logger, config *viper.Viper, redis *redis.Client, next http.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			response := &helper.WebResponse[interface{}]{}

			token := r.Header.Get("Authorization")[7:]

			key, err := helper.PublicKey()
			if err != nil {
				response = helper.Response("500", err.Error(), nil)

				log.Warnf("[RESPONSE MIDDLEWARE-AUTH]%+v", err.Error())

				helper.ReturnResponse(w, response)
			}

			if token == "" {
				response = helper.Response("401", "token not found", nil)

				log.Warnf("[RESPONSE MIDDLEWARE-AUTH] %+v", "token not found")

				helper.ReturnResponse(w, response)
			}

			tok, err := jwt.ParseWithClaims(token, &ClaimsToken{}, func(jwtToken *jwt.Token) (interface{}, error) {
				if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
					log.Warnf("[RESPONSE MIDDLEWARE-AUTH] - failed parse token: %+v", fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"]))
					return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
				}

				return key, nil
			})

			if err != nil {
				response = helper.Response("402", err.Error(), nil)

				log.Warnf("[RESPONSE MIDDLEWARE-AUTH] %+v", err.Error())

				helper.ReturnResponse(w, response)
			}

			claims, ok := tok.Claims.(*ClaimsToken)
			if !ok || !tok.Valid {
				response = helper.Response("403", "invalid token", nil)

				log.Warnf("[RESPONSE MIDDLEWARE-AUTH] %+v", "invalid token")

				helper.ReturnResponse(w, response)
			}

			decryptKey := config.GetString("encrypt.key")

			userID, err := helper.Decrypt(claims.UserID, []byte(decryptKey))
			if err != nil {
				response = helper.Response("500", err.Error(), nil)

				log.Warnf("[RESPONSE MIDDLEWARE-AUTH] %+v", err.Error())

				helper.ReturnResponse(w, response)
			}

			fmt.Println("lewat dong")
			ctx := context.WithValue(r.Context(), userKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
