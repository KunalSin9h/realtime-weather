## Realtime Weather

Realtime Weather application is a Realtime data processing application for daily weather condition. Since we have
 `time-series` data, this project uses [TimescaleDB](https://www.timescale.com/) - its build upon Postgres. And [OpenWeatherAPI](https://openweathermap.org/).

### Architecture

![Architecture](https://i.imgur.com/mysjOxG.png)

### Setup

Requirements

- GoLang
- Cargo (rust's build tool, only for installing `sqlx-cli`)
- `sqlx-cli`
- Docker (option, but recommended)
- GNU Make (optional)

Clone the Repository, and go to rood directory.

#### Download required things

1. `sqlx-cli`: _sqlx_ is rust cli, sql migration tool.

```bash
cargo install sqlx-cli --no-default-features --features postgres

# if its fails, causing some OpenSSL use then download openssl first
# see this: https://docs.rs/crate/openssl-sys/0.9.19
```

#### Run Database

We use TimescaleDB. (docker hub: `timescale/timescaledb:latest-pg16`). Its based on Postgres, so the DB `connection string`,
can have `postgres://....`.

```bash
bash ./rundb.sh
```

OR

```bash
docker run \
    -e POSTGRES_USER=admin \
    -e POSTGRES_PASSWORD=admin \
    -e POSTGRES_DB=weather \
    -p 5432:5432 \
    -d timescale/timescaledb:latest-pg16
```

#### Apply Migrations

> Make sure db is up and running

```bash
# This is connection sting of database we started
# postgres:// works as timescale is based on postgres
# DATABASE_URL env var is required by sqlx for migrations
export DATABASE_URL=postgres://admin:admin@127.0.0.1:5432/weather

sqlx database create
# then
sqlx migrate run
```

`sqlx` uses [migrations]() folder. It holds the entire DB schema.

#### use `OPEN_WEATHER_API_KEY`

Copy [format.env]() into new `.env` file

```bash
cp format.env .env
```

Fill the [.env]() file with Open Weather API key. Application automatically loads `.env`.

#### Run application

```bash
make run

# or 
# go run ./...
```

or with environment variables

```bash
export PORT=5000
export HOST=0.0.0.0
export TIMESCALE=postgres://admin:admin@127.0.0.1:5432/weather
export OPEN_WEATHER_API_KEY=9u93ji...

make run
```

Now Open [http://localhost:5000](http://localhost:5000)

#### Environment Variables

| ENV       | DEFAULT VALUE                                     | USE CASE                    |
|-----------|---------------------------------------------------|-----------------------------|
| PORT      | 5000                                              | application port            |
| HOST      | "0.0.0.0"                                         | application host            |
| TIMESCALE | "postgresql://admin:admin@localhost:5432/weather" | timescale connection string |
| OPEN_WEATHER_API_KEY |  | API Key for Open Weather    |

## Stack

- UI : with React and Vite
- API and Backend : with GoLang
- DB is TimeScaleDB

## UI Demo

![demo](https://i.imgur.com/VPV7GiF.png)

## Results

### 1. Processing and Analysis:

   - [cmd/api/fetcher.go](cmd/api/fetcher.go) is where the concurrent data fetcher is implemented. It uses `interval` (default is _3 minutes_), 
    with user preference, which can be changed by user with API Call.
   
     Used in [cmd/api/main.go](cmd/api/main.go)
     ```go
      // concurrently fetch data at the app.interval interval
      go app.dataSourceFetcher(ctx)
     ```

   - We store all temperature data in the Database in `Kelvin`, to maintain consistency. But user can choose to receive data in API results
     in either `kelvin` or `celsius`, with user preference, which can be changed by user with API Call.

### 2. Rollups and Aggregates:

  - **Daily Weather Summary**
   
    Daily weather summary for {`temperature`, `humidity` and `wind speed`}'s `Avg`, `Max` and `Min` are calculated. 

    It is **not** calculated at `Application Level`, but on `Database Level`.  We use `TimescaleDB`'s _continuous aggregate_.
   And in the interval of `1 Hour`, daily weather summary of each city is calculated. Its a `Materialized View`.
    - [migrations/20241018155039_create_daily_weather_summary_view.sql](migrations/20241018155039_create_daily_weather_summary_view.sql)
    - [migrations/20241018163503_get_latest_daily_summary.sql](migrations/20241018163503_get_latest_daily_summary.sql)
    - used in go code: [cmd/api/weather_routes.go](cmd/api/weather_routes.go)
 
  - **Alerting Thresholds**:
    - User can define `Alert Thresholds`. We use `alert_thresholds` table to store them.
    - In the Data fetch ([cmd/api/fetcher.go](cmd/api/fetcher.go)). We concurrently check for any `thresholds violation`
         for each weather data entry.
        
        For every new data point in the table.
      ```go
        // in background
        // concurrently check for any alert threshold on this dbData (weather data)
        go checkThresholds(ctx, dbData, wdID, query)
      ```    
    - Upon breach of thresholds. we store the alerts in the `alert` table, with message. And user can `Pool` these alerts, for showing in dashboard. see [cmd/api/send_alert.go](cmd/api/send_alert.go)

### 3. Test Cases

Run test cases

> [!CAUTION]
> Make sure to have DB up and running, as same as earlier

```bash
make test

# or
# go test ./...
```

## Design Choices

1. Use of `TimescaleDB`.

    Our weather data was a `time-series` data. And for that a _time series_ database become a feasible choice. I have used
       _timescale_ as it is just an extension over Postgres. And it as easy abstraction called [Continuous Aggregate](https://docs.timescale.com/use-timescale/latest/continuous-aggregates/create-a-continuous-aggregate/),
        which can easily calculate the required `rollups and aggrigate` with scheduled interval.  See [migrations](migrations) folder.

    We have a sql query to refresh this aggregate view. It will recalculate with the fresh data that is under 1 hr.  see [query.sql](query.sql)
 
    ```sql
    -- name: RefreshDailyWeatherSummary :exec
    CALL refresh_continuous_aggregate('daily_weather_summary_view', NULL, NULL);
    -- Manually Refresh the daily_weather_summary_view table
    ```

2. Use of `sqlc` (not `sqlx`).
 
    I have used `sqlc`, it a golang type-safe ORM which generate from raw sql queries. See [query.sql](query.sql)

3. Use of React + Tailwind + Shadcn

    I have used React with Vite and Shadcd for their simplicity and quick building process.