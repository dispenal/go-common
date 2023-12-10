package jwt

import (
	"time"
)

type JWT interface {
	CreateToken(tokenType string, params PayloadParams, duration time.Duration) (string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
}
