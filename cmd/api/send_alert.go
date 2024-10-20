package main

import (
	"encoding/json"
	"github.com/kunalsin9h/realtime-weather/internal/db"
	"net/http"
)

// Possible implementation:
// Send Alert check alert in the alert queue (alert table working as queue)
// and send alert current on the dashboard but can be extended for email etc.

// Current workflow:
// User will check for updates using API Pooling.
// GET /alerts
// We will give alerts and make them active = false (means alerts are given)
func (c *Config) sendAlert(w http.ResponseWriter, r *http.Request) {
	// for all cities give alerts to the user in json format

	query := db.New(c.dbConn)

	alerts, err := query.DeactivateAndGetAlerts(r.Context())
	if err != nil {
		sendError(w, err, http.StatusInternalServerError)
		return
	}

	respData, err := json.Marshal(alerts)
	if err != nil {
		sendError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(respData)

	if err != nil {
		sendError(w, err, http.StatusInternalServerError)
		return
	}
}
