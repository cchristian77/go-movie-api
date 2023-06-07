package repository

import (
	"context"
	"github.com/google/uuid"
	"go-movie-api/domain"
	"go-movie-api/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type sessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(gormDB *gorm.DB) domain.SessionRepository {
	return &sessionRepository{db: gormDB}
}

func (repo *sessionRepository) Store(ctx context.Context, session *domain.Session) (domain.Session, error) {
	result := repo.db.WithContext(ctx).Clauses(clause.Returning{}).Create(&session)
	if result.Error != nil {
		utils.Logger.Error(result.Error.Error())
		return domain.Session{}, result.Error
	}

	return *session, nil
}

func (repo *sessionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := repo.db.WithContext(ctx).Unscoped().Where("uuid = ?", id.String()).Delete(&domain.Session{})
	if result.Error != nil {
		utils.Logger.Error(result.Error.Error())
		return result.Error
	}

	return nil
}
