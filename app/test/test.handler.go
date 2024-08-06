package test

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gookit/slog"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type TestHandler interface {
	TestHandler(w http.ResponseWriter, r *http.Request)
}

type testHandler struct {
	DB          *gorm.DB
	Log         *slog.Logger
	Validate    *validator.Validate
	TestUsecase TestUsecase
	Config      *viper.Viper
	Redis       *redis.Client
}

func NewTestHandler(db *gorm.DB, log *slog.Logger, validate *validator.Validate, config *viper.Viper, redis *redis.Client, TestUsecase TestUsecase) TestHandler {
	return &testHandler{
		DB:          db,
		Log:         log,
		Validate:    validate,
		TestUsecase: TestUsecase,
		Config:      config,
		Redis:       redis,
	}
}

func (u testHandler) TestHandler(w http.ResponseWriter, r *http.Request) {
	res := u.TestUsecase.TestUsecase(r)
	json.NewEncoder(w).Encode(res)
}
