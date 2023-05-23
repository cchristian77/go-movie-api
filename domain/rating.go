package domain

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Rating struct {
	ID        uint           `gorm:"primarykey" json:"-"`
	Uuid      uuid.UUID      `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	MovieID   uint           `json:"-"`
	Rating    float32        `json:"rating"`
	Comment   string         `json:"comment"`
	Movie     *Movie         `json:"movie,omitempty"`
}

type RatingService interface {
	FindByID(ctx context.Context, uuid uuid.UUID) (Rating, error)
	Store(ctx context.Context, rating *Rating) (Rating, error)
	Update(ctx context.Context, rating *Rating) error
	SoftDelete(ctx context.Context, uuid uuid.UUID) error
	Delete(ctx context.Context, uuid2 uuid.UUID) error
}

type RatingRepository interface {
	FindByID(ctx context.Context, uuid uuid.UUID) (Rating, error)
	FindByIDForUpdate(ctx context.Context, uuid uuid.UUID) (Rating, error)
	Store(ctx context.Context, rating *Rating) (Rating, error)
	Update(ctx context.Context, rating *Rating) error
	SoftDelete(ctx context.Context, uuid uuid.UUID) error
	Delete(ctx context.Context, uuid uuid.UUID) error
}
