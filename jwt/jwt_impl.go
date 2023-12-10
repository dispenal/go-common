package jwt

import (
	"fmt"
	"time"

	common_utils "github.com/dispenal/go-common/utils"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/chacha20poly1305"
)

type JWTMaker struct {
	symmetricKey []byte
	jwt          *jwt.Token
}

func NewJWTMaker(config *common_utils.BaseConfig) (JWT, error) {
	symmetricKey := config.JwtSecretKey

	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
	}

	symmetricKeyByte := []byte(symmetricKey)

	jwt := jwt.New(jwt.SigningMethodHS256)

	return &JWTMaker{
		symmetricKey: symmetricKeyByte,
		jwt:          jwt,
	}, nil
}

func (m *JWTMaker) CreateToken(params PayloadParams, duration time.Duration) (string, *Payload, error) {
	payload := NewPayload(params, duration)

	m.jwt.Claims = payload

	token, err := m.jwt.SignedString(m.symmetricKey)

	return token, payload, err

}

func (m *JWTMaker) VerifyToken(token string) (*Payload, error) {
	var payload = &Payload{}

	jwtAuth, err := jwt.ParseWithClaims(token,
		payload, func(token *jwt.Token) (interface{}, error) {
			return m.symmetricKey, nil
		})

	if err != nil {
		return nil, ErrInvalidToken
	}

	claims, ok := jwtAuth.Claims.(*Payload)

	if !ok && !jwtAuth.Valid {
		return nil, ErrInvalidToken
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return claims, nil

}
