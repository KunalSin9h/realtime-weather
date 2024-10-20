package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

type Config struct {
	dbConn   *pgx.Conn
	interval time.Duration
	UserPref *UserPreference
}

var (
	PORT                 = "5000"
	HOST                 = "0.0.0.0"
	OPEN_WEATHER_API_KEY = ""
	TIMESCALE            = "postgresql://admin:admin@localhost:5432/weather?sslmode=disable"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if os.Getenv("PORT") != "" {
		PORT = os.Getenv("PORT")
	}

	if os.Getenv("HOST") != "" {
		HOST = os.Getenv("HOST")
	}

	if os.Getenv("OPEN_WEATHER_API_KEY") != "" {
		OPEN_WEATHER_API_KEY = os.Getenv("OPEN_WEATHER_API_KEY")
	} else {
		crashWithError("Missing OPEN_WEATHER_API_KEY in env vars", fmt.Errorf(""))
	}

	if os.Getenv("POSTGRES") != "" {
		TIMESCALE = os.Getenv("POSTGRES")
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var app Config
	app.UserPref = getDefaultUserPreference()

	db, err := setupTimescaleDb(ctx, TIMESCALE)
	if err != nil {
		crashWithError("Failed to connect to database", err)
	}

	defer db.Close(ctx)

	app.dbConn = db

	// concurrently fetch data at the app.interval interval
	go app.dataSourceFetcher(ctx)

	// Setup and Run server
	// This will run Server on HOST:PORT
	log.Fatal(app.setUpAndRunServer())
}
