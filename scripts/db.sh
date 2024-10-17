#!/bin/bash
set -e
export PGPASSWORD=password;
psql -v ON_ERROR_STOP=1 --username "postgres" --dbname "orderbook" <<-EOSQL
  CREATE DATABASE orderbook;
  GRANT ALL PRIVILEGES ON DATABASE orderbook TO "postgres";
EOSQL