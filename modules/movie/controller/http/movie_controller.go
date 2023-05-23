package http

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go-movie-api/domain"
	"go-movie-api/utils/response"
	"net/http"
	"strconv"
)

type MovieController struct {
	domain.MovieService
}

func NewMovieController(router *echo.Echo, movieService domain.MovieService) {
	controller := &MovieController{
		MovieService: movieService,
	}

	group := router.Group("/movies")
	group.POST("", controller.Store)
	group.GET("", controller.Index)
	group.GET("/:uuid", controller.Show)
	group.PUT("/:uuid", controller.Update)
	group.DELETE("/:uuid", controller.Destroy)
}

func (controller *MovieController) Index(ec echo.Context) error {
	page, err := strconv.Atoi(ec.QueryParam("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	perPage, err := strconv.Atoi(ec.QueryParam("per_page"))
	if err != nil || perPage <= 0 {
		perPage = 100
	}

	ctx := ec.Request().Context()
	data, pagination, err := controller.MovieService.FetchPagination(ctx, page, perPage)
	if err != nil {
		return err
	}

	if data == nil {
		data = make([]domain.Movie, 0)
	}

	return ec.JSON(http.StatusOK, response.Result{
		Meta: pagination,
		Data: data,
	})
}

func (controller *MovieController) Show(ec echo.Context) error {
	id, err := uuid.Parse(ec.Param("uuid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "the id is not valid.")
	}

	data, err := controller.MovieService.FindByID(ec.Request().Context(), id)
	if err != nil {
		return err
	}

	return ec.JSON(http.StatusOK, data)
}

func (controller *MovieController) Store(ec echo.Context) error {
	var request storeRequest
	if err := ec.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	if err := ec.Validate(request); err != nil {
		return err
	}

	data, err := controller.MovieService.Store(ec.Request().Context(), &domain.Movie{
		Title:    request.Title,
		Duration: request.Duration,
		Year:     request.Year,
		Synopsis: request.Synopsis,
	})
	if err != nil {
		return err
	}

	return ec.JSON(http.StatusCreated, data)
}

func (controller *MovieController) Update(ec echo.Context) error {
	var request updateRequest
	if err := ec.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	if err := ec.Validate(request); err != nil {
		return err
	}

	id, err := uuid.Parse(ec.Param("uuid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "the id is not valid.")
	}

	err = controller.MovieService.Update(ec.Request().Context(), &domain.Movie{
		Uuid:     id,
		Title:    request.Title,
		Duration: request.Duration,
		Year:     request.Year,
		Synopsis: request.Synopsis,
	})
	if err != nil {
		return err
	}

	return ec.JSON(http.StatusOK, response.UpdateSuccess)
}

func (controller *MovieController) Destroy(ec echo.Context) error {
	id, err := uuid.Parse(ec.Param("uuid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "the id is not valid.")
	}

	err = controller.MovieService.SoftDelete(ec.Request().Context(), id)
	if err != nil {
		return err
	}

	return ec.JSON(http.StatusOK, response.DeleteSuccess)
}
