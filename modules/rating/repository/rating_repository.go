package repository

import (
	"context"
	"github.com/google/uuid"
	"go-movie-api/domain"
	"go-movie-api/utils"
	"go-movie-api/utils/helper"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ratingRepository struct {
	db *gorm.DB
}

func NewRatingRepository(gormDB *gorm.DB) domain.RatingRepository {
	return &ratingRepository{db: gormDB}
}

func (repo *ratingRepository) FindByID(ctx context.Context, uuid uuid.UUID) (domain.Rating, error) {
	var rating domain.Rating

	result := repo.db.WithContext(ctx).
		Preload("User").
		Preload("Movie").
		Where("uuid = ?", uuid.String()).
		First(&rating)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return domain.Rating{}, helper.NotFoundErr
		}
		return domain.Rating{}, result.Error
	}

	return rating, nil
}

func (repo *ratingRepository) FindByIDForUpdate(ctx context.Context, uuid uuid.UUID) (domain.Rating, error) {
	var rating domain.Rating

	result := repo.db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("uuid = ?", uuid.String()).
		First(&rating)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return domain.Rating{}, helper.NotFoundErr
		}
		utils.Logger.Error(result.Error.Error())
		return domain.Rating{}, result.Error
	}

	return rating, nil
}

func (repo *ratingRepository) Store(ctx context.Context, rating *domain.Rating) (domain.Rating, error) {
	var movie domain.Movie
	result := repo.db.WithContext(ctx).Where("uuid = ?", rating.Movie.Uuid.String()).First(&movie)
	if result.Error != nil {
		utils.Logger.Error(result.Error.Error())
		return domain.Rating{}, result.Error
	}

	rating.MovieID = movie.ID
	rating.Movie = &movie
	result = repo.db.WithContext(ctx).Clauses(clause.Returning{}).Omit("uuid").Create(rating)
	if result.Error != nil {
		utils.Logger.Error(result.Error.Error())
		return domain.Rating{}, result.Error
	}

	return *rating, nil
}

func (repo *ratingRepository) Update(ctx context.Context, rating *domain.Rating) error {
	result := repo.db.WithContext(ctx).Model(rating).Where("uuid = ?", rating.Uuid.String()).Updates(rating)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *ratingRepository) SoftDelete(ctx context.Context, uuid uuid.UUID) error {
	result := repo.db.WithContext(ctx).Where("uuid = ?", uuid.String()).Delete(&domain.Rating{})
	if result.Error != nil {
		utils.Logger.Error(result.Error.Error())
		return result.Error
	}

	return nil
}

func (repo *ratingRepository) Delete(ctx context.Context, uuid uuid.UUID) error {
	result := repo.db.WithContext(ctx).Unscoped().Where("uuid = ?", uuid.String()).Delete(&domain.Rating{})
	if result.Error != nil {
		utils.Logger.Error(result.Error.Error())
		return result.Error
	}

	return nil
}
