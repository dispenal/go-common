package jwt

import (
	"context"
	"net/http"

	common_utils "github.com/dispenal/go-common/utils"
	"github.com/google/uuid"
)

type key string

const JWT_PAYLOAD key = "jwt-payload"

func AppendRequestCtx(r *http.Request, ctxKey key, input interface{}) context.Context {
	return context.WithValue(r.Context(), ctxKey, input)
}

func GetRequestCtx(r *http.Request, ctxKey key) *Payload {
	return r.Context().Value(ctxKey).(*Payload)
}

func CheckIsAuthorize(r *http.Request, accessId uuid.UUID) {
	jwtPayload := GetRequestCtx(r, JWT_PAYLOAD)

	if jwtPayload.UserId != accessId && jwtPayload.Role != "admin" {
		common_utils.PanicIfError(common_utils.CustomError("not authorize to perform this operation", 403))
	}
}
