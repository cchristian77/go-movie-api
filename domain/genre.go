package domain

import (
	"context"
	"github.com/google/uuid"
	"go-movie-api/utils"
	"gorm.io/gorm"
	"time"
)

type Genre struct {
	ID        uint           `gorm:"primarykey" json:"-"`
	Uuid      uuid.UUID      `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Name      string         `json:"name"`
	Movies    []Movie        `json:"movies,omitempty" gorm:"many2many:movie_genres;" `
}

type GenreService interface {
	FetchPagination(ctx context.Context, page int, perPage int) ([]Genre, utils.Pagination, error)
	GetByID(ctx context.Context, uuid uuid.UUID) (Genre, error)
	Store(ctx context.Context, genre *Genre) (Genre, error)
	Update(ctx context.Context, genre *Genre) error
	SoftDelete(ctx context.Context, uuid uuid.UUID) error
	Delete(ctx context.Context, uuid2 uuid.UUID) error
}

type GenreRepository interface {
	FetchPagination(ctx context.Context, pagination *utils.Pagination) ([]Genre, error)
	GetByID(ctx context.Context, uuid uuid.UUID) (Genre, error)
	Store(ctx context.Context, genre *Genre) (Genre, error)
	Update(ctx context.Context, genre *Genre) error
	SoftDelete(ctx context.Context, uuid uuid.UUID) error
	Delete(ctx context.Context, uuid uuid.UUID) error
}
