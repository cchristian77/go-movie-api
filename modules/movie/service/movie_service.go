package service

import (
	"context"
	"go-movie-api/domain"
	"go-movie-api/utils"
	errors "go-movie-api/utils/helper"
	"gorm.io/gorm"
	"time"
)

type movieService struct {
	movieRepo domain.MovieRepository
	timeout   time.Duration
}

func NewMovieService(movieRepo domain.MovieRepository, timeout time.Duration) domain.MovieService {
	return &movieService{
		movieRepo: movieRepo,
		timeout:   timeout,
	}
}

func (service *movieService) FetchPagination(ctx context.Context, page int, perPage int) ([]domain.Movie, utils.Pagination, error) {
	ctx, cancel := context.WithTimeout(ctx, service.timeout)
	defer cancel()

	pagination := utils.Pagination{
		Page:    page,
		PerPage: perPage,
	}
	movies, err := service.movieRepo.FetchPagination(ctx, &pagination)
	if err != nil {
		return nil, utils.Pagination{}, err
	}

	return movies, pagination, nil
}

func (service *movieService) GetByID(ctx context.Context, uuid string) (domain.Movie, error) {
	ctx, cancel := context.WithTimeout(ctx, service.timeout)
	defer cancel()

	movie, err := service.movieRepo.GetByID(ctx, uuid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.Movie{}, errors.NotFoundErr
		}
		return domain.Movie{}, err
	}

	return movie, nil
}

func (service *movieService) Store(ctx context.Context, movie *domain.Movie) (domain.Movie, error) {
	ctx, cancel := context.WithTimeout(ctx, service.timeout)
	defer cancel()

	result, err := service.movieRepo.Store(ctx, movie)
	if err != nil {
		return domain.Movie{}, err
	}

	return result, nil
}

func (service *movieService) Update(ctx context.Context, movie *domain.Movie) error {
	ctx, cancel := context.WithTimeout(ctx, service.timeout)
	defer cancel()

	if err := service.movieRepo.Update(ctx, movie); err != nil {
		return err
	}

	return nil
}

func (service *movieService) SoftDelete(ctx context.Context, uuid string) error {
	ctx, cancel := context.WithTimeout(ctx, service.timeout)
	defer cancel()

	if _, err := service.movieRepo.GetByID(ctx, uuid); err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NotFoundErr
		}
		return err
	}

	if err := service.movieRepo.SoftDelete(ctx, uuid); err != nil {
		return err
	}

	return nil
}

func (service *movieService) Delete(ctx context.Context, uuid string) error {
	ctx, cancel := context.WithTimeout(ctx, service.timeout)
	defer cancel()

	if _, err := service.movieRepo.GetByID(ctx, uuid); err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NotFoundErr
		}
		return err
	}

	if err := service.movieRepo.Delete(ctx, uuid); err != nil {
		return err
	}

	return nil
}
