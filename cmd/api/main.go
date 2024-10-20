package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
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
	INTERVAL             = "1m" // time interval for data fetching from source
	OPEN_WEATHER_API_KEY = ""
	TIMESCALE            = "postgresql://admin:admin@localhost:5432/weather?sslmode=disable"
)

func init() {
	if os.Getenv("PORT") != "" {
		PORT = os.Getenv("PORT")
	}

	if os.Getenv("HOST") != "" {
		HOST = os.Getenv("HOST")
	}

	if os.Getenv("INTERVAL") != "" {
		INTERVAL = os.Getenv("INTERVAL")
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

	// parse interval from env
	interval, err := time.ParseDuration(INTERVAL)
	if err != nil {
		crashWithError("Invalid interval format", err)
	}

	// set the interval to application config
	app.interval = interval

	// Set User preference
	app.UserPref = &UserPreference{
		TempUnit: Celsius, // by default use Celsius for Temperature Unit
	}

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
