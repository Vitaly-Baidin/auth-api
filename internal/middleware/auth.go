package middleware

import (
	"github.com/Vitaly-Baidin/auth-api/internal/model"
	db "github.com/Vitaly-Baidin/auth-api/internal/model/db"
	"github.com/Vitaly-Baidin/auth-api/internal/service"
	"net/http"
	"strconv"
)

func JWTMiddleware(next http.Handler, ts service.TokenService) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Bearer-Token")
		tokenModel, err := ts.VerifyToken(r.Context(), token, db.TokenTypeAccess)
		if err != nil {
			model.SendErrorResponse(rw, http.StatusUnauthorized, err.Error())
			return
		}

		rw.Header().Add("userId", strconv.Itoa(int(tokenModel.UserID)))

		next.ServeHTTP(rw, r)
	})
}
