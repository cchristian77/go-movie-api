package domain

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Genre struct {
	ID        uint           `gorm:"primarykey" json:"-"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Uuid      uuid.UUID      `json:"id"`
	Name      string         `json:"name"`
	Movies    []Movie        `gorm:"many2many:movie_genres;"`
}

type GenreService interface {
	FetchPagination(ctx context.Context, page int, perPage int) ([]Genre, error)
	GetByID(ctx context.Context, id int32) (Genre, error)
	Store(ctx context.Context, genre *Genre) (Genre, error)
	Update(ctx context.Context, genre *Genre) error
	SoftDelete(ctx context.Context, id int32) error
	Delete(ctx context.Context, id int32) error
}

type GenreRepository interface {
	FetchPagination(ctx context.Context, limit int, offset int) ([]Genre, error)
	GetByID(ctx context.Context, id int32) (Genre, error)
	Store(ctx context.Context, genre *Genre) (Genre, error)
	Update(ctx context.Context, genre *Genre) error
	SoftDelete(ctx context.Context, id int32) error
	Delete(ctx context.Context, id int32) error
}
