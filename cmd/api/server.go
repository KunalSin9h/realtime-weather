package main

import (
	"encoding/json"
	"fmt"
	"github.com/kunalsin9h/realtime-weather/internal/db"
	"log/slog"
	"net/http"
)

func (c *Config) setUpAndRunServer() error {
	slog.Info("Starting server...")
	mux := http.NewServeMux()

	// GET /
	// Serve Static Frontend Dashboard
	mux.Handle("GET /", http.FileServer(http.Dir("./ui/dist/")))

	// GET /cities
	// Get all the cities
	mux.HandleFunc("GET /cities", c.getAllCities)

	// Get Daily Weather Summary for City with city_id
	mux.HandleFunc("GET /summary/{city_id}", c.getDailyWeatherSummary)

	// Refresh Daily Weather Summary
	mux.HandleFunc("POST /summary/refresh", c.refreshDailySummaryViewTable)

	// change user preference
	// temperature

	// ALERTS

	// Send alerts on API Pooling
	// alerts are created
	mux.HandleFunc("GET /alerts", c.sendAlert)

	// Create Alert Thresholds
	mux.HandleFunc("POST /alert", c.createAlert)

	// Delete Alert Thresholds with ALERT ID
	mux.HandleFunc("DELETE /alert/{alert_threshold_id}", c.deleteAlert)

	// User Preference
	// Since there are only two user preferences, I am doing simple way.
	// else it would be better to have single user preference API and database entry
	mux.HandleFunc("GET /preference", c.getUserPreference)
	// change Interval
	mux.HandleFunc("POST /preference/interval/{new_interval}", c.changeInterval)
	// change Temperature Unit
	mux.HandleFunc("POST /preference/temp_unit/{new_unit}", c.changeTempUnit)

	// SERVER
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", HOST, PORT),
		Handler: mux,
	}

	slog.Info(fmt.Sprintf("Listening on http://%s", server.Addr))
	return server.ListenAndServe()
}

// getAllCities give list of citie in the database
func (c *Config) getAllCities(w http.ResponseWriter, r *http.Request) {
	query := db.New(c.dbConn)
	cities, err := query.GetAllCities(r.Context())

	if err != nil {
		sendError(w, err, http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(cities)
	if err != nil {
		sendError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
