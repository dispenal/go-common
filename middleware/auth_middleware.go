package middleware

import (
	"net/http"
	"strings"

	jwtMaker "github.com/dispenal/go-common/jwt"
	common_utils "github.com/dispenal/go-common/utils"
)

type AuthMiddleware interface {
	CheckIsAuthenticated(handler func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc
}

type AuthMiddlewareImpl struct {
	jwt jwtMaker.JWT
}

func NewAuthMiddleware(jwt jwtMaker.JWT) AuthMiddleware {
	return &AuthMiddlewareImpl{
		jwt: jwt,
	}
}

func (m *AuthMiddlewareImpl) CheckIsAuthenticated(handler func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")

		if header == "" || !strings.Contains(header, "Bearer ") {
			common_utils.PanicIfError(common_utils.CustomError("unauthorized", 401))
		}

		token := strings.Split(header, " ")[1]
		payload, err := m.jwt.VerifyToken(token)
		if err != nil {
			common_utils.PanicIfError(common_utils.CustomErrorWithTrace(err, "invalid token", 401))
		}

		if payload.Status != "active" {
			common_utils.PanicIfError(common_utils.CustomError("inactive user can't access this route", 403))
		}

		ctx := jwtMaker.AppendRequestCtx(r, jwtMaker.JWT_PAYLOAD, payload)
		handler(w, r.WithContext(ctx))
	}
}
