package middleware

import (
	"net/http"
	"strings"

	jwtMaker "github.com/dispenal/go-common/jwt"
	common_utils "github.com/dispenal/go-common/utils"
)

const (
	invalid_token           = "invalid token"
	invalid_token_type      = "invalid token type"
	forbidden_inactive_user = "inactive user can't access this route"
)

type AuthMiddleware interface {
	CheckIsAuthenticated(handler http.Handler) http.Handler
	CheckIsRefresh(handler http.Handler) http.Handler
}

type AuthMiddlewareImpl struct {
	jwt jwtMaker.JWT
}

func NewAuthMiddleware(jwt jwtMaker.JWT) AuthMiddleware {
	return &AuthMiddlewareImpl{
		jwt: jwt,
	}
}

func (m *AuthMiddlewareImpl) CheckIsAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")

		if header == "" || !strings.Contains(header, "Bearer ") {
			common_utils.PanicIfError(common_utils.CustomError("unauthorized", 401))
		}

		token := strings.Split(header, " ")[1]
		payload, err := m.jwt.VerifyToken(token)
		if err != nil {
			common_utils.PanicIfError(common_utils.CustomErrorWithTrace(err, invalid_token, 401))
		}

		if payload.TokenType != "access" {
			common_utils.PanicIfError(common_utils.CustomError(invalid_token_type, 401))
		}

		if payload.Status != "active" {
			common_utils.PanicIfError(common_utils.CustomError(forbidden_inactive_user, 403))
		}

		ctx := r.WithContext(jwtMaker.AppendRequestCtx(r, jwtMaker.JWT_PAYLOAD, payload))

		next.ServeHTTP(w, ctx)
	})
}

func (m *AuthMiddlewareImpl) CheckIsRefresh(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")

		if header == "" || !strings.Contains(header, "Bearer ") {
			common_utils.PanicIfError(common_utils.CustomError("unauthorized", 401))
		}

		token := strings.Split(header, " ")[1]
		payload, err := m.jwt.VerifyToken(token)
		if err != nil {
			common_utils.PanicIfError(common_utils.CustomErrorWithTrace(err, invalid_token, 401))
		}

		if payload.TokenType != "refresh" {
			common_utils.PanicIfError(common_utils.CustomError(invalid_token_type, 401))
		}

		if payload.Status != "active" {
			common_utils.PanicIfError(common_utils.CustomError(forbidden_inactive_user, 403))
		}

		ctx := r.WithContext(jwtMaker.AppendRequestCtx(r, jwtMaker.JWT_PAYLOAD, payload))

		next.ServeHTTP(w, ctx)
	})
}

func (m *AuthMiddlewareImpl) CheckIsAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")

		if header == "" || !strings.Contains(header, "Bearer ") {
			common_utils.PanicIfError(common_utils.CustomError("unauthorized", 401))
		}

		token := strings.Split(header, " ")[1]
		payload, err := m.jwt.VerifyToken(token)
		if err != nil {
			common_utils.PanicIfError(common_utils.CustomErrorWithTrace(err, invalid_token, 401))
		}

		if payload.TokenType != "access" {
			common_utils.PanicIfError(common_utils.CustomError(invalid_token_type, 401))
		}

		if payload.Role != "admin" {
			common_utils.PanicIfError(common_utils.CustomError("invalid role", 403))
		}

		if payload.Status != "active" {
			common_utils.PanicIfError(common_utils.CustomError(forbidden_inactive_user, 403))
		}

		ctx := r.WithContext(jwtMaker.AppendRequestCtx(r, jwtMaker.JWT_PAYLOAD, payload))

		next.ServeHTTP(w, ctx)
	})
}
