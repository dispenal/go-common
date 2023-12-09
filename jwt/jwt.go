package jwt

import (
	"time"
)

type JWT interface {
	CreateToken(params PayloadParams, duration time.Duration) (string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
}
