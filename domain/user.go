package domain

import (
	"context"
	"github.com/google/uuid"
	"go-movie-api/utils"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID                uint           `gorm:"primarykey" json:"-"`
	Uuid              uuid.UUID      `json:"id"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `json:"-" gorm:"index"`
	Username          string         `json:"username"`
	Email             string         `json:"email"`
	FullName          string         `json:"full_name"`
	Password          string         `json:"-"`
	IsAdmin           bool           `json:"is_admin"`
	IsEmailVerified   bool           `json:"is_email_verified"`
	PasswordChangedAt time.Time      `json:"password_changed_at"`
}

type UserService interface {
	FetchPagination(ctx context.Context, page int, perPage int) ([]User, utils.Pagination, error)
	FindByID(ctx context.Context, uuid uuid.UUID) (User, error)
	Store(ctx context.Context, user *User) (User, error)
	Update(ctx context.Context, user *User) error
	SoftDelete(ctx context.Context, uuid uuid.UUID) error
	Delete(ctx context.Context, uuid2 uuid.UUID) error
	CreateSession(ctx context.Context, session *Session) (Session, error)
	Login(ctx context.Context, user *User) (User, error)
}

type UserRepository interface {
	FindByID(ctx context.Context, uuid uuid.UUID) (User, error)
	FindByUsernameOrEmail(ctx context.Context, username string, email string) (User, error)
	FindByUsername(ctx context.Context, username string) (User, error)
	FindByEmail(ctx context.Context, email string) (User, error)
	FetchPagination(ctx context.Context, pagination *utils.Pagination) ([]User, error)
	FindByIDForUpdate(ctx context.Context, uuid uuid.UUID) (User, error)
	Store(ctx context.Context, user *User) (User, error)
	Update(ctx context.Context, user *User) error
	SoftDelete(ctx context.Context, uuid uuid.UUID) error
	Delete(ctx context.Context, uuid uuid.UUID) error
}
