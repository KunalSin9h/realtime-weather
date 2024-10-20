package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kunalsin9h/realtime-weather/internal/db"
	"log/slog"
)

// checkThresholds check for any alert threshold which is in action
// for a weather data.
// Usage
// it is used in fetcher.go
// after every new data entry
func checkThresholds(ctx context.Context, weatherData db.AddWeatherDataParams, conditionID int32, query *db.Queries) {
	slog.Info(fmt.Sprintf("checking for alert on city"))
	// alert threshold that might be applicable for this weatherData
	alertThs, err := query.CheckAlertThreshold(ctx, db.CheckAlertThresholdParams{
		CityID:         weatherData.CityID,
		ConditionID:    pgtype.Int4{Int32: conditionID, Valid: true},
		MinTemperature: weatherData.Temperature,
		MaxTemperature: weatherData.Temperature,
		MinHumidity:    weatherData.Humidity,
		MaxHumidity:    weatherData.Humidity,
		MinWindSpeed:   weatherData.WindSpeed,
		MaxWindSpeed:   weatherData.WindSpeed,
	})

	slog.Info(fmt.Sprintf("Found %d alert applicable for this city", len(alertThs)))

	if err != nil {
		slog.Warn("Failed to find all alert thresholds")
		slog.Error(err.Error())
		return
	}

	// all these alerts are application for this weatherData
	for _, alertTh := range alertThs {
		// make an entry in the alerts table for sending these alerts
		err := query.CreateAlert(ctx, db.CreateAlertParams{
			ThresholdID: alertTh.ID,
			Message:     generateAlertMessage(alertTh),
		})

		if err != nil {
			slog.Warn("Failed to create alert threshold")
			slog.Error(err.Error())
			return
		}

		slog.Info("Created alert!")
	}
}

// generateAlertMessage create message for indicating the cause of alert
func generateAlertMessage(alert db.AlertThreshold) string {
	// We are doing a hack, to reduce db query and any other operation
	// by just seeing "On what point we have set the alert"
	// and if there is alert, then we can say "that the only point that passes threshold"

	msg := "Alert exceeds "

	switch {
	case alert.MaxTemperature.Valid:
		msg += "max temperature"
	case alert.MinTemperature.Valid:
		msg += "min temperature"
	case alert.MinHumidity.Valid:
		msg += "min humidity"
	case alert.MaxHumidity.Valid:
		msg += "max humidity"
	case alert.MinWindSpeed.Valid:
		msg += "min wind speed"
	case alert.MaxWindSpeed.Valid:
		msg += "max wind speed"
	case alert.ConditionID.Valid:
		msg += "condition for weather"
	default:
		msg = "Unknown alert occur"
	}

	return msg
}
