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

type genreRepository struct {
	db *gorm.DB
}

func NewGenreRepository(gormDB *gorm.DB) domain.GenreRepository {
	return &genreRepository{db: gormDB}
}

func (repo *genreRepository) FetchPagination(ctx context.Context, pagination *utils.Pagination) ([]domain.Genre, error) {
	var genres []domain.Genre

	result := repo.db.WithContext(ctx).
		Scopes(utils.Paginate(genres, pagination, repo.db)).
		Order("id asc").
		Find(&genres)
	if result.Error != nil {
		utils.Logger.Error(result.Error.Error())
		return nil, result.Error
	}

	return genres, nil
}

func (repo *genreRepository) FindByID(ctx context.Context, uuid uuid.UUID) (domain.Genre, error) {
	var genre domain.Genre

	result := repo.db.WithContext(ctx).
		Preload("Movies").
		Where("uuid = ?", uuid.String()).
		First(&genre)
	if result.Error != nil {
		return domain.Genre{}, result.Error
	}

	return genre, nil
}

func (repo *genreRepository) FindByIDForUpdate(ctx context.Context, uuid uuid.UUID) (domain.Genre, error) {
	var genre domain.Genre

	result := repo.db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("uuid = ?", uuid.String()).
		First(&genre)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return domain.Genre{}, helper.NotFoundErr
		}
		utils.Logger.Error(result.Error.Error())
		return domain.Genre{}, result.Error
	}

	return genre, nil
}

func (repo *genreRepository) Store(ctx context.Context, genre *domain.Genre) (domain.Genre, error) {
	result := repo.db.WithContext(ctx).Clauses(clause.Returning{}).Omit("uuid").Create(&genre)
	if result.Error != nil {
		utils.Logger.Error(result.Error.Error())
		return domain.Genre{}, result.Error
	}

	return *genre, nil
}

func (repo *genreRepository) Update(ctx context.Context, genre *domain.Genre) error {
	result := repo.db.WithContext(ctx).Model(genre).Where("uuid = ?", genre.Uuid.String()).Updates(genre)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *genreRepository) SoftDelete(ctx context.Context, uuid uuid.UUID) error {
	result := repo.db.WithContext(ctx).Where("uuid = ?", uuid.String()).Delete(&domain.Genre{})
	if result.Error != nil {
		utils.Logger.Error(result.Error.Error())
		return result.Error
	}

	return nil
}

func (repo *genreRepository) Delete(ctx context.Context, uuid uuid.UUID) error {
	result := repo.db.WithContext(ctx).Unscoped().Where("uuid = ?", uuid.String()).Delete(&domain.Genre{})
	if result.Error != nil {
		utils.Logger.Error(result.Error.Error())
		return result.Error
	}

	return nil
}
