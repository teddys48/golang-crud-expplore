package helper

import (
	"fmt"
	"math/rand"
)

type WebResponse[T any] struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

func GenerateSessionID() string {
	return fmt.Sprint(rand.Int())
}

type contextKey string

var userKey contextKey = "user_id"

func GetContextKey() contextKey {
	return "user_id"
}
