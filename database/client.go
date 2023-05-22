package database

import (
	"database/sql"
	"fmt"
	"go-movie-api/utils"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var counts int64

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func ConnectToDB() *sql.DB {
	// get dsn from env.json
	//dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
	//	config.String("database.host"),
	//	config.Int("database.port"),
	//	config.String("database.user"),
	//	config.String("database.password"),
	//	config.String("database.db_name"),
	//)

	dsn := os.Getenv("DSN")
	fmt.Println(dsn)
	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready...")
			counts++
		} else {
			log.Println("Connected to Postgres")
			return connection
		}

		if counts > 10 {
			utils.Logger.Error(err.Error())
			return nil
		}

		log.Println("backing off for two seconds...")
		time.Sleep(2 * time.Second)
		continue
	}
}
