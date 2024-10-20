package main

type TemperatureUnit string

const (
	Kelvin  TemperatureUnit = "kelvin"
	Celsius TemperatureUnit = "celsius"
)

type UserPreference struct {
	TempUnit TemperatureUnit `json:"temp_unit"`
}
