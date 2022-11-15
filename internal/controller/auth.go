package controller

import (
	"encoding/json"
	"github.com/Vitaly-Baidin/auth-api/internal/model"
	db "github.com/Vitaly-Baidin/auth-api/internal/model/db"
	"github.com/Vitaly-Baidin/auth-api/internal/service"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
)

type AuthController interface {
	Register(rw http.ResponseWriter, r *http.Request)
	Login(rw http.ResponseWriter, r *http.Request)
	Refresh(rw http.ResponseWriter, r *http.Request)
}

type AuthContr struct {
	userService  service.UserService
	tokenService service.TokenService
}

func NewAuthController(us service.UserService, ts service.TokenService) *AuthContr {
	return &AuthContr{userService: us, tokenService: ts}
}

func (c *AuthContr) Register(rw http.ResponseWriter, r *http.Request) {
	var reqBody model.RegisterRequest
	ctx := r.Context()
	res := &model.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		res.Message = err.Error()
		res.SendResponse(rw)
		return
	}

	err = c.userService.CheckUserMail(ctx, reqBody.Email)
	if err != nil {
		res.Message = err.Error()
		res.SendResponse(rw)
		return
	}

	err = c.userService.CheckUserLogin(ctx, reqBody.Login)
	if err != nil {
		res.Message = err.Error()
		res.SendResponse(rw)
		return
	}

	reqBody.Login = strings.TrimSpace(reqBody.Login)
	user := db.NewUser(reqBody.Login, reqBody.Email, reqBody.Phone, reqBody.Password)

	user, err = c.userService.CreateUser(ctx, user)
	if err != nil {
		res.Message = err.Error()
		res.SendResponse(rw)
		return
	}

	// generate access tokens
	accessToken, refreshToken, err := c.tokenService.GenerateAccessTokens(ctx, user)
	if err != nil {
		res.Message = err.Error()
		res.SendResponse(rw)
		return
	}

	res.StatusCode = http.StatusCreated
	res.Success = true
	res.Data = map[string]any{
		"user": user,
		"token": map[string]any{
			"access":  accessToken.GetResponseJson(),
			"refresh": refreshToken.GetResponseJson(),
		},
	}

	res.SendResponse(rw)
}

func (c *AuthContr) Login(rw http.ResponseWriter, r *http.Request) {
	var reqBody model.LoginRequest
	ctx := r.Context()

	res := &model.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		res.Message = err.Error()
		res.SendResponse(rw)
		return
	}

	user, err := c.userService.FindUserByEmail(ctx, reqBody.Email)
	if err != nil {
		res.Message = err.Error()
		res.SendResponse(rw)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqBody.Password))
	if err != nil {
		res.Message = "email and password don't match"
		res.SendResponse(rw)
		return
	}

	// generate new access tokens
	accessToken, refreshToken, err := c.tokenService.GenerateAccessTokens(ctx, user)
	if err != nil {
		res.Message = err.Error()
		res.SendResponse(rw)
		return
	}

	res.StatusCode = http.StatusOK
	res.Success = true
	res.Data = map[string]any{
		"user": user,
		"token": map[string]any{
			"access":  accessToken.GetResponseJson(),
			"refresh": refreshToken.GetResponseJson()},
	}
	res.SendResponse(rw)
}

func (c *AuthContr) Refresh(rw http.ResponseWriter, r *http.Request) {
	var reqBody model.RefreshRequest
	ctx := r.Context()

	res := &model.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		res.Message = err.Error()
		res.SendResponse(rw)
		return
	}

	token, err := c.tokenService.VerifyToken(ctx, reqBody.Token, db.TokenTypeRefresh)
	if err != nil {
		res.Message = err.Error()
		res.SendResponse(rw)
		return
	}

	user, err := c.userService.FindUserById(ctx, token.UserID)
	if err != nil {
		res.Message = err.Error()
		res.SendResponse(rw)
		return
	}

	err = c.tokenService.DeleteToken(ctx, token.UserID, token.Type)
	if err != nil {
		res.Message = err.Error()
		res.SendResponse(rw)
		return
	}

	accessToken, refreshToken, err := c.tokenService.GenerateAccessTokens(ctx, user)
	res.StatusCode = http.StatusOK
	res.Success = true
	res.Data = map[string]any{
		"user": user,
		"token": map[string]any{
			"access":  accessToken.GetResponseJson(),
			"refresh": refreshToken.GetResponseJson()},
	}
	res.SendResponse(rw)
}
