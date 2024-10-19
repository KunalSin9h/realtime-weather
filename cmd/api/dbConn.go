package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log/slog"
	"time"
)

const dbName = "TimescaleDB"

// setupTimescaleDb connect to db
func setupTimescaleDb(ctx context.Context, connString string) (*pgx.Conn, error) {
	slog.Info(fmt.Sprintf("Connecting to %s...", dbName), "connection string", connString)

	for i := 1; i <= 5; i++ {
		conn, err := pgx.Connect(ctx, connString)

		if err != nil {
			slog.Warn(fmt.Sprintf("Failed to connect to %s, retrying...[%d/5]", dbName, i))

			backOff := i * 2
			time.Sleep(time.Duration(backOff) * time.Second)
			continue
		}

		slog.Info(fmt.Sprintf("Successfully connected to  %s", dbName))

		return conn, err
	}

	slog.Error(fmt.Sprintf("Failed to connect to %s", dbName))
	return nil, fmt.Errorf("failed to connect to %s, exiting", dbName)
}
