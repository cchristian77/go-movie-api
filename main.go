package main

import (
	"database/sql"
	"fmt"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/labstack/echo/v4"
	"go-movie-api/api"
	"go-movie-api/database"
	"go-movie-api/utils"
	"log"
)

var config = koanf.New(".")

type Config struct {
	DB     *sql.DB
	Router *echo.Echo
}

func init() {
	// Load Config JSON
	if err := config.Load(file.Provider("env.json"), json.Parser()); err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	log.Println("Starting service on port", config.String("app.port"))
}

func main() {
	// Initialized Logger
	utils.Logger = utils.InitializedLogger()
	defer utils.Logger.Sync()

	app := Config{
		DB:     database.ConnectToDB(),
		Router: api.InitializedRouter(),
	}
	defer app.DB.Close()

	// Run application
	app.Router.Logger.Fatal(app.Router.Start(fmt.Sprintf(":%s", config.String("app.port"))))
}
