package http

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go-movie-api/configs"
	"go-movie-api/domain"
	"go-movie-api/middleware"
	"go-movie-api/token"
	"net/http"
	"time"
)

type AuthController struct {
	domain.AuthService
	domain.UserService
}

func NewAuthController(router *echo.Echo, authService domain.AuthService, userService domain.UserService) {
	controller := &AuthController{
		AuthService: authService,
		UserService: userService,
	}

	authGroup := router.Group("auth")
	authGroup.POST("/register", controller.Store)
	authGroup.POST("/login", controller.Login)
	authGroup.POST("/renew-token", controller.RenewAccessToken)
	authGroup.POST("/logout", controller.Logout, middleware.AuthMiddleware.Handler)
	authGroup.GET("/current-user", controller.CurrentUser, middleware.AuthMiddleware.Handler)
}

func (controller *AuthController) Store(ec echo.Context) error {
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

func (controller *AuthController) Login(ec echo.Context) error {
	var request loginRequest
	if err := ec.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	if err := ec.Validate(request); err != nil {
		return err
	}

	ctx := ec.Request().Context()
	user, err := controller.AuthService.Authenticate(ctx, &domain.User{
		Username: request.Username,
		Password: request.Password,
	})
	if err != nil {
		return err
	}

	tokenID := uuid.Must(uuid.NewRandom())

	accessTokenDuration, _ := time.ParseDuration(configs.Env.Auth.AccessTokenExpiration)
	refreshTokenDuration, _ := time.ParseDuration(configs.Env.Auth.RefreshTokenExpiration)
	accessToken, accessPayload, err := token.TokenMaker.GenerateToken(
		tokenID,
		user.Uuid,
		accessTokenDuration,
	)
	if err != nil {
		return err
	}

	refreshToken, refreshPayload, err := token.TokenMaker.GenerateToken(
		tokenID,
		user.Uuid,
		refreshTokenDuration,
	)
	if err != nil {
		return err
	}

	accessTokenExpiresAt := time.Unix(accessPayload.StandardClaims.ExpiresAt, 0)
	accessTokenCreatedAt := time.Unix(accessPayload.StandardClaims.IssuedAt, 0)
	refreshTokenExpiresAt := time.Unix(refreshPayload.StandardClaims.ExpiresAt, 0)
	refreshTokenCreatedAt := time.Unix(refreshPayload.StandardClaims.IssuedAt, 0)

	session := domain.Session{
		ID:                    tokenID,
		UserID:                user.ID,
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  refreshTokenExpiresAt,
		AccessTokenCreatedAt:  accessTokenCreatedAt,
		RefreshTokenExpiresAt: refreshTokenExpiresAt,
		RefreshTokenCreatedAt: refreshTokenCreatedAt,
		UserAgent:             ec.Request().UserAgent(),
		ClientIp:              ec.RealIP(),
		IsRevoked:             false,
	}

	controller.AuthService.DeleteOldSession(ctx, &session)
	createdSession, err := controller.AuthService.CreateSession(ctx, &session)
	if err != nil {
		return err
	}

	return ec.JSON(http.StatusOK, authResponse{
		SessionID:             createdSession.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessTokenExpiresAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshTokenExpiresAt,
		User: userResponse{
			Uuid:     user.Uuid,
			Username: user.Username,
			FullName: user.FullName,
			Email:    user.Email,
		},
	})
}

func (controller *AuthController) Logout(ec echo.Context) error {
	authPayload := ec.Get(middleware.AuthPayloadKey).(*token.Payload)
	err := controller.AuthService.RevokeSession(ec.Request().Context(), authPayload.ID)
	if err != nil {
		return err
	}

	return ec.JSON(http.StatusOK, map[string]string{"message": "Logout success !"})
}

func (controller *AuthController) RenewAccessToken(ec echo.Context) error {
	var request renewAccessTokenRequest
	if err := ec.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	if err := ec.Validate(request); err != nil {
		return err
	}

	refreshPayload, err := token.TokenMaker.VerifyToken(request.RefreshToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	err = controller.AuthService.VerifySession(
		ec.Request().Context(),
		refreshPayload,
		request.RefreshToken,
	)
	if err != nil {
		return err
	}

	ctx := ec.Request().Context()
	user, err := controller.UserService.FindByID(ctx, refreshPayload.UserUuid)
	if err != nil {
		return err
	}

	tokenID := uuid.Must(uuid.NewRandom())

	accessTokenDuration, _ := time.ParseDuration(configs.Env.Auth.AccessTokenExpiration)
	refreshTokenDuration, _ := time.ParseDuration(configs.Env.Auth.RefreshTokenExpiration)
	accessToken, accessPayload, err := token.TokenMaker.GenerateToken(
		tokenID,
		user.Uuid,
		accessTokenDuration,
	)
	if err != nil {
		return err
	}

	refreshToken, refreshPayload, err := token.TokenMaker.GenerateToken(
		tokenID,
		user.Uuid,
		refreshTokenDuration,
	)
	if err != nil {
		return err
	}

	accessTokenExpiresAt := time.Unix(accessPayload.StandardClaims.ExpiresAt, 0)
	accessTokenCreatedAt := time.Unix(accessPayload.StandardClaims.IssuedAt, 0)
	refreshTokenExpiresAt := time.Unix(refreshPayload.StandardClaims.ExpiresAt, 0)
	refreshTokenCreatedAt := time.Unix(refreshPayload.StandardClaims.IssuedAt, 0)

	session := domain.Session{
		ID:                    tokenID,
		UserID:                user.ID,
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  refreshTokenExpiresAt,
		AccessTokenCreatedAt:  accessTokenCreatedAt,
		RefreshTokenExpiresAt: refreshTokenExpiresAt,
		RefreshTokenCreatedAt: refreshTokenCreatedAt,
		UserAgent:             ec.Request().UserAgent(),
		ClientIp:              ec.RealIP(),
		IsRevoked:             false,
	}

	controller.AuthService.DeleteOldSession(ctx, &session)
	createdSession, err := controller.AuthService.CreateSession(ctx, &session)
	if err != nil {
		return err
	}

	return ec.JSON(http.StatusOK, authResponse{
		SessionID:             createdSession.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessTokenExpiresAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshTokenExpiresAt,
		User: userResponse{
			Uuid:     user.Uuid,
			Username: user.Username,
			FullName: user.FullName,
			Email:    user.Email,
		},
	})
}

func (controller *AuthController) CurrentUser(ec echo.Context) error {
	return ec.JSON(http.StatusOK, ec.Get(middleware.AuthUserKey).(*domain.User))
}
