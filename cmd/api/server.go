package main

import (
	"fmt"
	"log/slog"
	"net/http"
)

func (c *Config) setUpAndRunServer() error {
	slog.Info("Starting server...")
	mux := http.NewServeMux()

	// GET /
	// Serve Static Frontend Dashboard
	mux.Handle("/", http.FileServer(http.Dir("./ui/dist/")))

	// Get Daily Weather Summary
	// Refresh Daily Weather Summary

	// ALERTS

	// Send alerts on API Pooling
	// alerts are created
	mux.HandleFunc("GET /alerts", c.sendAlert)

	// Create Alert Thresholds
	mux.HandleFunc("POST /alert", c.createAlert)

	// Delete Alert Thresholds with ALERT ID
	mux.HandleFunc("DELETE /alert/{id}", c.deleteAlert)

	// SERVER
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", HOST, PORT),
		Handler: mux,
	}

	slog.Info(fmt.Sprintf("Listening on http://%s", server.Addr))
	return server.ListenAndServe()
}
