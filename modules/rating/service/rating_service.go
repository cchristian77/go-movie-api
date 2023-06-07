package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go-movie-api/domain"
	errors "go-movie-api/utils/helper"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type ratingService struct {
	ratingRepo domain.RatingRepository
	movieRepo  domain.MovieRepository
	timeout    time.Duration
}

func NewRatingService(ratingRepo domain.RatingRepository, movieRepo domain.MovieRepository, timeout time.Duration) domain.RatingService {
	return &ratingService{
		ratingRepo: ratingRepo,
		movieRepo:  movieRepo,
		timeout:    timeout,
	}
}

func (service *ratingService) FindByID(ctx context.Context, uuid uuid.UUID) (domain.Rating, error) {
	ctx, cancel := context.WithTimeout(ctx, service.timeout)
	defer cancel()

	rating, err := service.ratingRepo.FindByID(ctx, uuid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.Rating{}, errors.NotFoundErr
		}
		return domain.Rating{}, err
	}

	return rating, nil
}

func (service *ratingService) Store(ctx context.Context, rating *domain.Rating) (domain.Rating, error) {
	ctx, cancel := context.WithTimeout(ctx, service.timeout)
	defer cancel()

	movie, err := service.movieRepo.FindByID(ctx, rating.Movie.Uuid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.Rating{}, echo.NewHTTPError(http.StatusBadRequest, "The movie is not valid.")
		}

		return domain.Rating{}, err
	}

	result, err := service.ratingRepo.Store(ctx, rating)
	if err != nil {
		return domain.Rating{}, err
	}
	result.Movie = &movie

	return result, nil
}

func (service *ratingService) Update(ctx context.Context, rating *domain.Rating) error {
	ctx, cancel := context.WithTimeout(ctx, service.timeout)
	defer cancel()

	if _, err := service.ratingRepo.FindByIDForUpdate(ctx, rating.Uuid); err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NotFoundErr
		}
		return err
	}

	if err := service.ratingRepo.Update(ctx, rating); err != nil {
		return err
	}

	return nil
}

func (service *ratingService) SoftDelete(ctx context.Context, uuid uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, service.timeout)
	defer cancel()

	if _, err := service.ratingRepo.FindByIDForUpdate(ctx, uuid); err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NotFoundErr
		}
		return err
	}

	if err := service.ratingRepo.SoftDelete(ctx, uuid); err != nil {
		return err
	}

	return nil
}

func (service *ratingService) Delete(ctx context.Context, uuid uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, service.timeout)
	defer cancel()

	if _, err := service.ratingRepo.FindByIDForUpdate(ctx, uuid); err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NotFoundErr
		}
		return err
	}

	if err := service.ratingRepo.Delete(ctx, uuid); err != nil {
		return err
	}

	return nil
}
