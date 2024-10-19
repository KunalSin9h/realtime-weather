package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kunalsin9h/realtime-weather/internal/db"
	"io"
	"log/slog"
	"net/http"
	"time"
)

// fetcher works concurrently with main go route.
// it will continuously fetch the API and update the data on the given interval
func (c *Config) dataSourceFetcher(ctx context.Context, apiKey string) {
	slog.Info("Started fetching weather data...")

	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	query := db.New(c.dbConn)

	// Run this upfront
	err := c.findAllCitiesAndGetWeather(ctx, query, apiKey)
	if err != nil {
		slog.Error(err.Error())
	}

	// Infinite loop: continuously fetch latest weather data for all cities
	for {
		select {
		case <-ticker.C:
			// Now run it on interval
			err := c.findAllCitiesAndGetWeather(ctx, query, apiKey)
			if err != nil {
				slog.Error(err.Error())
			}
		case <-ctx.Done():
			return
		}

		// Keep updated the interval, user can change this
		ticker.Reset(c.interval)
	}
}

func (c *Config) findAllCitiesAndGetWeather(ctx context.Context, query *db.Queries, apiKey string) error {
	// fetch all weather data for all city again
	// get all the cities
	cities, err := query.GetAllCities(ctx)

	if err != nil {
		return err
	}

	// For all cities
	for _, city := range cities {
		err = fetchWeatherData(ctx, &city, query, apiKey)
		if err != nil {
			slog.Error(err.Error())
		}
	}

	return nil
}

type WeatherData struct {
	Weather  []WeatherCondition `json:"weather"`
	Main     WeatherParameters  `json:"main"`
	Wind     WindCondition      `json:"wind"`
	DT       int64              `json:"dt"` // UNIX timestamp for weather data update
	Timezone int64              `json:"timezone"`
}

type WeatherCondition struct {
	Main string `json:"main"`
}

type WeatherParameters struct {
	Temp      pgtype.Numeric `json:"temp"`
	FeelsLike pgtype.Numeric `json:"feels_like"`
	Humidity  pgtype.Numeric `json:"humidity"`
}

type WindCondition struct {
	Speed pgtype.Numeric `json:"speed"`
}

func fetchWeatherData(ctx context.Context, city *db.City, query *db.Queries, apiKey string) error {
	slog.Info(fmt.Sprintf("Fetching weather data for city %s", city.Name))

	lat, lon, err := getFloatLatLon(city)
	if err != nil {
		return err
	}

	apiUrl := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%v&lon=%v&appid=%s&units=metric",
		lat, lon, apiKey)

	resp, err := http.Get(apiUrl)

	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("open weather api request got: status %s", resp.Status)
	}

	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	var weatherData WeatherData

	err = json.Unmarshal(respBody, &weatherData)

	if err != nil {
		return err
	}

	// find what weather condition id (in db) we have
	// weather condition id
	wcId, err := query.GetWeatherConditionID(ctx, weatherData.Weather[0].Main) // like "Rain", "Haze"

	if err != nil {
		return err
	}

	// insert the weather data into db
	err = query.AddWeatherData(ctx, db.AddWeatherDataParams{
		Time:        calTimestampWithTZ(weatherData.DT, weatherData.Timezone),
		CityID:      city.ID,
		ConditionID: wcId,
		Temperature: weatherData.Main.Temp,
		FeelsLike:   weatherData.Main.FeelsLike,
		Humidity:    weatherData.Main.Humidity,
		WindSpeed:   weatherData.Wind.Speed,
	})

	slog.Info("Done")
	return err
}
