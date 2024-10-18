-- docs: https://docs.timescale.com/use-timescale/latest/continuous-aggregates/create-a-continuous-aggregate/

-- Create a materialized view for the daily weather summary
CREATE MATERIALIZED VIEW daily_weather_summary_view
WITH (timescaledb.continuous) AS -- for timescaleDB to manage this as a continuous aggregate.
SELECT
    time_bucket('1 day', time) AS bucket,
    city_id,
    AVG(temperature)::DECIMAL(5,2) AS avg_temperature,
    MAX(temperature)::DECIMAL(5,2) AS max_temperature,
    MIN(temperature)::DECIMAL(5,2) AS min_temperature,
    AVG(humidity)::DECIMAL(5,2) AS avg_humidity,
    MAX(humidity)::DECIMAL(5,2) AS max_humidity,
    MIN(humidity)::DECIMAL(5,2) AS min_humidity,
    AVG(wind_speed)::DECIMAL(5,2) AS avg_wind_speed,
    MAX(wind_speed)::DECIMAL(5,2) AS max_wind_speed,
    MIN(wind_speed)::DECIMAL(5,2) AS min_wind_speed,
    mode() WITHIN GROUP (ORDER BY condition_id) AS dominant_condition_id
FROM weather_data
GROUP BY bucket, city_id WITH NO DATA;

-- Create a policy to refresh the view
-- This will update in every 1 hour
SELECT add_continuous_aggregate_policy('daily_weather_summary_view',
   start_offset => INTERVAL '3 days', -- This means we'll recompute data as far back as 3 days ago.
   end_offset => INTERVAL '1 hour', --  We'll compute up to 1 hour before the current time.
   schedule_interval => INTERVAL '1 hour'); -- The policy will run evey hour
