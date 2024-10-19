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

	// Create Alert Thresholds
	mux.HandleFunc("POST /alert", c.CreateAlert)

	// Delete Alert Thresholds
	mux.HandleFunc("DELETE /alert", c.DeleteAlert)

	// SERVER
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", HOST, PORT),
		Handler: mux,
	}

	slog.Info(fmt.Sprintf("Listening on http://%s", server.Addr))
	return server.ListenAndServe()
}
