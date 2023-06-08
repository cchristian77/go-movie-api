package token

import (
	"github.com/google/uuid"
	"time"
)

var TokenMaker Maker

// Maker is an interface for managing tokens
type Maker interface {
	// GenerateToken creates a new token for a specific username and duration
	GenerateToken(id uuid.UUID, userUuid uuid.UUID, duration time.Duration) (string, *Payload, error)

	// VerifyToken checks if the token is valid or not
	VerifyToken(token string) (*Payload, error)
}
