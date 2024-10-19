package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kunalsin9h/realtime-weather/internal/db"
	"io"
	"net/http"
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
func (c *Config) CreateAlert(w http.ResponseWriter, r *http.Request) {
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

	err = query.CreateAlertThreshold(ctx, db.CreateAlertThresholdParams{
		Name:           reqPayload.Name,
		CityID:         reqPayload.CityID,
		Condition:      reqPayload.Condition,
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

func (c *Config) DeleteAlert(w http.ResponseWriter, r *http.Request) {

}
