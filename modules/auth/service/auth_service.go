package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go-movie-api/domain"
	"go-movie-api/token"
	"go-movie-api/utils"
	errorHelper "go-movie-api/utils/helper"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type authService struct {
	userRepo    domain.UserRepository
	sessionRepo domain.SessionRepository
	timeout     time.Duration
}

func NewAuthService(userRepo domain.UserRepository, sessionRepo domain.SessionRepository, timeout time.Duration) domain.AuthService {
	return &authService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		timeout:     timeout,
	}
}

func (service *authService) Authenticate(ctx context.Context, user *domain.User) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, service.timeout)
	defer cancel()

	authUser, err := service.userRepo.FindByUsernameOrEmail(ctx, user.Username, user.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.User{}, errorHelper.IncorrectCredentialErr
		}

		return domain.User{}, err
	}

	err = utils.CheckPassword(user.Password, authUser.Password)
	if err != nil {
		return domain.User{}, err
	}

	return authUser, nil
}

func (service *authService) CreateSession(ctx context.Context, session *domain.Session) (domain.Session, error) {
	ctx, cancel := context.WithTimeout(ctx, service.timeout)
	defer cancel()

	result, err := service.sessionRepo.Store(ctx, session)
	if err != nil {
		return domain.Session{}, err
	}

	return result, nil
}

func (service *authService) VerifySession(ctx context.Context, payload *token.Payload, refreshToken string) error {
	session, err := service.sessionRepo.FindByID(ctx, payload.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errorHelper.NotFoundErr
		}

		return err
	}

	if session.RefreshToken != refreshToken {
		return token.InvalidTokenErr
	}

	if time.Now().After(session.ExpiresAt) {
		return echo.NewHTTPError(http.StatusUnauthorized, errors.New("session is expired"))
	}

	user, err := service.userRepo.FindByID(ctx, payload.UserUuid)
	if session.IsBlocked {
		if err == gorm.ErrRecordNotFound {
			return errorHelper.NotFoundErr
		}

		return err
	}

	if session.UserID != user.ID {
		return token.InvalidTokenErr
	}

	return nil
}

func (service *authService) BlockSession(ctx context.Context, sessionID uuid.UUID) error {
	fmt.Println(sessionID)
	session, err := service.sessionRepo.FindByID(ctx, sessionID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errorHelper.NotFoundErr
		}

		return err
	}

	if err = service.sessionRepo.BlockSession(ctx, session.ID); err != nil {
		return err
	}

	return nil
}
