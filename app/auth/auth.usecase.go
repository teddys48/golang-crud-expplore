package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gookit/slog"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"github.com/teddys48/kmpro/helper"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthUseCase interface {
	Login(r *http.Request) *helper.WebResponse[interface{}]
	RefreshToken(r *http.Request) *helper.WebResponse[interface{}]
}

type authUseCase struct {
	DB             *gorm.DB
	Validate       *validator.Validate
	AuthRepository AuthRepository
	Config         *viper.Viper
	Redis          *redis.Client
}

func NewAuthUseCase(db *gorm.DB, validate *validator.Validate, AuthRepository AuthRepository, viper *viper.Viper, redis *redis.Client) AuthUseCase {
	return &authUseCase{
		DB:             db,
		Validate:       validate,
		AuthRepository: AuthRepository,
		Config:         viper,
		Redis:          redis,
	}
}

func (u authUseCase) Login(r *http.Request) *helper.WebResponse[interface{}] {
	response := &helper.WebResponse[interface{}]{}
	request := new(LoginRequest)

	session := helper.GenerateRandomString()

	err := helper.ValidateRequest(r, u.Validate, request)
	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [AUTH LOGIN] RESPONSE : %+v", session, err.Error())
		return response
	}

	slog.Infof("[%+v] [AUTH LOGIN] REQUEST : %+v", session, request)

	tx := u.DB.WithContext(r.Context())

	user := new(LoginUsers)
	err = u.AuthRepository.CheckUsersByEmailOrNIP(tx, user, request.EmailOrNIP)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response = helper.Response("400", "User not found", nil)
		slog.Warnf("[%+v] [AUTH LOGIN] RESPONSE : %+v", session, err.Error())
		return response
	} else if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [AUTH LOGIN] RESPONSE : %+v", session, err.Error())
		return response
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		response = helper.Response("400", "Wrong password", nil)
		slog.Warnf("[%+v] [AUTH LOGIN] RESPONSE : %+v", session, response)
		return response
	}

	encryptKey := u.Config.GetString("encrypt.key")
	jwtKey, err := helper.PrivateKey()
	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [AUTH LOGIN] RESPONSE : %+v", session, err.Error())
		return response
	}

	userIDEnc, err := helper.Encrypt([]byte(fmt.Sprint(user.ID)), []byte(encryptKey))
	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [AUTH LOGIN] RESPONSE : %+v", session, err.Error())
		return response
	}

	var wg sync.WaitGroup
	var errChanToken = make(chan error)
	var accessToken, refreshToken string

	accessTokenExp := u.Config.GetInt("jwt.tokenExpiration")
	refreshTokenExp := u.Config.GetInt("jwt.refreshTokenExpiration")

	wg.Add(2)
	go func() {
		defer wg.Done()
		accessTokenClaim := ClaimsToken{
			UserID:           userIDEnc,
			RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(accessTokenExp)))},
		}

		accessToken, err = helper.GenerateToken(jwtKey, accessTokenClaim)
		if err != nil {
			errChanToken <- err
			return
		}
	}()

	go func() {
		defer wg.Done()
		refeshTokenClaim := ClaimsToken{
			UserID:           userIDEnc,
			RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(refreshTokenExp)))},
		}

		refreshToken, err = helper.GenerateToken(jwtKey, refeshTokenClaim)
		if err != nil {
			errChanToken <- err
			return
		}
	}()

	go func() {
		wg.Wait()
		close(errChanToken)
	}()

	for e := range errChanToken {
		if e != nil {
			err = e
			break
		}
	}

	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [AUTH LOGIN] RESPONSE : %+v", session, err.Error())
		return response
	}

	menu := new([]Menu)

	err = u.AuthRepository.GetRoleDetailData(tx, menu, int(user.RoleID))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		slog.Printf("[%+v] [AUTH LOGIN] RESPONSE : %+v", session, err.Error())
	} else if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [AUTH LOGIN] RESPONSE : %+v", session, err.Error())
		return response
	}

	userData := new(UsersData)
	err = u.AuthRepository.CheckUsersByEmailOrNIP2(tx, userData, request.EmailOrNIP)
	// if errors.Is(err, gorm.ErrRecordNotFound) {
	// 	slog.Printf("[%+v] [AUTH LOGIN] RESPONSE : %+v", session, err.Error())
	// } else
	fmt.Println(userData)
	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [AUTH LOGIN] RESPONSE : %+v", session, err.Error())
		return response
	}

	response = helper.Response("00", "success", LoginResponse{AccessToken: accessToken, RefreshToken: refreshToken, Menu: menu, User: userData})
	slog.Infof("[%+v] [AUTH LOGIN] RESPONSE : %+v", session, response)
	return response
}

func (u authUseCase) RefreshToken(r *http.Request) *helper.WebResponse[interface{}] {
	response := &helper.WebResponse[interface{}]{}
	userID := context.Context.Value(r.Context(), "user_id")
	fmt.Println("user_id", userID)
	session := helper.GenerateRandomString()

	encryptKey := u.Config.GetString("encrypt.key")
	jwtKey, err := helper.PrivateKey()
	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [AUTH REFRESH TOKEN] RESPONSE : %+v", session, err.Error())
		return response
	}

	slog.Infof("[%+v] [AUTH REFRESH TOKEN] REQUEST : %+v", session, nil)

	userIDEnc, err := helper.Encrypt([]byte(fmt.Sprint(userID)), []byte(encryptKey))
	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [AUTH REFRESH TOKEN] RESPONSE : %+v", session, err.Error())
		return response
	}

	var wg sync.WaitGroup
	var errChanToken = make(chan error)
	var accessToken, refreshToken string

	accessTokenExp := u.Config.GetInt("jwt.tokenExpiration")
	refreshTokenExp := u.Config.GetInt("jwt.refreshTokenExpiration")

	wg.Add(2)
	go func() {
		defer wg.Done()
		accessTokenClaim := ClaimsToken{
			UserID:           userIDEnc,
			RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(accessTokenExp)))},
		}

		accessToken, err = helper.GenerateToken(jwtKey, accessTokenClaim)
		if err != nil {
			errChanToken <- err
			return
		}
	}()

	go func() {
		defer wg.Done()
		refeshTokenClaim := ClaimsToken{
			UserID:           userIDEnc,
			RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(refreshTokenExp)))},
		}

		refreshToken, err = helper.GenerateToken(jwtKey, refeshTokenClaim)
		if err != nil {
			errChanToken <- err
			return
		}
	}()

	go func() {
		wg.Wait()
		close(errChanToken)
	}()

	for e := range errChanToken {
		if e != nil {
			err = e
			break
		}
	}

	if err != nil {
		response = helper.Response("500", err.Error(), nil)
		slog.Warnf("[%+v] [AUTH REFRESH TOKEN] RESPONSE : %+v", session, err.Error())
		return response
	}

	response = helper.Response("00", "success", LoginResponse{AccessToken: accessToken, RefreshToken: refreshToken})
	slog.Warnf("[%+v] [AUTH REFRESH TOKEN] RESPONSE : %+v", session, response)
	return response
}
