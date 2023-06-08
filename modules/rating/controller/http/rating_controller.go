package http

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go-movie-api/domain"
	"go-movie-api/middleware"
	"go-movie-api/utils/response"
	"net/http"
)

type RatingController struct {
	domain.RatingService
}

func NewRatingController(router *echo.Echo, genreService domain.RatingService) {
	controller := &RatingController{
		RatingService: genreService,
	}

	group := router.Group("/ratings")
	group.GET("/:uuid", controller.Show)
	group.POST("", controller.Store, middleware.AuthMiddleware.Handler)
	group.PUT("/:uuid", controller.Update, middleware.AuthMiddleware.Handler)
	group.DELETE("/:uuid", controller.Destroy, middleware.AuthMiddleware.Handler)
}

func (controller *RatingController) Show(ec echo.Context) error {
	id, err := uuid.Parse(ec.Param("uuid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "the id is not valid.")
	}

	data, err := controller.RatingService.FindByID(ec.Request().Context(), id)
	if err != nil {
		return err
	}

	return ec.JSON(http.StatusOK, data)
}

func (controller *RatingController) Store(ec echo.Context) error {
	var request storeRequest
	if err := ec.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	if err := ec.Validate(request); err != nil {
		return err
	}

	data, err := controller.RatingService.Store(ec.Request().Context(), &domain.Rating{
		Rating:  request.Rating,
		Comment: request.Comment,
		Movie: &domain.Movie{
			Uuid: request.MovieUuid,
		},
	})
	if err != nil {
		return err
	}

	return ec.JSON(http.StatusCreated, data)
}

func (controller *RatingController) Update(ec echo.Context) error {
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

	err = controller.RatingService.Update(ec.Request().Context(), &domain.Rating{
		Uuid:    id,
		Rating:  request.Rating,
		Comment: request.Comment,
	})
	if err != nil {
		return err
	}

	return ec.JSON(http.StatusOK, response.UpdateSuccess)
}

func (controller *RatingController) Destroy(ec echo.Context) error {
	id, err := uuid.Parse(ec.Param("uuid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "the id is not valid.")
	}

	err = controller.RatingService.SoftDelete(ec.Request().Context(), id)
	if err != nil {
		return err
	}

	return ec.JSON(http.StatusOK, response.DeleteSuccess)
}
