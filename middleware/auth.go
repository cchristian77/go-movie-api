package middleware

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"go-movie-api/token"
	"go-movie-api/utils/helper"
	"net/http"
	"strings"
)

const (
	authHeaderKey  = "authorization"
	authTypeBearer = "bearer"
	AuthPayloadKey = "auth_payload"
	AuthUserKey    = "auth_user"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ec echo.Context) error {
		authHeader := ec.Request().Header.Get(authHeaderKey)

		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, helper.UnauthorizedErr)
		}

		authFields := strings.Fields(authHeader)
		if len(authFields) < 2 {
			return echo.NewHTTPError(http.StatusUnauthorized, helper.UnauthorizedErr)
		}

		authorizationType := strings.ToLower(authFields[0])
		if authorizationType != authTypeBearer {
			return echo.NewHTTPError(http.StatusUnauthorized, fmt.Sprintf("unsupported authorization type %s", authorizationType))
		}

		accessToken := authFields[1]
		payload, err := token.TokenMaker.VerifyToken(accessToken)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}

		ec.Set(AuthPayloadKey, payload)

		return next(ec)
	}
}
