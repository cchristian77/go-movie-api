package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go-movie-api/utils"
	"net/http"
)

func InitializedRouter() *echo.Echo {
	router := echo.New()

	// Config CORS
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:          middleware.DefaultSkipper,
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, echo.HeaderXCSRFToken},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	router.Use(middleware.Recover())

	// Config Rate Limiter allows 100 requests/sec
	router.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(100)))

	// Config Validator to Router
	router.Validator = &utils.RequestValidator{Validator: validator.New()}

	// Register RequestLog to Router Middleware
	router.Use(utils.RequestLog)

	// Register HTTP Error Handler function
	router.HTTPErrorHandler = utils.ErrorHandler

	router.GET("/ping", func(ec echo.Context) error {
		return ec.JSON(http.StatusOK, map[string]string{
			"message": "Ping!",
		})
	})

	return router
}
