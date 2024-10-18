-- A function to get latest daily summary for a city
CREATE OR REPLACE FUNCTION get_latest_daily_summary(city_id_param INTEGER)
RETURNS TABLE (
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
    dominant_condition_id INTEGER
) AS $$
BEGIN
RETURN QUERY
SELECT
    bucket::DATE AS date,
        avg_temperature,
        max_temperature,
        min_temperature,
        avg_humidity,
        max_humidity,
        min_humidity,
        avg_wind_speed,
        max_wind_speed,
        min_wind_speed,
        dominant_condition_id
FROM daily_weather_summary_view
WHERE city_id = city_id_param
ORDER BY bucket DESC
    LIMIT 1;
END;
$$ LANGUAGE plpgsql;
