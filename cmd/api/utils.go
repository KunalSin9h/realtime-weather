package main

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kunalsin9h/realtime-weather/internal/db"
	"log/slog"
	"os"
	"time"
)

func calTimestampWithTZ(unix, timezone int64) pgtype.Timestamptz {
	// convert unix timestamp to time.Time
	timestamp := time.Unix(unix, 0)
	location := time.FixedZone("API Timezone", int(timezone))

	timeStampWithTimezone := timestamp.In(location)

	return pgtype.Timestamptz{
		Time:  timeStampWithTimezone,
		Valid: true,
	}
}

func getFloatLatLon(city *db.City) (float64, float64, error) {
	latitude, err := city.Latitude.Float64Value()
	if err != nil {
		return 0, 0, err
	}

	longitude, err := city.Longitude.Float64Value()
	if err != nil {
		return 0, 0, err
	}

	return latitude.Float64, longitude.Float64, nil
}

func crashWithError(msg string, err error) {
	slog.Warn(msg)
	slog.Error(err.Error())
	os.Exit(1)
}
