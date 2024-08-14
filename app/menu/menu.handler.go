package menu

import (
	"net/http"

	"github.com/teddys48/kmpro/helper"
)

type Handler interface {
	All(w http.ResponseWriter, r *http.Request)
	Find(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	UseCase UseCase
}

func NewHandler(useCase UseCase) Handler {
	return &handler{
		UseCase: useCase,
	}
}

func (h *handler) All(w http.ResponseWriter, r *http.Request) {
	res := h.UseCase.All(r)
	helper.ReturnResponse(w, res)
}

func (h *handler) Find(w http.ResponseWriter, r *http.Request) {
	res := h.UseCase.Find(r)
	helper.ReturnResponse(w, res)
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	res := h.UseCase.Create(r)
	helper.ReturnResponse(w, res)
}

func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	res := h.UseCase.Update(r)
	helper.ReturnResponse(w, res)
}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	res := h.UseCase.Delete(r)
	helper.ReturnResponse(w, res)
}
