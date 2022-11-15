package controller

import (
	"github.com/Vitaly-Baidin/auth-api/internal/model"
	"net/http"
)

func Ping(rw http.ResponseWriter, r *http.Request) {
	response := &model.Response{
		StatusCode: http.StatusOK,
		Success:    true,
		Message:    "pong",
	}

	response.SendResponse(rw)
}
