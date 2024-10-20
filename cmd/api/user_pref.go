package main

import (
	"encoding/json"
	"fmt"
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

func (c *Config) changeInterval(w http.ResponseWriter, r *http.Request) {
	newInterval := r.PathValue("new_interval")

	interval, err := time.ParseDuration(newInterval)
	if err != nil {
		sendError(w, err, http.StatusBadRequest)
		return
	}

	c.UserPref.Interval = interval

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Interval changed!"))
}

func (c *Config) changeTempUnit(w http.ResponseWriter, r *http.Request) {
	newUnit := r.PathValue("new_unit")

	switch newUnit {
	case "kelvin":
		c.UserPref.TempUnit = Kelvin
	case "celsius":
		c.UserPref.TempUnit = Celsius
	default:
		sendError(w, fmt.Errorf("unknown temperature unit: %s", newUnit), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Temperature Unit changed!"))
}
