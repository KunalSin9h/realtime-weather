-- name: AddWeatherData :exec
INSERT INTO weather_data (
    time, city_id, temperature, feels_like, humidity, wind_speed, condition_id
) VALUES (
  $1, $2, $3, $4, $5, $6, (
        SELECT id FROM weather_conditions WHERE condition = $7
    )
);

-- name: GetTodayWeatherSummary :one
SELECT * FROM get_latest_daily_summary($1);

-- name: RefreshDailyWeatherSummary :exec
CALL refresh_continuous_aggregate('daily_weather_summary_view', localtimestamp - INTERVAL '1 hour', localtimestamp);
-- Manually Refresh the daily_weather_summary_view of past 1 hour

-- CITIES
-- name: GetAllCities :many
SELECT * FROM cities;

-- ALERTS
-- name: CreateAlertThreshold :exec
INSERT INTO alert_thresholds (
    name,
    city_id,
    min_temperature,
    max_temperature,
    min_humidity,
    max_humidity,
    min_wind_speed,
    max_wind_speed,
    occur_limit,
    condition_id
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9,
    (
        SELECT id FROM weather_conditions WHERE condition = $10
    )
);
