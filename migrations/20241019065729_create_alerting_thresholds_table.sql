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
    active BOOLEAN NOT NULL DEFAULT true, -- for confirmed alert delivery, we can make this active to false
    FOREIGN KEY (threshold_id) REFERENCES alert_thresholds(id)
);
