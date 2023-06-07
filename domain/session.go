package domain

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type Session struct {
	ID           uuid.UUID
	UserID       uint
	RefreshToken string
	UserAgent    string
	ClientIp     string
	IsBlocked    bool
	ExpiresAt    time.Time
	CreatedAt    time.Time
}

type SessionRepository interface {
	Store(ctx context.Context, session *Session) (Session, error)
	FindByID(ctx context.Context, id uuid.UUID) (Session, error)
	Delete(ctx context.Context, id uuid.UUID) error
	BlockSession(ctx context.Context, id uuid.UUID) error
}
