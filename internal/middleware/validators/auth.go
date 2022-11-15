package validators

import (
	"bytes"
	"encoding/json"
	"github.com/Vitaly-Baidin/auth-api/internal/model"
	"io"
	"net/http"
)

func RegisterValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		var reqBody model.RegisterRequest
		res := &model.Response{
			StatusCode: http.StatusBadRequest,
			Success:    false,
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			res.Message = err.Error()
			res.SendResponse(rw)
			return
		}

		err = json.Unmarshal(body, &reqBody)
		if err != nil {
			res.Message = err.Error()
			res.SendResponse(rw)
			return
		}

		r.Body = io.NopCloser(bytes.NewBuffer(body))

		if err := reqBody.Validate(); err != nil {
			model.SendErrorResponse(rw, http.StatusBadRequest, err.Error())
			return
		}

		next.ServeHTTP(rw, r)
	})
}

func LoginValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		var reqBody model.LoginRequest
		res := &model.Response{
			StatusCode: http.StatusBadRequest,
			Success:    false,
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			res.Message = err.Error()
			res.SendResponse(rw)
			return
		}

		err = json.Unmarshal(body, &reqBody)
		if err != nil {
			res.Message = err.Error()
			res.SendResponse(rw)
			return
		}

		r.Body = io.NopCloser(bytes.NewBuffer(body))

		if err := reqBody.Validate(); err != nil {
			model.SendErrorResponse(rw, http.StatusBadRequest, err.Error())
			return
		}

		next.ServeHTTP(rw, r)
	})
}

func RefreshValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		var reqBody model.RefreshRequest
		res := &model.Response{
			StatusCode: http.StatusBadRequest,
			Success:    false,
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			res.Message = err.Error()
			res.SendResponse(rw)
			return
		}

		err = json.Unmarshal(body, &reqBody)
		if err != nil {
			res.Message = err.Error()
			res.SendResponse(rw)
			return
		}

		r.Body = io.NopCloser(bytes.NewBuffer(body))

		if err := reqBody.Validate(); err != nil {
			model.SendErrorResponse(rw, http.StatusBadRequest, err.Error())
			return
		}

		next.ServeHTTP(rw, r)
	})
}
