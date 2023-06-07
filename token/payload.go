package token

import (
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
	ID        uuid.UUID `json:"id"`
	UserUuid  uuid.UUID `json:"user_id"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

// NewPayload creates a new token payload with a specific username and duration
func NewPayload(tokenID uuid.UUID, userUuid uuid.UUID, duration time.Duration) (*Payload, error) {
	return &Payload{
		ID:        tokenID,
		UserUuid:  userUuid,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}, nil
}

// Valid checks if the token payload is valid or not
func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return InvalidTokenErr
	}

	return nil
}
