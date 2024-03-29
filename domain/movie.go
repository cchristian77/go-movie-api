package domain

import (
	"context"
	"github.com/google/uuid"
	"go-movie-api/utils"
	"gorm.io/gorm"
	"time"
)

type Movie struct {
	ID        uint           `gorm:"primarykey" json:"-"`
	Uuid      uuid.UUID      `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	Title     string         `json:"title"`
	Duration  int32          `json:"duration"`
	Year      int32          `json:"year"`
	Synopsis  string         `json:"synopsis"`
	Genres    []Genre        `json:"genres,omitempty" gorm:"many2many:movie_genres;"`
	Ratings   []Rating       `json:"ratings,omitempty"`
}

type MovieService interface {
	FetchPagination(ctx context.Context, page int, perPage int) ([]Movie, utils.Pagination, error)
	FindByID(ctx context.Context, uuid uuid.UUID) (Movie, error)
	Store(ctx context.Context, movie *Movie) (Movie, error)
	Update(ctx context.Context, movie *Movie) error
	SoftDelete(ctx context.Context, uuid uuid.UUID) error
	Delete(ctx context.Context, uuid uuid.UUID) error
}

type MovieRepository interface {
	FetchPagination(ctx context.Context, pagination *utils.Pagination) ([]Movie, error)
	FindByID(ctx context.Context, uuid uuid.UUID) (Movie, error)
	FindByIDForUpdate(ctx context.Context, uuid uuid.UUID) (Movie, error)
	Store(ctx context.Context, movie *Movie) (Movie, error)
	Update(ctx context.Context, movie *Movie) error
	SoftDelete(ctx context.Context, uuid uuid.UUID) error
	Delete(ctx context.Context, uuid uuid.UUID) error
}
