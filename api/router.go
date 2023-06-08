package api

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go-movie-api/configs"
	m "go-movie-api/middleware"
	_authController "go-movie-api/modules/auth/controller/http"
	_authService "go-movie-api/modules/auth/service"
	_sessionRepo "go-movie-api/modules/session/repository"
	_userController "go-movie-api/modules/user/controller/http"
	_userRepo "go-movie-api/modules/user/repository"
	_userService "go-movie-api/modules/user/service"
	"go-movie-api/token"
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
	_ratingController "go-movie-api/modules/rating/controller/http"
	_ratingRepo "go-movie-api/modules/rating/repository"
	_ratingService "go-movie-api/modules/rating/service"
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

	// Register API Routes
	registerRoutes(router, db)

	return router
}

func registerRoutes(router *echo.Echo, db *gorm.DB) {
	timeout, _ := time.ParseDuration(configs.Env.Context.Timeout)
	router.GET("/ping", func(ec echo.Context) error {
		return ec.JSON(http.StatusOK, map[string]string{
			"message": "Ping!",
		})
	})

	tokenMaker, err := token.NewJWTMaker(configs.Env.JWTKey)
	if err != nil {
		utils.Logger.Fatal(fmt.Sprintf("failed to create token maker: %s", err))
	}
	token.TokenMaker = tokenMaker

	userRepo := _userRepo.NewUserRepository(db)
	sessionRepo := _sessionRepo.NewSessionRepository(db)
	m.AuthMiddleware = m.NewAuthMiddleware(sessionRepo, userRepo)

	// User
	userService := _userService.NewUserService(userRepo, timeout)
	_userController.NewUserController(router, userService)

	// Auth
	authService := _authService.NewAuthService(userRepo, sessionRepo, timeout)
	_authController.NewAuthController(router, authService, userService)

	// Genre
	genreRepo := _genreRepo.NewGenreRepository(db)
	genreService := _genreService.NewGenreService(genreRepo, timeout)
	_genreController.NewGenreController(router, genreService)

	// Movies
	movieRepo := _movieRepo.NewMovieRepository(db)
	movieService := _movieService.NewMovieService(movieRepo, genreRepo, timeout)
	_movieController.NewMovieController(router, movieService)

	// Rating
	ratingRepo := _ratingRepo.NewRatingRepository(db)
	ratingService := _ratingService.NewRatingService(ratingRepo, movieRepo, timeout)
	_ratingController.NewRatingController(router, ratingService)
}
