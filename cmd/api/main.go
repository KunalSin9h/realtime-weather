package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Config struct {
	dbConn   *pgx.Conn
	interval time.Duration
}

var (
	PORT                 = "7000"
	HOST                 = "0.0.0.0"
	INTERVAL             = "3m" // time interval for data fetching from source
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
	// Use signals to gracefully shut down all the running go routines
	// and clear resource
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	var app Config

	// parse interval from env
	interval, err := time.ParseDuration(INTERVAL)
	if err != nil {
		crashWithError("Invalid interval format", err)
	}

	// set the interval to application config
	app.interval = interval

	db, err := setupTimescaleDb(ctx, TIMESCALE)
	if err != nil {
		crashWithError("Failed to connect to database", err)
	}

	defer db.Close(ctx)

	app.dbConn = db

	// concurrently fetch data at the app.interval interval
	app.dataSourceFetcher(ctx, OPEN_WEATHER_API_KEY)
}

func crashWithError(msg string, err error) {
	slog.Warn(msg)
	slog.Error(err.Error())
	os.Exit(1)
}
