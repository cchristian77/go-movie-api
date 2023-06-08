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

func (repo *sessionRepository) FindByID(ctx context.Context, id uuid.UUID) (domain.Session, error) {
	var session domain.Session

	result := repo.db.WithContext(ctx).First(&session, id)
	if result.Error != nil {
		utils.Logger.Error(result.Error.Error())
		return domain.Session{}, result.Error
	}

	return session, nil
}

func (repo *sessionRepository) BlockSession(ctx context.Context, id uuid.UUID) error {
	session := domain.Session{
		ID:        id,
		IsBlocked: true,
	}
	result := repo.db.WithContext(ctx).Model(&session).Updates(&session)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *sessionRepository) Delete(ctx context.Context, session *domain.Session) error {
	result := repo.db.WithContext(ctx).
		Where(
			"user_id = ? AND user_agent = ? AND client_ip = ?",
			session.UserID, session.UserAgent, session.ClientIp,
		).
		Delete(&domain.Session{})
	if result.Error != nil {
		utils.Logger.Error(result.Error.Error())
		return result.Error
	}

	return nil
}
