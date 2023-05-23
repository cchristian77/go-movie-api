package service

import (
	"context"
	"github.com/google/uuid"
	"go-movie-api/domain"
	"go-movie-api/utils"
	errors "go-movie-api/utils/helper"
	"gorm.io/gorm"
	"time"
)

type genreService struct {
	genreRepo domain.GenreRepository
	timeout   time.Duration
}

func NewGenreService(genreRepo domain.GenreRepository, timeout time.Duration) domain.GenreService {
	return &genreService{
		genreRepo: genreRepo,
		timeout:   timeout,
	}
}

func (service *genreService) FetchPagination(ctx context.Context, page int, perPage int) ([]domain.Genre, utils.Pagination, error) {
	ctx, cancel := context.WithTimeout(ctx, service.timeout)
	defer cancel()

	pagination := utils.Pagination{
		Page:    page,
		PerPage: perPage,
	}
	genres, err := service.genreRepo.FetchPagination(ctx, &pagination)
	if err != nil {
		return nil, utils.Pagination{}, err
	}

	return genres, pagination, nil
}

func (service *genreService) GetByID(ctx context.Context, uuid uuid.UUID) (domain.Genre, error) {
	ctx, cancel := context.WithTimeout(ctx, service.timeout)
	defer cancel()

	genre, err := service.genreRepo.GetByID(ctx, uuid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.Genre{}, errors.NotFoundErr
		}
		return domain.Genre{}, err
	}

	return genre, nil
}

func (service *genreService) Store(ctx context.Context, genre *domain.Genre) (domain.Genre, error) {
	ctx, cancel := context.WithTimeout(ctx, service.timeout)
	defer cancel()

	result, err := service.genreRepo.Store(ctx, genre)
	if err != nil {
		return domain.Genre{}, err
	}

	return result, nil
}

func (service *genreService) Update(ctx context.Context, genre *domain.Genre) error {
	ctx, cancel := context.WithTimeout(ctx, service.timeout)
	defer cancel()

	if _, err := service.genreRepo.GetByID(ctx, genre.Uuid); err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NotFoundErr
		}
		return err
	}

	if err := service.genreRepo.Update(ctx, genre); err != nil {
		return err
	}

	return nil
}

func (service *genreService) SoftDelete(ctx context.Context, uuid uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, service.timeout)
	defer cancel()

	if _, err := service.genreRepo.GetByID(ctx, uuid); err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NotFoundErr
		}
		return err
	}

	if err := service.genreRepo.SoftDelete(ctx, uuid); err != nil {
		return err
	}

	return nil
}

func (service *genreService) Delete(ctx context.Context, uuid uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, service.timeout)
	defer cancel()

	if _, err := service.genreRepo.GetByID(ctx, uuid); err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NotFoundErr
		}
		return err
	}

	if err := service.genreRepo.Delete(ctx, uuid); err != nil {
		return err
	}

	return nil
}
