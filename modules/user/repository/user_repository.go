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

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(gormDB *gorm.DB) domain.UserRepository {
	return &userRepository{db: gormDB}
}

func (repo *userRepository) FindByID(ctx context.Context, uuid uuid.UUID) (domain.User, error) {
	var user domain.User

	result := repo.db.WithContext(ctx).Where("uuid = ?", uuid.String()).First(&user)
	if result.Error != nil {
		return domain.User{}, result.Error
	}

	return user, nil
}

func (repo *userRepository) FindByUsernameOrEmail(ctx context.Context, username string, email string) (domain.User, error) {
	var user domain.User

	result := repo.db.WithContext(ctx).Where("username = ?", username).Or("email = ?", email).First(&user)
	if result.Error != nil {
		return domain.User{}, result.Error
	}

	return user, nil
}

func (repo *userRepository) FindByUsername(ctx context.Context, username string) (domain.User, error) {
	var user domain.User

	result := repo.db.WithContext(ctx).Where("username = ?", username).First(&user)
	if result.Error != nil {
		return domain.User{}, result.Error
	}

	return user, nil
}

func (repo *userRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	var user domain.User

	result := repo.db.WithContext(ctx).Where("email = ?", email).First(&user)
	if result.Error != nil {
		return domain.User{}, result.Error
	}

	return user, nil
}

func (repo *userRepository) FetchPagination(ctx context.Context, pagination *utils.Pagination) ([]domain.User, error) {
	var users []domain.User

	result := repo.db.WithContext(ctx).
		Scopes(utils.Paginate(users, pagination, repo.db)).
		Order("id asc").
		Find(&users)
	if result.Error != nil {
		utils.Logger.Error(result.Error.Error())
		return nil, result.Error
	}

	return users, nil
}

func (repo *userRepository) FindByIDForUpdate(ctx context.Context, uuid uuid.UUID) (domain.User, error) {
	var user domain.User

	result := repo.db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("uuid = ?", uuid.String()).
		First(&user)
	if result.Error != nil {
		utils.Logger.Error(result.Error.Error())

		if result.Error == gorm.ErrRecordNotFound {
			return domain.User{}, helper.NotFoundErr
		}

		return domain.User{}, result.Error
	}

	return user, nil
}

func (repo *userRepository) Store(ctx context.Context, user *domain.User) (domain.User, error) {
	result := repo.db.WithContext(ctx).Clauses(clause.Returning{}).Omit("uuid").Create(&user)
	if result.Error != nil {
		utils.Logger.Error(result.Error.Error())
		return domain.User{}, result.Error
	}

	return *user, nil
}

func (repo *userRepository) Update(ctx context.Context, user *domain.User) error {
	result := repo.db.WithContext(ctx).Model(user).Where("uuid = ?", user.Uuid.String()).Updates(user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *userRepository) SoftDelete(ctx context.Context, uuid uuid.UUID) error {
	result := repo.db.WithContext(ctx).Where("uuid = ?", uuid.String()).Delete(&domain.User{})
	if result.Error != nil {
		utils.Logger.Error(result.Error.Error())
		return result.Error
	}

	return nil
}

func (repo *userRepository) Delete(ctx context.Context, uuid uuid.UUID) error {
	result := repo.db.WithContext(ctx).Unscoped().Where("uuid = ?", uuid.String()).Delete(&domain.User{})
	if result.Error != nil {
		utils.Logger.Error(result.Error.Error())
		return result.Error
	}

	return nil
}
