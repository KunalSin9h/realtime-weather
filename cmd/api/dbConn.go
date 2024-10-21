package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"time"
)

const dbName = "TimescaleDB"

// setupTimescaleDb connect to db
func setupTimescaleDb(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	slog.Info(fmt.Sprintf("Connecting to %s...", dbName), "connection string", connString)

	for i := 1; i <= 5; i++ {
		pool, err := pgxpool.New(ctx, connString)

		if err != nil {
			slog.Warn(fmt.Sprintf("Failed to connect to %s, retrying...[%d/5]", dbName, i))

			backOff := i * 2
			time.Sleep(time.Duration(backOff) * time.Second)
			continue
		}

		slog.Info(fmt.Sprintf("Successfully connected to  %s", dbName))

		return pool, err
	}

	slog.Error(fmt.Sprintf("Failed to connect to %s", dbName))
	return nil, fmt.Errorf("failed to connect to %s, exiting", dbName)
}
