package main

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type TemperatureUnit string

const (
	Kelvin  TemperatureUnit = "kelvin"
	Celsius TemperatureUnit = "celsius"
)

type UserPreference struct {
	TempUnit TemperatureUnit `json:"temp_unit"`
	Interval time.Duration   `json:"interval"`
}

func getDefaultUserPreference() *UserPreference {
	return &UserPreference{
		TempUnit: Celsius,
		Interval: 3 * time.Minute,
	}
}

func (c *Config) getUserPreference(w http.ResponseWriter, r *http.Request) {
	var data struct {
		TimeUnit TemperatureUnit `json:"time_unit"`
		Interval string          `json:"interval"`
	}

	data.TimeUnit = c.UserPref.TempUnit
	data.Interval = c.UserPref.Interval.String()

	respData, err := json.Marshal(data)
	if err != nil {
		sendError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(respData)
}

type UpdateUserPreference struct {
	TempUnit string `json:"temp_unit"`
	Interval string `json:"interval"`
}

// updateUserPreference updates the user settings, for making it simple they are only
// on memory (application level) and not on storege (database)
func (c *Config) updateUserPreference(w http.ResponseWriter, r *http.Request) {
	slog.Info("Updating user preference")
	data, err := io.ReadAll(r.Body)
	if err != nil {
		sendError(w, err, http.StatusInternalServerError)
		return
	}

	var newPreference UpdateUserPreference
	err = json.Unmarshal(data, &newPreference)

	if err != nil {
		sendError(w, err, http.StatusInternalServerError)
		return
	}

	newInterval, err := time.ParseDuration(newPreference.Interval)
	if err != nil {
		sendError(w, err, http.StatusInternalServerError)
		return
	}

	if newPreference.TempUnit == "kelvin" {
		c.UserPref.TempUnit = Kelvin
	} else {
		c.UserPref.TempUnit = Celsius
	}
	c.UserPref.Interval = newInterval
	slog.Info("Done with preference update!")

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Preference changed!"))
}
