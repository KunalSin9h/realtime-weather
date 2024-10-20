-- Daily summary type, used by get_latest_daily_summary
CREATE TYPE daily_summary AS (
    date DATE,
    avg_temperature DECIMAL(5,2),
    max_temperature DECIMAL(5,2),
    min_temperature DECIMAL(5,2),
    avg_humidity DECIMAL(5,2),
    max_humidity DECIMAL(5,2),
    min_humidity DECIMAL(5,2),
    avg_wind_speed DECIMAL(5,2),
    max_wind_speed DECIMAL(5,2),
    min_wind_speed DECIMAL(5,2),
    dominant_condition VARCHAR(50)
);

-- A function to get latest daily summary for a city
CREATE OR REPLACE FUNCTION get_latest_daily_summary(city_id_param INTEGER)
RETURNS SETOF daily_summary AS $$
BEGIN
RETURN QUERY
SELECT
    bucket::DATE AS date,
    daily_weather_summary_view.avg_temperature,
    daily_weather_summary_view.max_temperature,
    daily_weather_summary_view.min_temperature,
    daily_weather_summary_view.avg_humidity,
    daily_weather_summary_view.max_humidity,
    daily_weather_summary_view.min_humidity,
    daily_weather_summary_view.avg_wind_speed,
    daily_weather_summary_view.max_wind_speed,
    daily_weather_summary_view.min_wind_speed,
    (
        SELECT condition FROM weather_conditions
        WHERE id = daily_weather_summary_view.dominant_condition_id
    )
FROM daily_weather_summary_view
WHERE city_id = city_id_param
ORDER BY bucket DESC
    LIMIT 1;
END;
$$ LANGUAGE plpgsql;
