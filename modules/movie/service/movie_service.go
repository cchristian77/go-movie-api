package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go-movie-api/domain"
	"go-movie-api/utils"
	errors "go-movie-api/utils/helper"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type movieService struct {
	movieRepo domain.MovieRepository
	genreRepo domain.GenreRepository
	timeout   time.Duration
}

func NewMovieService(movieRepo domain.MovieRepository, genreRepo domain.GenreRepository, timeout time.Duration) domain.MovieService {
	return &movieService{
		movieRepo: movieRepo,
		genreRepo: genreRepo,
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

func (service *movieService) FindByID(ctx context.Context, uuid uuid.UUID) (domain.Movie, error) {
	ctx, cancel := context.WithTimeout(ctx, service.timeout)
	defer cancel()

	movie, err := service.movieRepo.FindByID(ctx, uuid)
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

	var ids []uuid.UUID
	for _, genre := range movie.Genres {
		ids = append(ids, genre.Uuid)
	}
	genres, err := service.genreRepo.FindByIDs(ctx, ids)
	if err != nil {
		return domain.Movie{}, nil
	}
	if len(ids) != len(genres) {
		return domain.Movie{}, echo.NewHTTPError(http.StatusBadRequest, "The genre(s) is not valid.")
	}

	movie.Genres = genres
	result, err := service.movieRepo.Store(ctx, movie)
	if err != nil {
		return domain.Movie{}, err
	}

	return result, nil
}

func (service *movieService) Update(ctx context.Context, movie *domain.Movie) error {
	ctx, cancel := context.WithTimeout(ctx, service.timeout)
	defer cancel()

	result, err := service.movieRepo.FindByIDForUpdate(ctx, movie.Uuid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NotFoundErr
		}
		return err
	}
	movie.ID = result.ID

	var ids []uuid.UUID
	for _, genre := range movie.Genres {
		ids = append(ids, genre.Uuid)
	}
	genres, err := service.genreRepo.FindByIDs(ctx, ids)
	if err != nil {
		return nil
	}
	if len(ids) != len(genres) {
		return echo.NewHTTPError(http.StatusBadRequest, "The genre(s) is not valid.")
	}

	movie.Genres = genres
	if err = service.movieRepo.Update(ctx, movie); err != nil {
		return err
	}

	return nil
}

func (service *movieService) SoftDelete(ctx context.Context, uuid uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, service.timeout)
	defer cancel()

	if _, err := service.movieRepo.FindByIDForUpdate(ctx, uuid); err != nil {
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

func (service *movieService) Delete(ctx context.Context, uuid uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, service.timeout)
	defer cancel()

	if _, err := service.movieRepo.FindByIDForUpdate(ctx, uuid); err != nil {
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
