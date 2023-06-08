package domain

import (
	"context"
	"github.com/google/uuid"
	"go-movie-api/token"
)

type AuthService interface {
	Authenticate(ctx context.Context, user *User) (User, error)
	CreateSession(ctx context.Context, session *Session) (Session, error)
	VerifySession(ctx context.Context, payload *token.Payload, refreshToken string) error
	RevokeSession(ctx context.Context, sessionID uuid.UUID) error
	DeleteOldSession(ctx context.Context, session *Session) error
}
