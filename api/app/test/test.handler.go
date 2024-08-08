package test

import (
	"net/http"

	"github.com/teddys48/kmpro/helper"
)

type TestHandler interface {
	TestHandler(w http.ResponseWriter, r *http.Request)
}

type testHandler struct {
	TestUsecase TestUsecase
}

func NewTestHandler(TestUsecase TestUsecase) TestHandler {
	return &testHandler{
		TestUsecase: TestUsecase,
	}
}

func (u testHandler) TestHandler(w http.ResponseWriter, r *http.Request) {
	res := u.TestUsecase.TestUsecase(r)
	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(res)
	helper.ReturnResponse(w, res)
}
