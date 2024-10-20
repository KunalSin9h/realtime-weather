## Realtime Weather

Install Development Dependencies

```bash
# for ubuntu
sudo apt-get install pkg-config libssl-dev
# for macos
# brew install openssl
# more details https://docs.rs/crate/openssl-sys/0.9.19
cargo install sqlx-cli --no-default-features --features postgres

# for psql
sudo apt-get install -y postgresql-client
```

Design Choices

1. Use of `CURRENT_TIMESAMP` as `DEFAULT` value in `time` colum.

To simply the applicaion layer login

```bash
#run migration

sqlx database create
# sqlx migrate add create_subscription_table -- On the shell to create migration file
sqlx migrate run
```

Kelin in DB
return type is kelvin vs celsius