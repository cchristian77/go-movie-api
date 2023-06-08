package domain

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type Session struct {
	ID                    uuid.UUID
	UserID                uint
	AccessToken           string
	RefreshToken          string
	AccessTokenExpiresAt  time.Time
	AccessTokenCreatedAt  time.Time
	RefreshTokenExpiresAt time.Time
	RefreshTokenCreatedAt time.Time
	UserAgent             string
	ClientIp              string
	IsRevoked             bool
}

type SessionRepository interface {
	Store(ctx context.Context, session *Session) (Session, error)
	FindByID(ctx context.Context, id uuid.UUID) (Session, error)
	Delete(ctx context.Context, session *Session) error
	BlockSession(ctx context.Context, id uuid.UUID) error
}
