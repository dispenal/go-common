package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

type Payload struct {
	jwt.RegisteredClaims
	ID        uuid.UUID `json:"id"`
	UserId    uuid.UUID `json:"userId"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Status    string    `json:"status"`
	TokenType string    `json:"tokenType"`
	IssuedAt  time.Time `json:"issuedAt"`
	ExpiredAt time.Time `json:"expiredAt"`
}

type PayloadParams struct {
	UserId    uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Status    string    `json:"status"`
	TokenType string    `json:"tokenType"`
}

func NewPayload(params PayloadParams, duration time.Duration) *Payload {
	tokenId := uuid.New()

	return &Payload{
		ID:        tokenId,
		UserId:    params.UserId,
		Email:     params.Email,
		Role:      params.Role,
		Status:    params.Status,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
