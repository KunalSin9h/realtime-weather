-- Add default cities
INSERT INTO cities (name, latitude, longitude)
VALUES
    ('Delhi', 28.6139, 77.2090),
    ('Mumbai', 19.0760, 72.8777),
    ('Chennai', 13.0827, 80.2707),
    ('Bangalore', 12.9716, 77.5946),
    ('Kolkata', 22.5726, 88.3639),
    ('Hyderabad', 17.3850, 78.4867);

-- Add default weather conditions
INSERT INTO weather_conditions (condition)
VALUES
    ('Thunderstorm'),
    ('Drizzle'),
    ('Rain'),
    ('Snow'),
    ('Mist'),
    ('Smoke'),
    ('Haze'),
    ('Dust'),
    ('Fog'),
    ('Sand'),
    ('Ash'),
    ('Squall'),
    ('Tornado'),
    ('Clear'),
    ('Clouds');
