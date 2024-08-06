package test

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gookit/slog"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"github.com/teddys48/kmpro/helper"
	"gorm.io/gorm"
)

type TestUsecase interface {
	TestUsecase(r *http.Request) *helper.WebResponse[interface{}]
}

type testUsecase struct {
	DB             *gorm.DB
	Log            *slog.Logger
	Validate       *validator.Validate
	TestRepository TestRepository
	Config         *viper.Viper
	Redis          *redis.Client
}

func NewTestUsecase(db *gorm.DB, log *slog.Logger, validate *validator.Validate, config *viper.Viper, redis *redis.Client, TestRepository TestRepository) TestUsecase {
	return &testUsecase{
		DB:             db,
		Log:            log,
		Validate:       validate,
		TestRepository: TestRepository,
		Config:         config,
		Redis:          redis,
	}
}

func (u testUsecase) TestUsecase(r *http.Request) *helper.WebResponse[interface{}] {
	response := new(helper.WebResponse[interface{}])

	response = helper.Response(0, "success", nil)

	return response
}
