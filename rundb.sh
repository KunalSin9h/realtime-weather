#!usr/bin/env bash

set -eox pipefail

# Check if the sqlx-cli is install on the system
if ! [ -x "$(command -v sqlx)" ]; then
    echo >&2 'Error: sqlx is not installed.'
    echo >&2 'Use:'
    echo >&2 '    cargo install sqlx-cli --no-default-features --features postgres'
    echo >&2 'to install it.'
    exit 1
fi

# Check if the custom user has been set, otherwise default to 'admin'
DB_USER="${POSTGRES_USER:=admin}"

# Check if the custom password has been set, otherwise default to 'admin'
DB_PASSWORD="${POSTGRES_PASSWORD:=admin}"

# Check if a custom database name has been set, otherwise default to 'weather'
DB_NAME="${POSTGRES_DB:=weather}"

# Check if a custom post has been set, otherwise default to '5432'
DB_PORT="${POSTGRES_PORT:=5432}"

# Launch postgres using docker if not running
if [[ -z "${SKIP_DOCKER}" ]]; then
docker run \
    -e POSTGRES_USER=${DB_USER} \
    -e POSTGRES_PASSWORD=${DB_PASSWORD} \
    -e POSTGRES_DB=${DB_NAME} \
    -p "${DB_PORT}":5432 \
    -d timescale/timescaledb:latest-pg16  \
    postgres -N 1000
    # ^^^^^^^^^^^^^ Increased maximum number of connections for testing purpose
fi

>&2 echo "TimescaleDB is up and running on port ${DB_PORT}"

echo ""
echo postgres://${DB_USER}:${DB_PASSWORD}@127.0.0.1:${DB_PORT}/${DB_NAME}
echo ""
