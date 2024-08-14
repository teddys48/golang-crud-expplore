package users

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
	"golang.org/x/crypto/bcrypt"
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

	slog.Infof("[%+v] [USERS ALL] REQUEST : %+v", session, nil)

	users := new([]Users)
	err := u.Repository.All(tx, users)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		users = nil
	}

	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [USERS ALL] RESPONSE : %+v", session, err.Error())
		return response
	}

	response = helper.Response("00", "success", users)
	slog.Warnf("[%+v] [USERS ALL] RESPONSE : %+v", session, response)
	return response
}

func (u *useCase) Find(r *http.Request) *helper.WebResponse[interface{}] {
	response := &helper.WebResponse[interface{}]{}
	session := helper.GenerateRandomString()
	tx := u.DB.WithContext(r.Context())
	request := r.URL.Query().Get("id")

	slog.Infof("[%+v] [USERS FIND] REQUEST : %+v", session, request)

	id, err := strconv.ParseInt(request, 10, 64)
	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [USERS FIND] RESPONSE : %+v", session, err.Error())
		return response
	}

	users := new(Users)
	err = u.Repository.CheckUsersByID(tx, users, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response = helper.Response("400", "User not found", nil)
		slog.Warnf("[%+v] [USERS FIND] RESPONSE : %+v", session, response)
		return response
	}

	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [USERS FIND] RESPONSE : %+v", session, err.Error())
		return response
	}

	response = helper.Response("00", "success", users)
	slog.Warnf("[%+v] [USERS FIND] RESPONSE : %+v", session, response)
	return response
}

func (u *useCase) Create(r *http.Request) *helper.WebResponse[interface{}] {
	response := &helper.WebResponse[interface{}]{}
	session := helper.GenerateRandomString()
	tx := u.DB.WithContext(r.Context())
	request := new(UserCreateRequest)

	err := helper.ValidateRequest(r, u.Validate, request)
	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [USERS CREATE] RESPONSE : %+v", session, err.Error())
		return response
	}

	slog.Infof("[%+v] [USERS CREATE] REQUEST : %+v", session, request)

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), 12)
	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [USERS CREATE] RESPONSE : %+v", session, err.Error())
		return response
	}
	fmt.Println(context.Context.Value(r.Context(), helper.GetContextKey()))
	getUserID := fmt.Sprint(context.Context.Value(r.Context(), helper.GetContextKey()))
	userID, err := strconv.ParseInt(getUserID, 10, 64)
	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [USERS CREATE] RESPONSE : %+v", session, err.Error())
		return response
	}

	perusahaanID := int64(*request.PerusahaanID)
	dataInsert := entity.Users{
		Code:          "USER" + fmt.Sprint(time.Now().Unix()),
		Username:      request.Username,
		Nip:           request.NIP,
		CorporationID: &perusahaanID,
		NRK:           request.NRK,
		Instance:      request.Instansi,
		Position:      request.Jabatan,
		FullName:      *request.Fullname,
		Password:      string(hashPassword),
		RoleID:        int64(request.RoleID),
		Status:        "ACTIVE",
		Phone:         request.Phone,
		Email:         request.Email,
		CreatedBy:     userID,
	}

	err = u.Repository.Create(tx, &dataInsert)
	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [USERS CREATE] RESPONSE : %+v", session, err.Error())
		return response
	}

	response = helper.Response("00", "success", nil)
	slog.Infof("[%+v] [USERS CREATE] RESPONSE : %+v", session, response)
	return response
}

func (u *useCase) Update(r *http.Request) *helper.WebResponse[interface{}] {
	response := &helper.WebResponse[interface{}]{}
	session := helper.GenerateRandomString()
	tx := u.DB.WithContext(r.Context())
	request := new(UserUpdateRequest)
	requestID := r.URL.Query().Get("id")

	err := helper.ValidateRequest(r, u.Validate, request)
	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [USERS UPDATE] RESPONSE : %+v", session, err.Error())
		return response
	}

	slog.Infof("[%+v] [USERS UPDATE] REQUEST : %+v", session, request)

	id, err := strconv.ParseInt(requestID, 10, 64)
	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [USERS FIND] RESPONSE : %+v", session, err.Error())
		return response
	}

	users := new(Users)
	err = u.Repository.CheckUsersByID(tx, users, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response = helper.Response("400", "User not found", nil)
		slog.Warnf("[%+v] [USERS Update] RESPONSE : %+v", session, response)
		return response
	}

	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [USERS Update] RESPONSE : %+v", session, err.Error())
		return response
	}

	getUserID := fmt.Sprint(context.Context.Value(r.Context(), helper.GetContextKey()))
	userID, err := strconv.ParseInt(getUserID, 10, 64)
	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [USERS Update] RESPONSE : %+v", session, err.Error())
		return response
	}

	perusahaanID := int64(*request.PerusahaanID)
	dataUpdate := entity.Users{
		Username:      request.Username,
		Nip:           request.NIP,
		CorporationID: &perusahaanID,
		NRK:           request.NRK,
		Instance:      request.Instansi,
		Position:      request.Jabatan,
		FullName:      *request.Fullname,
		RoleID:        int64(request.RoleID),
		Status:        "ACTIVE",
		Phone:         request.Phone,
		Email:         request.Email,
		UpdatedBy:     userID,
		UpdatedOn:     time.Now(),
	}

	err = u.Repository.Update(tx, &dataUpdate, id)
	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [USERS Update] RESPONSE : %+v", session, err.Error())
		return response
	}

	response = helper.Response("00", "success", nil)
	slog.Infof("[%+v] [USERS UPDATE] REQUEST : %+v", session, response)
	return response
}

func (u *useCase) Delete(r *http.Request) *helper.WebResponse[interface{}] {
	response := &helper.WebResponse[interface{}]{}
	session := helper.GenerateRandomString()
	tx := u.DB.WithContext(r.Context())
	request := r.URL.Query().Get("id")

	slog.Infof("[%+v] [USERS DELETE] REQUEST : %+v", session, request)

	id, err := strconv.ParseInt(request, 10, 64)
	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [USERS FIND] RESPONSE : %+v", session, err.Error())
		return response
	}

	users := new(Users)
	err = u.Repository.CheckUsersByID(tx, users, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		users = nil
	}

	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [USERS FIND] RESPONSE : %+v", session, err.Error())
		return response
	}

	getUserID := fmt.Sprint(context.Context.Value(r.Context(), helper.GetContextKey()))
	userID, err := strconv.ParseInt(getUserID, 10, 64)
	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [USERS Update] RESPONSE : %+v", session, err.Error())
		return response
	}

	dataUpdate := entity.Users{
		Status:    "DELETED",
		UpdatedBy: userID,
		UpdatedOn: time.Now(),
	}

	err = u.Repository.Update(tx, &dataUpdate, id)
	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [USERS Update] RESPONSE : %+v", session, err.Error())
		return response
	}

	response = helper.Response("00", "success", nil)
	slog.Infof("[%+v] [USERS DELETE] RESPONSE : %+v", session, request)
	return response
}
