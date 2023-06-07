package http

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go-movie-api/domain"
	"go-movie-api/token"
	"go-movie-api/utils/response"
	"net/http"
	"strconv"
	"time"
)

type UserController struct {
	domain.UserService
	TokenMaker token.Maker
}

func NewUserController(router *echo.Echo, tokenMaker token.Maker, userService domain.UserService) {
	controller := &UserController{
		UserService: userService,
		TokenMaker:  tokenMaker,
	}

	group := router.Group("/users")
	group.GET("", controller.Index)
	group.GET("/:uuid", controller.Show)
	group.PUT("/:uuid", controller.Update)
	group.DELETE("/:uuid", controller.Destroy)

	authGroup := router.Group("auth")
	authGroup.POST("/login", controller.Login)
	authGroup.POST("/logout", controller.Logout)
	authGroup.POST("/register", controller.Store)
}

func (controller *UserController) Login(ec echo.Context) error {
	var request loginRequest
	if err := ec.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	if err := ec.Validate(request); err != nil {
		return err
	}

	ctx := ec.Request().Context()
	user, err := controller.UserService.Login(ctx, &domain.User{
		Username: request.Username,
		Password: request.Password,
	})
	if err != nil {
		return err
	}

	accessToken, accessPayload, err := controller.TokenMaker.GenerateToken(
		user.Username,
		15*time.Minute,
	)
	if err != nil {
		return err
	}

	refreshToken, refreshPayload, err := controller.TokenMaker.GenerateToken(
		user.Username,
		24*time.Hour,
	)
	if err != nil {
		return err
	}

	session, err := controller.UserService.CreateSession(ctx, &domain.Session{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    ec.Request().UserAgent(),
		ClientIp:     ec.RealIP(),
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		return err
	}

	return ec.JSON(http.StatusOK, authResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User: userResponse{
			Username: user.Username,
			FullName: user.FullName,
			Email:    user.Email,
		},
	})
}

func (controller *UserController) Logout(ec echo.Context) error {

	return nil
}

func (controller *UserController) Index(ec echo.Context) error {
	page, err := strconv.Atoi(ec.QueryParam("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	perPage, err := strconv.Atoi(ec.QueryParam("per_page"))
	if err != nil || perPage <= 0 {
		perPage = 100
	}

	ctx := ec.Request().Context()
	data, pagination, err := controller.UserService.FetchPagination(ctx, page, perPage)
	if err != nil {
		return err
	}

	if data == nil {
		data = make([]domain.User, 0)
	}

	return ec.JSON(http.StatusOK, response.Result{
		Meta: pagination,
		Data: data,
	})
}

func (controller *UserController) Show(ec echo.Context) error {
	id, err := uuid.Parse(ec.Param("uuid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "the id is not valid.")
	}

	data, err := controller.UserService.FindByID(ec.Request().Context(), id)
	if err != nil {
		return err
	}

	return ec.JSON(http.StatusOK, data)
}

func (controller *UserController) Store(ec echo.Context) error {
	var request storeRequest
	if err := ec.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	if err := ec.Validate(request); err != nil {
		return err
	}

	data, err := controller.UserService.Store(ec.Request().Context(), &domain.User{
		Username: request.Username,
		Email:    request.Email,
		FullName: request.FullName,
		Password: request.Password,
		IsAdmin:  request.IsAdmin,
	})
	if err != nil {
		return err
	}

	return ec.JSON(http.StatusCreated, data)
}

func (controller *UserController) Update(ec echo.Context) error {
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

	err = controller.UserService.Update(ec.Request().Context(), &domain.User{
		Uuid:     id,
		FullName: request.FullName,
		Username: request.Username,
		Email:    request.Email,
	})
	if err != nil {
		return err
	}

	return ec.JSON(http.StatusOK, response.UpdateSuccess)
}

func (controller *UserController) Destroy(ec echo.Context) error {
	id, err := uuid.Parse(ec.Param("uuid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "the id is not valid.")
	}

	err = controller.UserService.SoftDelete(ec.Request().Context(), id)
	if err != nil {
		return err
	}

	return ec.JSON(http.StatusOK, response.DeleteSuccess)
}
