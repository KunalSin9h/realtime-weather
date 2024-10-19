package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kunalsin9h/realtime-weather/internal/db"
	"io"
	"net/http"
	"strconv"
)

type CreateAlertPayload struct {
	Name           string         `json:"name"`
	CityID         int32          `json:"city_id"`
	Condition      string         `json:"condition,omitempty"` // Rain, snow etc
	MinTemperature pgtype.Numeric `json:"min_temperature,omitempty"`
	MaxTemperature pgtype.Numeric `json:"max_temperature,omitempty"`
	MinHumidity    pgtype.Numeric `json:"min_humidity,omitempty"`
	MaxHumidity    pgtype.Numeric `json:"max_humidity,omitempty"`
	MinWindSpeed   pgtype.Numeric `json:"min_wind_speed,omitempty"`
	MaxWindSpeed   pgtype.Numeric `json:"max_wind_speed,omitempty"`
	OccurLimit     int32          `json:"occur_limit,omitempty"`
}

// CreateAlert creates a new Alert Threshold entry in the database
func (c *Config) createAlert(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// ready request payload
	data, err := io.ReadAll(r.Body)
	if err != nil {
		sendError(w, err)
		return
	}

	// unmarshal data
	var reqPayload CreateAlertPayload

	err = json.Unmarshal(data, &reqPayload)
	if err != nil {
		sendError(w, err, http.StatusInternalServerError)
		return
	}

	query := db.New(c.dbConn)

	// occurLimit it limit, on how many times this condition need to occur before alerting
	occurLimit := int32(1) // default
	if reqPayload.OccurLimit > 0 {
		occurLimit = reqPayload.OccurLimit
	}

	// weather condition id
	wdId, err := query.GetConditionID(ctx, reqPayload.Condition)
	if err != nil {
		sendError(w, err, http.StatusInternalServerError)
		return
	}

	err = query.CreateAlertThreshold(ctx, db.CreateAlertThresholdParams{
		Name:           reqPayload.Name,
		CityID:         reqPayload.CityID,
		ConditionID:    pgtype.Int4{Int32: wdId, Valid: true},
		MinTemperature: reqPayload.MinTemperature,
		MaxTemperature: reqPayload.MaxTemperature,
		MinHumidity:    reqPayload.MinHumidity,
		MaxHumidity:    reqPayload.MaxHumidity,
		MinWindSpeed:   reqPayload.MinWindSpeed,
		MaxWindSpeed:   reqPayload.MaxWindSpeed,
		OccurLimit:     occurLimit,
	})

	if err != nil {
		sendError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "Alert Threshold Created")
}

// deleteAlert is a API endpoint handler for DELETE /alert/{id}
func (c *Config) deleteAlert(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// DELETE /alert/3
	// 3 here is alert id
	alertId := r.PathValue("id")

	alertIdInt, err := strconv.ParseInt(alertId, 10, 32)
	if err != nil {
		sendError(w, err, http.StatusInternalServerError)
		return
	}

	query := db.New(c.dbConn)

	err = query.DeleteAlertThreshold(ctx, int32(alertIdInt))
	if err != nil {
		sendError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "Alert Threshold Deleted")
}
