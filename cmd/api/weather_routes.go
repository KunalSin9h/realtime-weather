package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/kunalsin9h/realtime-weather/internal/db"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

// getDailyWeatherSummary give the today's (current day) weather summary.
// we don't process this on application level, we use TimeScaleDB. It's a
// time-series database on top of Postgres that does the calculation on interval (currently 1 hr)
func (c *Config) getDailyWeatherSummary(w http.ResponseWriter, r *http.Request) {
	cityId := r.PathValue("city_id")

	cityIdInt, err := strconv.ParseInt(cityId, 10, 64)
	if err != nil {
		sendError(w, err)
		return
	}

	var data WeatherSummaryData

	// We are not using sqlc ORM, its giving some problem with get_latest_daily_summary function's return table type.
	// TODO: open issue in the sqlc project
	// so we are doing row sql query with pgx
	err = c.dbConn.QueryRow(r.Context(), `SELECT 
    	date,
		avg_temperature,
		max_temperature,
		min_temperature,
		avg_humidity,
		max_humidity,
		min_humidity,
		avg_wind_speed,
		max_wind_speed,
		min_wind_speed,
		dominant_condition
	FROM get_latest_daily_summary($1)`, cityIdInt).Scan(
		&data.Date,
		&data.AvgTemperature,
		&data.MaxTemperature,
		&data.MinTemperature,
		&data.AvgHumidity,
		&data.MaxHumidity,
		&data.MinHumidity,
		&data.AvgWindSpeed,
		&data.MaxWindSpeed,
		&data.MinWindSpeed,
		&data.DominantCondition,
	)

	// err might be of no rows in result set
	if errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if err != nil {
		sendError(w, err, http.StatusInternalServerError)
		return
	}

	// convert kelvin to Celsius according to User preference
	if c.UserPref.TempUnit == Celsius {
		data.AvgTemperature = convertToCelsius(data.AvgTemperature)
		data.MaxTemperature = convertToCelsius(data.MaxTemperature)
		data.MinTemperature = convertToCelsius(data.MinTemperature)
	}

	respData, err := json.Marshal(data)
	if err != nil {
		sendError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(respData)
}

// refreshDailySummaryViewTable refreshed / recompute the daily_summary_view table.
// This is helpful when we need to process data which is not yet process by timescaleDB interval policy.
func (c *Config) refreshDailySummaryViewTable(w http.ResponseWriter, r *http.Request) {
	slog.Info("Refreshing daily summary data view (table data)")
	query := db.New(c.dbConn)

	if err := query.RefreshDailyWeatherSummary(r.Context()); err != nil {
		sendError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte("Refreshed daily summary data"))
	slog.Info("Done with refreshing!")
}

type WeatherSummaryData struct {
	Date              time.Time `json:"date"`
	AvgTemperature    float64   `json:"avg_temperature"`
	MaxTemperature    float64   `json:"max_temperature"`
	MinTemperature    float64   `json:"min_temperature"`
	AvgHumidity       float64   `json:"avg_humidity"`
	MaxHumidity       float64   `json:"max_humidity"`
	MinHumidity       float64   `json:"min_humidity"`
	AvgWindSpeed      float64   `json:"avg_wind_speed"`
	MaxWindSpeed      float64   `json:"max_wind_speed"`
	MinWindSpeed      float64   `json:"min_wind_speed"`
	DominantCondition string    `json:"dominant_condition"`
}
