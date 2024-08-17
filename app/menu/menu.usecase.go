package menu

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
	"github.com/teddys48/kmpro/entity"
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

	slog.Infof("[%+v] [MENU ALL] REQUEST : %+v", session, nil)

	menu := new([]Menu)
	err := u.Repository.All(tx, menu)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		menu = nil
	}

	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [MENU ALL] RESPONSE : %+v", session, err.Error())
		return response
	}

	response = helper.Response("00", "success", menu)
	slog.Warnf("[%+v] [MENU ALL] RESPONSE : %+v", session, response)
	return response
}

func (u *useCase) Find(r *http.Request) *helper.WebResponse[interface{}] {
	response := &helper.WebResponse[interface{}]{}
	session := helper.GenerateRandomString()
	tx := u.DB.WithContext(r.Context())
	request := r.URL.Query().Get("id")

	slog.Infof("[%+v] [MENU FIND] REQUEST : %+v", session, request)

	id, err := strconv.ParseInt(request, 10, 64)
	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [MENU FIND] RESPONSE : %+v", session, err.Error())
		return response
	}

	users := new(Menu)
	err = u.Repository.CheckByID(tx, users, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response = helper.Response("400", "Data not found", nil)
		slog.Warnf("[%+v] [MENU FIND] RESPONSE : %+v", session, response)
		return response
	}

	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [MENU FIND] RESPONSE : %+v", session, err.Error())
		return response
	}

	response = helper.Response("00", "success", users)
	slog.Warnf("[%+v] [MENU FIND] RESPONSE : %+v", session, response)
	return response
}

func (u *useCase) Create(r *http.Request) *helper.WebResponse[interface{}] {
	response := &helper.WebResponse[interface{}]{}

	session := helper.GenerateRandomString()
	tx := u.DB.WithContext(r.Context())
	request := new(MenuCreateRequest)

	err := helper.ValidateRequest(r, u.Validate, request)
	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [MENU CREATE] RESPONSE : %+v", session, err.Error())
		return response
	}

	slog.Infof("[%+v] [MENU CREATE] REQUEST : %+v", session, request)

	getUserID := fmt.Sprint(context.Context.Value(r.Context(), helper.GetContextKey()))
	userID, err := strconv.ParseInt(getUserID, 10, 64)
	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [MENU CREATE] RESPONSE : %+v", session, err.Error())
		return response
	}

	dataInsert := entity.Menu{
		Code:         "MENU" + fmt.Sprint(time.Now().Unix()),
		Name:         request.Name,
		MenuParentID: request.MenuParentID,
		Icon:         request.Icon,
		PathURL:      request.PathURL,
		Sort:         request.Sort,
		HiddenData:   request.HiddenData,
		Description:  request.Description,
		CreatedBy:    userID,
		CreatedOn:    time.Now(),
	}

	err = u.Repository.Create(tx, &dataInsert)
	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [MENU CREATE] RESPONSE : %+v", session, err.Error())
		return response
	}

	response = helper.Response("00", "success", nil)
	slog.Infof("[%+v] [MENU CREATE] RESPONSE : %+v", session, response)
	return response
}

func (u *useCase) Update(r *http.Request) *helper.WebResponse[interface{}] {
	response := &helper.WebResponse[interface{}]{}

	session := helper.GenerateRandomString()
	tx := u.DB.WithContext(r.Context())
	request := new(MenuUpdateRequest)
	requestID := r.URL.Query().Get("id")

	err := helper.ValidateRequest(r, u.Validate, request)
	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [MENU UPDATE] RESPONSE : %+v", session, err.Error())
		return response
	}

	slog.Infof("[%+v] [MENU UPDATE] REQUEST : %+v", session, request)

	id, err := strconv.ParseInt(requestID, 10, 64)
	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [USERS UPDATE] RESPONSE : %+v", session, err.Error())
		return response
	}

	users := new(Menu)
	err = u.Repository.CheckByID(tx, users, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response = helper.Response("400", "Data not found", nil)
		slog.Warnf("[%+v] [MENU UPDATE] RESPONSE : %+v", session, response)
		return response
	}

	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [MENU UPDATE] RESPONSE : %+v", session, err.Error())
		return response
	}

	getUserID := fmt.Sprint(context.Context.Value(r.Context(), helper.GetContextKey()))
	userID, err := strconv.ParseInt(getUserID, 10, 64)
	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [MENU UPDATE] RESPONSE : %+v", session, err.Error())
		return response
	}

	dataUpdate := entity.Menu{
		Name:         request.Name,
		MenuParentID: request.MenuParentID,
		Icon:         request.Icon,
		PathURL:      request.PathURL,
		Sort:         request.Sort,
		HiddenData:   request.HiddenData,
		Description:  request.Description,
		UpdatedBy:    userID,
		UpdatedOn:    time.Now(),
	}

	err = u.Repository.Update(tx, &dataUpdate, id)
	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [MENU UPDATE] RESPONSE : %+v", session, err.Error())
		return response
	}

	response = helper.Response("00", "success", nil)
	slog.Infof("[%+v] [MENU UPDATE] REQUEST : %+v", session, response)
	return response
}

func (u *useCase) Delete(r *http.Request) *helper.WebResponse[interface{}] {
	response := &helper.WebResponse[interface{}]{}

	session := helper.GenerateRandomString()
	tx := u.DB.WithContext(r.Context())
	request := r.URL.Query().Get("id")

	slog.Infof("[%+v] [MENU DELETE] REQUEST : %+v", session, request)

	id, err := strconv.ParseInt(request, 10, 64)
	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [MENU DELETE] RESPONSE : %+v", session, err.Error())
		return response
	}

	users := new(Menu)
	err = u.Repository.CheckByID(tx, users, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		users = nil
	}

	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [MENU DELETE] RESPONSE : %+v", session, err.Error())
		return response
	}

	getUserID := fmt.Sprint(context.Context.Value(r.Context(), helper.GetContextKey()))
	userID, err := strconv.ParseInt(getUserID, 10, 64)
	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [MENU DELETE] RESPONSE : %+v", session, err.Error())
		return response
	}

	dataUpdate := entity.Menu{
		Status:    "DELETED",
		UpdatedBy: userID,
		UpdatedOn: time.Now(),
	}

	err = u.Repository.Delete(tx, &dataUpdate, id)
	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [MENU DELETE] RESPONSE : %+v", session, err.Error())
		return response
	}

	response = helper.Response("00", "success", nil)
	slog.Infof("[%+v] [MENU DELETE] RESPONSE : %+v", session, request)
	return response
}
