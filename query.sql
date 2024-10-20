-- name: GetConditionID :one
SELECT id FROM weather_conditions WHERE condition = $1;

-- name: AddWeatherData :exec
INSERT INTO weather_data (
    time, condition_id, city_id, temperature, feels_like, humidity, wind_speed
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
);

-- It's a function in migrations folder, is used continuous aggregate feature of timescale
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
    condition_id,
    min_temperature,
    max_temperature,
    min_humidity,
    max_humidity,
    min_wind_speed,
    max_wind_speed
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9
);

-- name: DeleteAlertThreshold :exec
UPDATE alert_thresholds SET active = false WHERE id = $1;

-- Give all the alerts that are in action for a given weather data
-- name: CheckAlertThreshold :many
SELECT * FROM alert_thresholds
WHERE city_id = $1 AND active = true AND (
    condition_id = $2
    OR
    min_temperature >= $3
    OR
    max_temperature <= $4
    OR
    min_humidity >= $5
    OR
    max_humidity <= $6
    OR
    min_wind_speed >= $7
    OR
    max_wind_speed <= $8
);

-- add an entry to alert table for user notification
-- name: CreateAlert :exec
INSERT INTO alerts (
    threshold_id,
    message
) VALUES (
    $1,
    $2
);

-- Get all the alerts in the alerts table, make then active = false (means they are processed)
-- name: DeactivateAndGetAlerts :many
UPDATE alerts al
SET active = false
FROM alert_thresholds th
WHERE al.threshold_id = th.id
  AND al.active = true
  AND th.active = true
RETURNING th.id as threshold_id, th.name as name, al.time as time, al.message as message;
