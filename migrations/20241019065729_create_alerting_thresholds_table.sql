-- Alert Thresholds
CREATE TABLE alert_thresholds (
    id SERIAL PRIMARY KEY,
    name VARCHAR(125) NOT NULL,
    city_id INTEGER NOT NULL,
    condition_id INTEGER, -- For condition thresholds, like alert me when it "Rains"
    min_temperature DECIMAL(5,2),
    max_temperature DECIMAL(5,2),
    min_humidity DECIMAL(5,2),
    max_humidity DECIMAL(5,2),
    min_wind_speed DECIMAL(5,2),
    max_wind_speed DECIMAL(5,2),
    occur_limit INTEGER NOT NULL DEFAULT 1, -- For watching number of alerts, if only >= update_limit then send alerts
    -- 35 degrees Celsius for $update_limit consecutive updates then only send alert., here update_limit can be 2
    occur_count INTEGER NOT NULL DEFAULT 0, -- Number of times this threshold is passed, ony alert then occur_count = occur_limit
    active BOOLEAN NOT NULL DEFAULT true,
    FOREIGN KEY (city_id) REFERENCES cities(id),
    FOREIGN KEY (condition_id) REFERENCES weather_conditions(id)
);

CREATE INDEX ON alert_thresholds (city_id);

-- Alerts
CREATE TABLE alerts (
    id SERIAL PRIMARY KEY,
    time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    threshold_id INTEGER NOT NULL,
    message TEXT NOT NULL,
    active BOOLEAN NOT NULL DEFAULT true,
    FOREIGN KEY (threshold_id) REFERENCES alert_thresholds(id)
);
