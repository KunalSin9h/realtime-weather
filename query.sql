-- name: AddWeatherData :exec
INSERT INTO weather_data (
    time, city_id, condition_id, temperature, feels_like, humidity, wind_speed
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
);

-- name: GetTodayWeatherSummary :one
SELECT * FROM get_latest_daily_summary($1);

-- name: RefreshDailyWeatherSummary :exec
CALL refresh_continuous_aggregate('daily_weather_summary_view', localtimestamp - INTERVAL '1 hour', localtimestamp);
-- Manually Refresh the daily_weather_summary_view of past 1 hour
