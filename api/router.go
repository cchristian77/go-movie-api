package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go-movie-api/utils"
	"gorm.io/gorm"
	"net/http"
	"time"

	_genreController "go-movie-api/modules/genre/controller/http"
	_genreRepo "go-movie-api/modules/genre/repository"
	_genreService "go-movie-api/modules/genre/service"
	_movieController "go-movie-api/modules/movie/controller/http"
	_movieRepo "go-movie-api/modules/movie/repository"
	_movieService "go-movie-api/modules/movie/service"
)

func InitializedRouter(db *gorm.DB) *echo.Echo {
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

	// Register API Routes
	registerRoutes(router, db)

	return router
}

func registerRoutes(router *echo.Echo, db *gorm.DB) {

	timeout := 5 * time.Second

	// Movies
	movieRepo := _movieRepo.NewMovieRepository(db)
	movieService := _movieService.NewMovieService(movieRepo, timeout)
	_movieController.NewMovieController(router, movieService)

	// Genre
	genreRepo := _genreRepo.NewGenreRepository(db)
	genreService := _genreService.NewGenreService(genreRepo, timeout)
	_genreController.NewGenreController(router, genreService)
}
