package jobs

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
	Approve(r *http.Request) *helper.WebResponse[interface{}]
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

	slog.Infof("[%+v] [JOBS ALL] REQUEST : %+v", session, nil)

	data := new([]Jobs)
	err := u.Repository.All(tx, data)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		data = nil
	}

	responseData := []JobsAllResponse{}
	project := new([]Project)
	projectID := []int64{}
	for _, v := range *data {

		projectID = append(projectID, v.ProjectID)
		// err = u.Repository.GetProjectByID(tx, project, v.ProjectID)
		// if err != nil {
		// 	response = helper.Response("500", err.Error(), nil)
		// 	slog.Warnf("[%+v] [JOBS ALL] RESPONSE : %+v", session, response)
		// 	return response
		// }
		// responseData = append(responseData, JobsAllResponse{
		// 	ProjectName:                   project.ActivityName,
		// 	TargetPercentage:              project.TargetPercentage,
		// 	ProgressPlan:                  project.ProgressPlan,
		// 	TotalCumulativeProgressPlan:   project.TotalCumulativeProgressPlan,
		// 	ActualProgress:                project.ActualProgress,
		// 	TotalCumulativeActualProgress: project.TotalCumulativeActualProgress,
		// 	Deviation:                     project.Deviation,
		// 	// JobList:                       data,
		// })
	}

	err = u.Repository.GetProjectByID(tx, project, projectID)
	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [JOBS ALL] RESPONSE : %+v", session, err.Error())
		return response
	}

	for _, v := range *project {
		err = u.Repository.CheckByProjectID(tx, data, v.ID)
		if err != nil {
			response = helper.Response("500", err.Error(), nil)
			slog.Warnf("[%+v] [JOBS ALL] RESPONSE : %+v", session, err.Error())
			return response
		}

		responseData = append(responseData, JobsAllResponse{
			ProjectName:                   v.ActivityName,
			ProgressPlan:                  v.ProgressPlan,
			TargetPercentage:              v.TargetPercentage,
			TotalCumulativeProgressPlan:   v.TotalCumulativeActualProgress,
			ActualProgress:                v.ActualProgress,
			TotalCumulativeActualProgress: v.TotalCumulativeProgressPlan,
			Deviation:                     v.Deviation,
			JobList:                       data,
		})
	}

	// for _, v := range responseData {
	// 	err = u.Repository.CheckByID(tx, data, v.)
	// }

	// responseData.JobList = data

	response = helper.Response("00", "success", responseData)
	slog.Infof("[%+v] [JOBS ALL] RESPONSE : %+v", session, response)
	return response
}

func (u *useCase) Find(r *http.Request) *helper.WebResponse[interface{}] {
	response := &helper.WebResponse[interface{}]{}
	session := helper.GenerateRandomString()
	tx := u.DB.WithContext(r.Context())
	request := r.URL.Query().Get("id")

	slog.Infof("[%+v] [JOBS FIND] REQUEST : %+v", session, request)

	id, err := strconv.ParseInt(request, 10, 64)
	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [JOBS FIND] RESPONSE : %+v", session, err.Error())
		return response
	}

	data := new(Jobs)
	err = u.Repository.CheckByID(tx, data, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response = helper.Response("400", "Data not found", nil)
		slog.Warnf("[%+v] [JOBS FIND] RESPONSE : %+v", session, response)
		return response
	}

	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [JOBS FIND] RESPONSE : %+v", session, err.Error())
		return response
	}

	response = helper.Response("00", "success", data)
	slog.Warnf("[%+v] [JOBS FIND] RESPONSE : %+v", session, response)
	return response
}

func (u *useCase) Create(r *http.Request) *helper.WebResponse[interface{}] {
	response := &helper.WebResponse[interface{}]{}
	session := helper.GenerateRandomString()
	tx := u.DB.WithContext(r.Context())
	request := new(Jobs)

	fmt.Println("print request", request)

	err := helper.ValidateRequest(r, u.Validate, request)
	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [JOBS CREATE] RESPONSE : %+v", session, err.Error())
		return response
	}

	slog.Infof("[%+v] [JOBS CREATE] REQUEST : %+v", session, request)

	getUserID := fmt.Sprint(context.Context.Value(r.Context(), helper.GetContextKey()))
	userID, err := strconv.ParseInt(getUserID, 10, 64)
	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [JOBS CREATE] RESPONSE : %+v", session, err.Error())
		return response
	}

	request.CreatedBy = userID
	request.CreatedOn = time.Now()
	request.Status = "ACTIVE"
	request.Code = "JOBS" + fmt.Sprint(time.Now().Unix())

	err = u.Repository.Create(tx, request)
	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [JOBS CREATE] RESPONSE : %+v", session, err.Error())
		return response
	}

	response = helper.Response("00", "success", nil)
	slog.Infof("[%+v] [JOBS CREATE] RESPONSE : %+v", session, response)
	return response
}

func (u *useCase) Update(r *http.Request) *helper.WebResponse[interface{}] {
	response := &helper.WebResponse[interface{}]{}
	session := helper.GenerateRandomString()
	tx := u.DB.WithContext(r.Context())
	request := new(Jobs)
	requestID := r.URL.Query().Get("id")

	err := helper.ValidateRequest(r, u.Validate, request)
	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [JOBS UPDATE] RESPONSE : %+v", session, err.Error())
		return response
	}

	slog.Infof("[%+v] [JOBS UPDATE] REQUEST : %+v", session, request)

	id, err := strconv.ParseInt(requestID, 10, 64)
	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [JOBS UPDATE] RESPONSE : %+v", session, err.Error())
		return response
	}

	data := new(Jobs)
	err = u.Repository.CheckByID(tx, data, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response = helper.Response("400", "Data not found", nil)
		slog.Warnf("[%+v] [JOBS UPDATE] RESPONSE : %+v", session, response)
		return response
	}

	getUserID := fmt.Sprint(context.Context.Value(r.Context(), helper.GetContextKey()))
	userID, err := strconv.ParseInt(getUserID, 10, 64)
	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [JOBS UPDATE] RESPONSE : %+v", session, err.Error())
		return response
	}

	now := time.Now()

	request.UpdatedBy = &userID
	request.UpdatedOn = &now

	err = u.Repository.Update(tx, request, id)
	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [JOBS UPDATE] RESPONSE : %+v", session, err.Error())
		return response
	}

	response = helper.Response("00", "success", nil)
	slog.Infof("[%+v] [JOBS UPDATE] RESPONSE : %+v", session, response)
	return response
}

func (u *useCase) Delete(r *http.Request) *helper.WebResponse[interface{}] {
	response := &helper.WebResponse[interface{}]{}
	session := helper.GenerateRandomString()
	tx := u.DB.WithContext(r.Context())
	request := r.URL.Query().Get("id")

	slog.Infof("[%+v] [JOBS DELETE] REQUEST : %+v", session, request)

	getUserID := fmt.Sprint(context.Context.Value(r.Context(), helper.GetContextKey()))
	userID, err := strconv.ParseInt(getUserID, 10, 64)
	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [JOBS DELETE] RESPONSE : %+v", session, err.Error())
		return response
	}

	id, err := strconv.ParseInt(request, 10, 64)
	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [JOBS DELETE] RESPONSE : %+v", session, err.Error())
		return response
	}

	data := new(Jobs)
	err = u.Repository.CheckByID(tx, data, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response = helper.Response("400", "Data not found", nil)
		slog.Warnf("[%+v] [JOBS DELETE] RESPONSE : %+v", session, response)
		return response
	}

	now := time.Now()

	dataUpdate := Jobs{
		UpdatedBy: &userID,
		UpdatedOn: &now,
		Status:    "DELETED",
	}

	err = u.Repository.Update(tx, &dataUpdate, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response = helper.Response("400", "User not found", nil)
		slog.Warnf("[%+v] [JOBS DELETE] RESPONSE : %+v", session, response)
		return response
	}

	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [JOBS DELETE] RESPONSE : %+v", session, err.Error())
		return response
	}

	response = helper.Response("00", "success", nil)
	slog.Warnf("[%+v] [JOBS DELETE] RESPONSE : %+v", session, response)
	return response
}

func (u *useCase) Approve(r *http.Request) *helper.WebResponse[interface{}] {
	response := &helper.WebResponse[interface{}]{}
	session := helper.GenerateRandomString()
	tx := u.DB.WithContext(r.Context())
	request := r.URL.Query().Get("id")

	slog.Infof("[%+v] [JOBS APPROVE] REQUEST : %+v", session, request)

	getUserID := fmt.Sprint(context.Context.Value(r.Context(), helper.GetContextKey()))
	userID, err := strconv.ParseInt(getUserID, 10, 64)
	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [JOBS APPROVE] RESPONSE : %+v", session, err.Error())
		return response
	}

	id, err := strconv.ParseInt(request, 10, 64)
	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [JOBS APPROVE] RESPONSE : %+v", session, err.Error())
		return response
	}

	data := new(Jobs)
	err = u.Repository.CheckByID(tx, data, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response = helper.Response("400", "Data not found", nil)
		slog.Warnf("[%+v] [JOBS APPROVE] RESPONSE : %+v", session, response)
		return response
	}

	now := time.Now()

	dataUpdate := new(Jobs)
	if data.ApprovedBy != nil {
		fmt.Println("cek1")
		err = u.Repository.Disapprove(tx, dataUpdate, id)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response = helper.Response("400", "User not found", nil)
			slog.Warnf("[%+v] [JOBS APPROVE] RESPONSE : %+v", session, response)
			return response
		}
	} else {
		fmt.Println("cek2")
		dataUpdate = &Jobs{
			ApprovedBy: &userID,
			ApprovedOn: &now,
		}
		err = u.Repository.Update(tx, dataUpdate, id)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response = helper.Response("400", "User not found", nil)
			slog.Warnf("[%+v] [JOBS APPROVE] RESPONSE : %+v", session, response)
			return response
		}
	}

	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [JOBS APPROVE] RESPONSE : %+v", session, err.Error())
		return response
	}

	response = helper.Response("00", "success", nil)
	slog.Warnf("[%+v] [JOBS APPROVE] RESPONSE : %+v", session, response)
	return response
}
