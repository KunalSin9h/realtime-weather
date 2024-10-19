-- name: AddWeatherData :exec
INSERT INTO weather_data (
    time, city_id, condition_id, temperature, feels_like, humidity, wind_speed
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
);

-- name: GetTodayWeatherSummary :one
SELECT * FROM get_latest_daily_summary($1);
