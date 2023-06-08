package middleware

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"go-movie-api/domain"
	"go-movie-api/token"
	"go-movie-api/utils/helper"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
)

const (
	authHeaderKey  = "authorization"
	authTypeBearer = "bearer"
	AuthPayloadKey = "auth_payload"
	AuthUserKey    = "auth_user"
)

var AuthMiddleware *authMiddleware

type authMiddleware struct {
	userRepo    domain.UserRepository
	sessionRepo domain.SessionRepository
}

func NewAuthMiddleware(sessionRepo domain.SessionRepository, userRepo domain.UserRepository) *authMiddleware {
	return &authMiddleware{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

func (middleware *authMiddleware) Handler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ec echo.Context) error {
		authHeader := ec.Request().Header.Get(authHeaderKey)
		ctx := ec.Request().Context()

		if authHeader == "" {
			return helper.UnauthorizedErr
		}

		authFields := strings.Fields(authHeader)
		if len(authFields) < 2 {
			return helper.UnauthorizedErr
		}

		authorizationType := strings.ToLower(authFields[0])
		if authorizationType != authTypeBearer {
			return echo.NewHTTPError(http.StatusUnauthorized, fmt.Sprintf("unsupported authorization type %s", authorizationType))
		}

		accessToken := authFields[1]
		payload, err := token.TokenMaker.VerifyToken(accessToken)
		if err != nil {
			return err
		}

		session, err := middleware.sessionRepo.FindByID(ctx, payload.ID)
		if err != nil {
			return token.InvalidTokenErr
		}

		if session.AccessToken != accessToken {
			return token.InvalidTokenErr
		}

		if time.Now().After(session.AccessTokenExpiresAt) {
			return token.InvalidTokenErr
		}

		if session.IsRevoked {
			return token.InvalidTokenErr
		}

		user, err := middleware.userRepo.FindByID(ctx, payload.UserUuid)
		if session.IsRevoked {
			if err == gorm.ErrRecordNotFound {
				return echo.NewHTTPError(http.StatusNotFound, "user not found")
			}

			return err
		}

		if session.UserID != user.ID {
			return token.InvalidTokenErr
		}

		ec.Set(AuthPayloadKey, payload)
		ec.Set(AuthUserKey, &user)

		return next(ec)
	}
}
