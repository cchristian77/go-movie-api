package main

import (
	"database/sql"
	"fmt"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/labstack/echo/v4"
	"go-movie-api/api"
	"go-movie-api/configs"
	"go-movie-api/database"
	"go-movie-api/utils"
	"gorm.io/gorm"
	"log"
)

var config = koanf.New(".")

type Server struct {
	DB     *sql.DB
	GormDB *gorm.DB
	Router *echo.Echo
}

func init() {
	// Load Config JSON
	if err := config.Load(file.Provider("./configs/env.json"), json.Parser()); err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	if err := config.UnmarshalWithConf("env", &configs.Env, koanf.UnmarshalConf{Tag: "koanf"}); err != nil {
		utils.Logger.Fatal(fmt.Sprintf("failed to read env.json file: %v", err))
	}

	log.Println("Starting service on port", configs.Env.App.Port)
}

func main() {
	// Initialized Logger
	utils.Logger = utils.InitializedLogger()
	defer utils.Logger.Sync()

	db := database.ConnectToDB()
	if db == nil {
		utils.Logger.Fatal("Can't connect to Postgres!")
	}

	gormDB, err := database.OpenGormDB(db)
	if err != nil {
		utils.Logger.Fatal(fmt.Sprintf("gorm driver errror: %v", err))
	}

	app := Server{
		DB:     db,
		GormDB: gormDB,
		Router: api.InitializedRouter(gormDB),
	}
	defer app.DB.Close()

	// Run application
	app.Router.Logger.Fatal(app.Router.Start(fmt.Sprintf(":%d", configs.Env.App.Port)))
}
