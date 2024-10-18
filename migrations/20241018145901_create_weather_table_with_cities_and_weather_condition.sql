-- Create the cities table
CREATE TABLE cities (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);

-- Create the weather_conditions table
-- Like: rain, clear, clouds, snow etc...
CREATE TABLE weather_conditions (
    id SERIAL PRIMARY KEY,
    condition VARCHAR(50) NOT NULL
);

-- Create the weather_data hypertable (TimescaleDB Abstraction)
CREATE TABLE weather_data (
    time TIMESTAMPTZ NOT NULL,
    city_id INTEGER NOT NULL,
    condition_id INTEGER NOT NULL,
    temperature DECIMAL(5,2) NOT NULL,
    feels_like DECIMAL(5,2) NOT NULL,
    humidity DECIMAL(5, 2) NOT NULL,
    wind_speed DECIMAL(5, 2) NOT NULL,
    FOREIGN KEY (city_id) REFERENCES cities(id),
    FOREIGN KEY (condition_id) REFERENCES weather_conditions(id)
);

-- Convert weather_data table to a timescaleDB hypertable
SELECT create_hypertable('weather_data', 'time');
