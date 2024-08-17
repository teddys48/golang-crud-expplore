package corporation

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gookit/slog"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"github.com/teddys48/kmpro/helper"
	"gorm.io/gorm"
)

type UseCase interface {
	Create(r *http.Request) *helper.WebResponse[interface{}]
	Find(r *http.Request) *helper.WebResponse[interface{}]
	All(r *http.Request) *helper.WebResponse[interface{}]
	Update(r *http.Request) *helper.WebResponse[interface{}]
	Delete(r *http.Request) *helper.WebResponse[interface{}]
}

type useCase struct {
	DB         *gorm.DB
	Validate   *validator.Validate
	Repository Repository
	Config     *viper.Viper
	Redis      *redis.Client
}

func NewUseCase(db *gorm.DB, validate *validator.Validate, Repository Repository, viper *viper.Viper, redis *redis.Client) UseCase {
	return &useCase{
		DB:         db,
		Validate:   validate,
		Repository: Repository,
		Config:     viper,
		Redis:      redis,
	}
}

func (u *useCase) All(r *http.Request) *helper.WebResponse[interface{}] {
	response := &helper.WebResponse[interface{}]{}

	session := helper.GenerateRandomString()
	tx := u.DB.WithContext(r.Context())

	slog.Infof("[%+v] [CORPORATION ALL] REQUEST : %+v", session, nil)

	data := new([]Corporation)
	err := u.Repository.All(tx, data)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		data = nil
	}

	response = helper.Response("00", "success", data)
	slog.Infof("[%+v] [CORPORATION ALL] RESPONSE : %+v", session, response)
	return response
}

func (u *useCase) Find(r *http.Request) *helper.WebResponse[interface{}] {
	response := &helper.WebResponse[interface{}]{}
	session := helper.GenerateRandomString()
	tx := u.DB.WithContext(r.Context())
	request := r.URL.Query().Get("id")

	slog.Infof("[%+v] [CORPORATION FIND] REQUEST : %+v", session, request)

	id, err := strconv.ParseInt(request, 10, 64)
	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [CORPORATION FIND] RESPONSE : %+v", session, err.Error())
		return response
	}

	data := new(Corporation)
	err = u.Repository.CheckByID(tx, data, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response = helper.Response("400", "Data not found", nil)
		slog.Warnf("[%+v] [CORPORATION FIND] RESPONSE : %+v", session, response)
		return response
	}

	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [CORPORATION FIND] RESPONSE : %+v", session, err.Error())
		return response
	}

	response = helper.Response("00", "success", data)
	slog.Warnf("[%+v] [CORPORATION FIND] RESPONSE : %+v", session, response)
	return response
}

func (u *useCase) Create(r *http.Request) *helper.WebResponse[interface{}] {
	response := &helper.WebResponse[interface{}]{}
	session := helper.GenerateRandomString()
	tx := u.DB.WithContext(r.Context())
	request := new(Corporation)

	err := helper.ValidateRequest(r, u.Validate, request)
	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [CORPORATION CREATE] RESPONSE : %+v", session, err.Error())
		return response
	}

	slog.Infof("[%+v] [CORPORATION CREATE] REQUEST : %+v", session, request)

	getUserID := fmt.Sprint(context.Context.Value(r.Context(), helper.GetContextKey()))
	userID, err := strconv.ParseInt(getUserID, 10, 64)
	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [CORPORATION CREATE] RESPONSE : %+v", session, err.Error())
		return response
	}

	request.CreatedBy = userID
	request.CreatedOn = time.Now()
	request.Status = "ACTIVE"
	request.Code = "CORPORATION" + fmt.Sprint(time.Now().Unix())

	err = u.Repository.Create(tx, request)
	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [CORPORATION CREATE] RESPONSE : %+v", session, err.Error())
		return response
	}

	response = helper.Response("00", "success", nil)
	slog.Infof("[%+v] [CORPORATION CREATE] RESPONSE : %+v", session, response)
	return response
}

func (u *useCase) Update(r *http.Request) *helper.WebResponse[interface{}] {
	response := &helper.WebResponse[interface{}]{}
	session := helper.GenerateRandomString()
	tx := u.DB.WithContext(r.Context())
	request := new(Corporation)
	requestID := r.URL.Query().Get("id")

	err := helper.ValidateRequest(r, u.Validate, request)
	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [CORPORATION UPDATE] RESPONSE : %+v", session, err.Error())
		return response
	}

	slog.Infof("[%+v] [CORPORATION UPDATE] REQUEST : %+v", session, request)

	id, err := strconv.ParseInt(requestID, 10, 64)
	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [USERS UPDATE] RESPONSE : %+v", session, err.Error())
		return response
	}

	data := new(Corporation)
	err = u.Repository.CheckByID(tx, data, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response = helper.Response("400", "Data not found", nil)
		slog.Warnf("[%+v] [CORPORATION UPDATE] RESPONSE : %+v", session, response)
		return response
	}

	getUserID := fmt.Sprint(context.Context.Value(r.Context(), helper.GetContextKey()))
	userID, err := strconv.ParseInt(getUserID, 10, 64)
	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [CORPORATION UPDATE] RESPONSE : %+v", session, err.Error())
		return response
	}

	now := time.Now()

	request.UpdatedBy = &userID
	request.UpdatedOn = &now

	err = u.Repository.Update(tx, request, id)
	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [CORPORATION UPDATE] RESPONSE : %+v", session, err.Error())
		return response
	}

	response = helper.Response("00", "success", nil)
	slog.Infof("[%+v] [CORPORATION UPDATE] RESPONSE : %+v", session, response)
	return response
}

func (u *useCase) Delete(r *http.Request) *helper.WebResponse[interface{}] {
	response := &helper.WebResponse[interface{}]{}
	session := helper.GenerateRandomString()
	tx := u.DB.WithContext(r.Context())
	request := r.URL.Query().Get("id")

	slog.Infof("[%+v] [CORPORATION DELETE] REQUEST : %+v", session, request)

	getUserID := fmt.Sprint(context.Context.Value(r.Context(), helper.GetContextKey()))
	userID, err := strconv.ParseInt(getUserID, 10, 64)
	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [CORPORATION DELETE] RESPONSE : %+v", session, err.Error())
		return response
	}

	id, err := strconv.ParseInt(request, 10, 64)
	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [CORPORATION DELETE] RESPONSE : %+v", session, err.Error())
		return response
	}

	data := new(Corporation)
	err = u.Repository.CheckByID(tx, data, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response = helper.Response("400", "Data not found", nil)
		slog.Warnf("[%+v] [CORPORATION DELETE] RESPONSE : %+v", session, response)
		return response
	}

	now := time.Now()

	dataUpdate := Corporation{
		UpdatedBy: &userID,
		UpdatedOn: &now,
		Status:    "DELETED",
	}

	err = u.Repository.Update(tx, &dataUpdate, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response = helper.Response("400", "User not found", nil)
		slog.Warnf("[%+v] [CORPORATION DELETE] RESPONSE : %+v", session, response)
		return response
	}

	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [CORPORATION DELETE] RESPONSE : %+v", session, err.Error())
		return response
	}

	response = helper.Response("00", "success", nil)
	slog.Warnf("[%+v] [CORPORATION DELETE] RESPONSE : %+v", session, response)
	return response
}
