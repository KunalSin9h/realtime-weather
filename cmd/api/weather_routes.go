package main

import (
	"encoding/json"
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
	err = c.dbConn.QueryRow(r.Context(), "SELECT * FROM get_latest_daily_summary($1)", cityIdInt).Scan(
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

	if err != nil {
		sendError(w, err, http.StatusInternalServerError)
		return
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
	AvgTemperature    float64
	MaxTemperature    float64
	MinTemperature    float64
	AvgHumidity       float64
	MaxHumidity       float64
	MinHumidity       float64
	AvgWindSpeed      float64
	MaxWindSpeed      float64
	MinWindSpeed      float64
	DominantCondition string
}
