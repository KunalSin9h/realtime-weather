// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Alert struct {
	ID          int32
	Time        pgtype.Timestamptz
	ThresholdID int32
	Message     string
	Active      bool
}

type AlertThreshold struct {
	ID             int32
	Name           string
	CityID         int32
	ConditionID    pgtype.Int4
	MinTemperature pgtype.Numeric
	MaxTemperature pgtype.Numeric
	MinHumidity    pgtype.Numeric
	MaxHumidity    pgtype.Numeric
	MinWindSpeed   pgtype.Numeric
	MaxWindSpeed   pgtype.Numeric
	Active         bool
}

type City struct {
	ID        int32
	Name      string
	Latitude  pgtype.Numeric
	Longitude pgtype.Numeric
}

type DailyWeatherSummaryView struct {
	Bucket              interface{}
	CityID              int32
	AvgTemperature      pgtype.Numeric
	MaxTemperature      pgtype.Numeric
	MinTemperature      pgtype.Numeric
	AvgHumidity         pgtype.Numeric
	MaxHumidity         pgtype.Numeric
	MinHumidity         pgtype.Numeric
	AvgWindSpeed        pgtype.Numeric
	MaxWindSpeed        pgtype.Numeric
	MinWindSpeed        pgtype.Numeric
	DominantConditionID interface{}
}

type WeatherCondition struct {
	ID        int32
	Condition string
}

type WeatherDatum struct {
	Time        pgtype.Timestamptz
	CityID      int32
	ConditionID int32
	Temperature pgtype.Numeric
	FeelsLike   pgtype.Numeric
	Humidity    pgtype.Numeric
	WindSpeed   pgtype.Numeric
}
