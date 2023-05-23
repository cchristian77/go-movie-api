package repository

import (
	"context"
	"github.com/google/uuid"
	"go-movie-api/domain"
	"go-movie-api/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type movieRepository struct {
	db *gorm.DB
}

func NewMovieRepository(gormDB *gorm.DB) domain.MovieRepository {
	return &movieRepository{db: gormDB}
}

func (repo *movieRepository) FetchPagination(ctx context.Context, pagination *utils.Pagination) ([]domain.Movie, error) {
	var movies []domain.Movie

	result := repo.db.WithContext(ctx).
		Scopes(utils.Paginate(movies, pagination, repo.db)).
		Preload("Genres").
		Order("id asc").
		Find(&movies)
	if result.Error != nil {
		utils.Logger.Error(result.Error.Error())
		return nil, result.Error
	}

	return movies, nil
}

func (repo *movieRepository) GetByID(ctx context.Context, uuid uuid.UUID) (domain.Movie, error) {
	var movie domain.Movie

	result := repo.db.WithContext(ctx).Where("uuid = ?", uuid.String()).First(&movie)
	if result.Error != nil {
		return domain.Movie{}, result.Error
	}

	return movie, nil
}

func (repo *movieRepository) Store(ctx context.Context, movie *domain.Movie) (domain.Movie, error) {
	result := repo.db.WithContext(ctx).Clauses(clause.Returning{}).Omit("uuid").Create(&movie)
	if result.Error != nil {
		utils.Logger.Error(result.Error.Error())
		return domain.Movie{}, result.Error
	}

	return *movie, nil
}

func (repo *movieRepository) Update(ctx context.Context, movie *domain.Movie) error {
	//result := repo.db.WithContext(ctx).
	//	Clauses(clause.Locking{Strength: "UPDATE"}).
	//	Where("uuid = ?", movie.Uuid).
	//	First(&movie)
	//if result.Error != nil {
	//	if result.Error == gorm.ErrRecordNotFound {
	//		return helper.NotFoundErr
	//	}
	//	utils.Logger.Error(result.Error.Error())
	//	return result.Error
	//}

	result := repo.db.WithContext(ctx).Model(movie).Where("uuid = ?", movie.Uuid.String()).Updates(movie)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *movieRepository) SoftDelete(ctx context.Context, uuid uuid.UUID) error {
	result := repo.db.WithContext(ctx).Where("uuid = ?", uuid.String()).Delete(&domain.Movie{})
	if result.Error != nil {
		utils.Logger.Error(result.Error.Error())
		return result.Error
	}

	return nil
}

func (repo *movieRepository) Delete(ctx context.Context, uuid uuid.UUID) error {
	result := repo.db.WithContext(ctx).Unscoped().Where("uuid = ?", uuid.String()).Delete(&domain.Movie{})
	if result.Error != nil {
		utils.Logger.Error(result.Error.Error())
		return result.Error
	}

	return nil
}
