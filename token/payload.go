package token

import (
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"time"
)

// Different types of error returned by the VerifyToken function
var (
	InvalidTokenErr = errors.New("token is invalid")
	ExpiredTokenErr = errors.New("token has expired")
)

// Payload contains the payload data of the token
type Payload struct {
	ID       uuid.UUID `json:"id"`
	UserUuid uuid.UUID `json:"user_id"`
	jwt.StandardClaims
}

// NewPayload creates a new token payload with a specific username and duration
func NewPayload(tokenID uuid.UUID, userUuid uuid.UUID, duration time.Duration) (*Payload, error) {
	return &Payload{
		ID:       tokenID,
		UserUuid: userUuid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "movie.api.auth",
		},
	}, nil
}

// Valid checks if the token payload is valid or not
func (payload *Payload) Valid() error {
	expiredAt := time.Unix(payload.StandardClaims.ExpiresAt, 0)

	if time.Now().After(expiredAt) {
		return InvalidTokenErr
	}

	return nil
}
