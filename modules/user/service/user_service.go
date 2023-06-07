package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go-movie-api/domain"
	"go-movie-api/utils"
	errorHelper "go-movie-api/utils/helper"
	"gorm.io/gorm"
	"time"
)

type userService struct {
	userRepo    domain.UserRepository
	sessionRepo domain.SessionRepository
	timeout     time.Duration
}

func NewUserService(userRepo domain.UserRepository, sessionRepo domain.SessionRepository, timeout time.Duration) domain.UserService {
	return &userService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		timeout:     timeout,
	}
}

func (service *userService) Login(ctx context.Context, user *domain.User) (domain.User, error) {
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
		return domain.User{}, errorHelper.IncorrectCredentialErr
	}

	return authUser, nil
}

func (service *userService) FetchPagination(ctx context.Context, page int, perPage int) ([]domain.User, utils.Pagination, error) {
	ctx, cancel := context.WithTimeout(ctx, service.timeout)
	defer cancel()

	pagination := utils.Pagination{
		Page:    page,
		PerPage: perPage,
	}
	users, err := service.userRepo.FetchPagination(ctx, &pagination)
	if err != nil {
		return nil, utils.Pagination{}, err
	}

	return users, pagination, nil
}

func (service *userService) FindByID(ctx context.Context, uuid uuid.UUID) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, service.timeout)
	defer cancel()

	user, err := service.userRepo.FindByID(ctx, uuid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.User{}, errorHelper.NotFoundErr
		}
		return domain.User{}, err
	}

	return user, nil
}

func (service *userService) Store(ctx context.Context, user *domain.User) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, service.timeout)
	defer cancel()

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return domain.User{}, errors.New(fmt.Sprintf("failed to hash password: %s", err))
	}
	user.Password = hashedPassword

	result, err := service.userRepo.Store(ctx, user)
	if err != nil {
		return domain.User{}, err
	}

	return result, nil
}

func (service *userService) Update(ctx context.Context, user *domain.User) error {
	ctx, cancel := context.WithTimeout(ctx, service.timeout)
	defer cancel()

	_, err := service.userRepo.FindByIDForUpdate(ctx, user.Uuid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errorHelper.NotFoundErr
		}
		return err
	}

	if err = service.userRepo.Update(ctx, user); err != nil {
		return err
	}

	return nil
}

func (service *userService) SoftDelete(ctx context.Context, uuid uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, service.timeout)
	defer cancel()

	if _, err := service.userRepo.FindByIDForUpdate(ctx, uuid); err != nil {
		if err == gorm.ErrRecordNotFound {
			return errorHelper.NotFoundErr
		}
		return err
	}

	if err := service.userRepo.SoftDelete(ctx, uuid); err != nil {
		return err
	}

	return nil
}

func (service *userService) Delete(ctx context.Context, uuid uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, service.timeout)
	defer cancel()

	if _, err := service.userRepo.FindByIDForUpdate(ctx, uuid); err != nil {
		if err == gorm.ErrRecordNotFound {
			return errorHelper.NotFoundErr
		}
		return err
	}

	if err := service.userRepo.Delete(ctx, uuid); err != nil {
		return err
	}

	return nil
}

func (service *userService) CreateSession(ctx context.Context, session *domain.Session) (domain.Session, error) {
	ctx, cancel := context.WithTimeout(ctx, service.timeout)
	defer cancel()

	result, err := service.sessionRepo.Store(ctx, session)
	if err != nil {
		return domain.Session{}, err
	}

	return result, nil
}
