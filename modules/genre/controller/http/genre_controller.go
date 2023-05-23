package http

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go-movie-api/domain"
	"go-movie-api/utils/response"
	"net/http"
	"strconv"
)

type GenreController struct {
	domain.GenreService
}

func NewGenreController(router *echo.Echo, genreService domain.GenreService) {
	controller := &GenreController{
		GenreService: genreService,
	}

	group := router.Group("/genres")
	group.POST("", controller.Store)
	group.GET("", controller.Index)
	group.GET("/:uuid", controller.Show)
	group.PUT("/:uuid", controller.Update)
	group.DELETE("/:uuid", controller.Destroy)
}

func (controller *GenreController) Index(ec echo.Context) error {
	page, err := strconv.Atoi(ec.QueryParam("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	perPage, err := strconv.Atoi(ec.QueryParam("per_page"))
	if err != nil || perPage <= 0 {
		perPage = 100
	}

	ctx := ec.Request().Context()
	data, pagination, err := controller.GenreService.FetchPagination(ctx, page, perPage)
	if err != nil {
		return err
	}

	if data == nil {
		data = make([]domain.Genre, 0)
	}

	return ec.JSON(http.StatusOK, response.Result{
		Meta: pagination,
		Data: data,
	})
}

func (controller *GenreController) Show(ec echo.Context) error {
	id, err := uuid.Parse(ec.Param("uuid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "the id is not valid.")
	}

	data, err := controller.GenreService.FindByID(ec.Request().Context(), id)
	if err != nil {
		return err
	}

	return ec.JSON(http.StatusOK, data)
}

func (controller *GenreController) Store(ec echo.Context) error {
	var request storeRequest
	if err := ec.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	if err := ec.Validate(request); err != nil {
		return err
	}

	data, err := controller.GenreService.Store(ec.Request().Context(), &domain.Genre{
		Name: request.Name,
	})
	if err != nil {
		return err
	}

	return ec.JSON(http.StatusCreated, data)
}

func (controller *GenreController) Update(ec echo.Context) error {
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

	err = controller.GenreService.Update(ec.Request().Context(), &domain.Genre{
		Uuid: id,
		Name: request.Name,
	})
	if err != nil {
		return err
	}

	return ec.JSON(http.StatusOK, response.UpdateSuccess)
}

func (controller *GenreController) Destroy(ec echo.Context) error {
	id, err := uuid.Parse(ec.Param("uuid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "the id is not valid.")
	}

	err = controller.GenreService.SoftDelete(ec.Request().Context(), id)
	if err != nil {
		return err
	}

	return ec.JSON(http.StatusOK, response.DeleteSuccess)
}
