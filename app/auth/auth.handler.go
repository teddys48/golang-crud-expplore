package auth

import (
	"net/http"

	"github.com/teddys48/kmpro/helper"
)

type AuthHandler interface {
	Login(w http.ResponseWriter, r *http.Request)
	RefreshToken(w http.ResponseWriter, r *http.Request)
}

type authHandler struct {
	UseCase AuthUseCase
}

func NewAuthHandler(useCase AuthUseCase) AuthHandler {
	return &authHandler{
		UseCase: useCase,
	}
}

func (h authHandler) Login(w http.ResponseWriter, r *http.Request) {
	res := h.UseCase.Login(r)
	helper.ReturnResponse(w, res)
}

func (h authHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	res := h.UseCase.RefreshToken(r)
	helper.ReturnResponse(w, res)
}
