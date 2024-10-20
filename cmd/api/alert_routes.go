package main

import (
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

	weadCondId := pgtype.Int4{Valid: false, Int32: 0}
	// weather condition id
	// if there is no condition, that is we don't have alert on condition
	// it ok
	wdId, err := query.GetConditionID(r.Context(), reqPayload.Condition)
	if err == nil {
		// if no error: that means we have condition
		weadCondId.Int32 = wdId
		weadCondId.Valid = true
	}

	err = query.CreateAlertThreshold(r.Context(), db.CreateAlertThresholdParams{
		Name:           reqPayload.Name,
		CityID:         reqPayload.CityID,
		ConditionID:    weadCondId,
		MinTemperature: reqPayload.MinTemperature,
		MaxTemperature: reqPayload.MaxTemperature,
		MinHumidity:    reqPayload.MinHumidity,
		MaxHumidity:    reqPayload.MaxHumidity,
		MinWindSpeed:   reqPayload.MinWindSpeed,
		MaxWindSpeed:   reqPayload.MaxWindSpeed,
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
	// DELETE /alert/3
	// 3 here is alert id
	alertId := r.PathValue("alert_threshold_id")

	alertIdInt, err := strconv.ParseInt(alertId, 10, 32)
	if err != nil {
		sendError(w, err, http.StatusInternalServerError)
		return
	}

	query := db.New(c.dbConn)

	err = query.DeleteAlertThreshold(r.Context(), int32(alertIdInt))
	if err != nil {
		sendError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "Alert Threshold Deleted")
}
