package main

import (
	"fmt"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"go-movie-api/api"
	"go-movie-api/database"
	"go-movie-api/utils"
	"log"
)

var config = koanf.New(".")

func init() {
	// Load Config JSON
	if err := config.Load(file.Provider("./configs/env.json"), json.Parser()); err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	if err := config.UnmarshalWithConf("env", &utils.Env, koanf.UnmarshalConf{Tag: "koanf"}); err != nil {
		utils.Logger.Fatal(fmt.Sprintf("failed to read env.json file: %v", err))
	}

	log.Println("Starting service on port", utils.Env.App.Port)
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

	api.Server.DB = db
	api.Server.GormDB = gormDB
	api.Server.Router = api.InitializedRouter(gormDB)
	defer api.Server.DB.Close()

	// Run application
	api.Server.Router.Logger.Fatal(api.Server.Router.Start(fmt.Sprintf(":%d", utils.Env.App.Port)))
}
