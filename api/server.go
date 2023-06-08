package api

import (
	"database/sql"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

var Server struct {
	DB     *sql.DB
	GormDB *gorm.DB
	Router *echo.Echo
}
