package main

import (
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kunalsin9h/realtime-weather/internal/db"
	"log/slog"
	"net/http"
	"slices"
	"strconv"
	"sync/atomic"
)

type LiveData struct {
	CityID      int32              `json:"city_id"`
	Time        pgtype.Timestamptz `json:"time"`
	Temperature float64            `json:"temperature"`
	Humidity    float64            `json:"humidity"`
	WindSpeed   float64            `json:"wind_speed"`
}

// LiveDataStreams is where the fetcher (fetcher.go) sends live data
// And we broadcast to client
var LiveDataStreams []*chan LiveData

// liveWeatherData send latest weather data for a city with SSE
func (c *Config) liveWeatherData(w http.ResponseWriter, r *http.Request) {
	slog.Info("Sending live data to cities")

	channelForThisVeryClient := make(chan LiveData)
	LiveDataStreams = append(LiveDataStreams, &channelForThisVeryClient)

	defer func() {
		// remove the client channel one connection is lost form client
		idx := slices.Index(LiveDataStreams, &channelForThisVeryClient)
		LiveDataStreams[idx] = nil
		slices.Delete(LiveDataStreams, idx, idx+1)
	}()

	// SSE Setup
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher := w.(http.Flusher)

	cityID := r.PathValue("city_id")
	cityIDInt, err := strconv.ParseInt(cityID, 10, 32)
	if err != nil {
		sendError(w, err, http.StatusInternalServerError)
		return
	}

	id := atomic.Int64{}

	// send today's data first
	query := db.New(c.dbConn)
	todayWeather, err := query.GetTodaysWeatherData(r.Context())
	if err != nil {
		sendError(w, err, http.StatusInternalServerError)
		return
	}

	for _, today := range todayWeather {
		liveData := LiveData{
			Time:        today.Time,
			CityID:      today.CityID,
			Temperature: pgToFloat(today.Temperature),
			Humidity:    pgToFloat(today.Humidity),
			WindSpeed:   pgToFloat(today.WindSpeed),
		}

		if c.UserPref.TempUnit == Celsius {
			liveData.Temperature = convertToCelsius(liveData.Temperature)
		}

		resp, err := json.Marshal(liveData)
		if err != nil {
			sendError(w, err, http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, getSSEData("message", string(resp), id.Load(), 1000))
		flusher.Flush()
		id.Add(1)
	}

	// then send new data in realtime
	for {
		select {
		case data := <-channelForThisVeryClient:
			if data.CityID != int32(cityIDInt) {
				// this data need not be sent on this client
				continue
			}

			if c.UserPref.TempUnit == Celsius {
				data.Temperature = convertToCelsius(data.Temperature)
			}

			resp, err := json.Marshal(data)
			if err != nil {
				sendError(w, err, http.StatusInternalServerError)
				return
			}
			fmt.Fprintf(w, getSSEData("message", string(resp), id.Load(), 1000))
			flusher.Flush()

			id.Add(1)
		case <-r.Context().Done():
			return
		}
	}
}
