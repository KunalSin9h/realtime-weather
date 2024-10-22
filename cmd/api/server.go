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

	// Serve Static Frontend Dashboard
	// ALL UI ROUTES
	// Since, the frontend also served by GO backend, we really need to do this mess. this is how it works :)
	mux.Handle("GET /", http.FileServer(http.Dir("./ui/dist/")))
	mux.HandleFunc("GET /settings", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./ui/dist/index.html")
	})
	mux.HandleFunc("GET /city/{city_name}/{city_id}", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./ui/dist/index.html")
	})

	///***********************************************************************
	// ALL THINGS API

	// GET /cities
	// Get all the cities
	mux.HandleFunc("GET /api/cities", enableCorsAnd(c.getAllCities))

	// Get Daily Weather Summary for City with city_id
	mux.HandleFunc("GET /api/cities/summary/{city_id}", enableCorsAnd(c.getDailyWeatherSummary))

	// Refresh Daily Weather Summary
	mux.HandleFunc("POST /api/cities/summary/refresh", enableCorsAnd(c.refreshDailySummaryViewTable))

	// GET Weather Data for A City
	// Realtime using SSE (Server Sent Events)
	mux.HandleFunc("GET /api/cities/live/{city_id}", enableCorsAnd(c.liveWeatherData))

	// Send alerts on API Pooling
	// alerts are created
	mux.HandleFunc("GET /api/alerts", enableCorsAnd(c.sendAlert))

	// Create Alert Thresholds
	mux.HandleFunc("POST /api/alert", enableCorsAnd(c.createAlert))

	// Delete Alert Thresholds with ALERT ID
	mux.HandleFunc("DELETE /api/alert/{alert_threshold_id}", enableCorsAnd(c.deleteAlert))

	// User Preference
	// Get user preference
	mux.HandleFunc("GET /api/preference", enableCorsAnd(c.getUserPreference))
	// change user preference
	mux.HandleFunc("POST /api/preference", enableCorsAnd(c.updateUserPreference))

	// SERVER
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", HOST, PORT),
		Handler: mux,
	}

	slog.Info(fmt.Sprintf("Listening on http://%s", server.Addr))
	return server.ListenAndServe()
}

func enableCorsAnd(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Vite app
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Accept, Origin, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// getAllCities give list of cities in the database
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
