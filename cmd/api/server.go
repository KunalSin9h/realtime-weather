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

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", HOST, PORT),
		Handler: mux,
	}

	slog.Info(fmt.Sprintf("Listening on http://%s", server.Addr))
	return server.ListenAndServe()
}
