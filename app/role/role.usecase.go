package role

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

	slog.Infof("[%+v] [ROLE ALL] REQUEST : %+v", session, nil)

	role := new([]Role)
	err := u.Repository.All(tx, role)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		role = nil
	}

	roleDataMenu := []RoleDataMenu{}

	for _, v := range *role {
		roleDetail := new([]RoleDetailData)
		err = u.Repository.GetRoleDetailData(tx, roleDetail, int(v.ID))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			slog.Infof("[%+v] [ROLE ALL] RESPONSE : %+v", session, err.Error())
		} else if err != nil {
			response = helper.Response("500", err.Error(), nil)
			slog.Warnf("[%+v] [ROLE ALL] RESPONSE : %+v", session, err.Error())
			return response
		}

		data := RoleDataMenu{
			ID:          v.ID,
			Code:        v.Code,
			Name:        v.Name,
			Description: v.Description,
			Status:      v.Status,
			CreatedBy:   v.CreatedBy,
			CreatedOn:   v.CreatedOn,
			UpdatedBy:   v.UpdatedBy,
			UpdatedOn:   v.UpdatedOn,
			MenuList:    *roleDetail,
		}

		roleDataMenu = append(roleDataMenu, data)
	}

	// if err != nil {
	// 	response = helper.Response("500", err.Error(), nil)
	// 	slog.Warnf("[%+v] [ROLE ALL] RESPONSE : %+v", session, err.Error())
	// 	return response
	// }

	response = helper.Response("00", "success", roleDataMenu)
	slog.Infof("[%+v] [ROLE ALL] RESPONSE : %+v", session, response)
	return response
}

func (u *useCase) Find(r *http.Request) *helper.WebResponse[interface{}] {
	response := &helper.WebResponse[interface{}]{}
	session := helper.GenerateRandomString()
	tx := u.DB.WithContext(r.Context())
	request := r.URL.Query().Get("id")

	slog.Infof("[%+v] [ROLE FIND] REQUEST : %+v", session, request)

	id, err := strconv.ParseInt(request, 10, 64)
	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [ROLE FIND] RESPONSE : %+v", session, err.Error())
		return response
	}

	users := new(Role)
	err = u.Repository.CheckByID(tx, users, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response = helper.Response("400", "User not found", nil)
		slog.Warnf("[%+v] [ROLE FIND] RESPONSE : %+v", session, response)
		return response
	}

	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [ROLE FIND] RESPONSE : %+v", session, err.Error())
		return response
	}

	response = helper.Response("00", "success", users)
	slog.Warnf("[%+v] [ROLE FIND] RESPONSE : %+v", session, response)
	return response
}

func (u *useCase) Create(r *http.Request) *helper.WebResponse[interface{}] {
	response := &helper.WebResponse[interface{}]{}
	session := helper.GenerateRandomString()
	tx := u.DB.WithContext(r.Context())
	request := new(RoleCreateRequest)

	err := helper.ValidateRequest(r, u.Validate, request)
	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [ROLE CREATE] RESPONSE : %+v", session, err.Error())
		return response
	}

	slog.Infof("[%+v] [ROLE CREATE] REQUEST : %+v", session, request)

	getUserID := fmt.Sprint(context.Context.Value(r.Context(), helper.GetContextKey()))
	userID, err := strconv.ParseInt(getUserID, 10, 64)
	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [ROLE CREATE] RESPONSE : %+v", session, err.Error())
		return response
	}

	dataRole := entity.Role{
		Code:        "ROLE" + fmt.Sprint(time.Now().Unix()),
		Name:        request.Name,
		Description: request.Description,
		Status:      "ACTIVE",
		CreatedBy:   userID,
		CreatedOn:   time.Now(),
	}

	dataRoleDetail := []entity.RoleDetail{}
	for _, v := range request.Menu {
		dataRoleDetail = append(dataRoleDetail, entity.RoleDetail{MenuID: v.MenuID, Action: v.Action})
	}

	fmt.Println("cek", dataRoleDetail)

	err = u.Repository.InsertUpdateRole(tx, &dataRole, &dataRoleDetail)
	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [ROLE CREATE] RESPONSE : %+v", session, err.Error())
		return response
	}

	response = helper.Response("00", "success", nil)
	slog.Infof("[%+v] [ROLE CREATE] RESPONSE : %+v", session, response)
	return response
}

func (u *useCase) Update(r *http.Request) *helper.WebResponse[interface{}] {
	response := &helper.WebResponse[interface{}]{}
	session := helper.GenerateRandomString()
	tx := u.DB.WithContext(r.Context())
	request := new(RoleCreateRequest)
	requestID := r.URL.Query().Get("id")

	err := helper.ValidateRequest(r, u.Validate, request)
	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [ROLE UPDATE] RESPONSE : %+v", session, err.Error())
		return response
	}

	slog.Infof("[%+v] [ROLE UPDATE] REQUEST : %+v", session, request)

	id, err := strconv.ParseInt(requestID, 10, 64)
	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [USERS FIND] RESPONSE : %+v", session, err.Error())
		return response
	}

	getUserID := fmt.Sprint(context.Context.Value(r.Context(), helper.GetContextKey()))
	userID, err := strconv.ParseInt(getUserID, 10, 64)
	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [ROLE UPDATE] RESPONSE : %+v", session, err.Error())
		return response
	}

	dataRole := entity.Role{
		// Code:        "ROLE" + fmt.Sprint(time.Now().Unix()),
		Name:        request.Name,
		Description: request.Description,
		Status:      "ACTIVE",
		UpdatedBy:   userID,
		UpdatedOn:   time.Now(),
	}

	dataRoleDetail := []entity.RoleDetail{}
	for _, v := range request.Menu {
		dataRoleDetail = append(dataRoleDetail, entity.RoleDetail{MenuID: v.MenuID, Action: v.Action, RoleID: id})
	}

	err = u.Repository.UpdateUpdateRole(tx, &dataRole, &dataRoleDetail, id)
	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [ROLE UPDATE] RESPONSE : %+v", session, err.Error())
		return response
	}

	response = helper.Response("00", "success", nil)
	slog.Infof("[%+v] [ROLE UPDATE] RESPONSE : %+v", session, response)
	return response
}

func (u *useCase) Delete(r *http.Request) *helper.WebResponse[interface{}] {
	response := &helper.WebResponse[interface{}]{}
	session := helper.GenerateRandomString()
	tx := u.DB.WithContext(r.Context())
	request := r.URL.Query().Get("id")

	slog.Infof("[%+v] [ROLE FIND] REQUEST : %+v", session, request)

	getUserID := fmt.Sprint(context.Context.Value(r.Context(), helper.GetContextKey()))
	userID, err := strconv.ParseInt(getUserID, 10, 64)
	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [ROLE CREATE] RESPONSE : %+v", session, err.Error())
		return response
	}

	id, err := strconv.ParseInt(request, 10, 64)
	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [ROLE FIND] RESPONSE : %+v", session, err.Error())
		return response
	}

	role := entity.Role{
		Status:    "DELETED",
		UpdatedBy: userID,
		UpdatedOn: time.Now(),
	}

	err = u.Repository.Update(tx, &role, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response = helper.Response("400", "User not found", nil)
		slog.Warnf("[%+v] [ROLE FIND] RESPONSE : %+v", session, response)
		return response
	}

	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [ROLE FIND] RESPONSE : %+v", session, err.Error())
		return response
	}

	response = helper.Response("00", "success", nil)
	slog.Warnf("[%+v] [ROLE FIND] RESPONSE : %+v", session, response)
	return response
}
