package main

import (
	"context"
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
	PORT      = "7000"
	HOST      = "0.0.0.0"
	INTERVAL  = "3m" // time interval for data fetching from source
	TIMESCALE = "postgresql://admin:admin@localhost:5432/weather?sslmode=disable"
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

	if os.Getenv("POSTGRES") != "" {
		TIMESCALE = os.Getenv("POSTGRES")
	}
}

func main() {
	ctx := context.Background()
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

	// Use signals to gracefully shut down all the running go routines
	// and clear resource
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// concurrently fetch data at the app.interval interval
	go app.dataSourceFetcher(ctx)
}

func crashWithError(msg string, err error) {
	slog.Warn(msg)
	slog.Error(err.Error())
	os.Exit(1)
}
