package main

import "time"

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
