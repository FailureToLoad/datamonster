#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE TABLE IF NOT EXISTS item (
        name VARCHAR (50) NOT NULL PRIMARY KEY,
        keywords VARCHAR (150) NOT NULL,
        source VARCHAR (25)
    );

    CREATE TABLE IF NOT EXISTS resource
    (
        name VARCHAR (50) NOT NULL PRIMARY KEY,
        keywords VARCHAR (150),
        source VARCHAR (50),
        strange boolean
    );

    CREATE TABLE If NOT EXISTS location
    (
        id VARCHAR (50) NOT NULL PRIMARY KEY,
        name VARCHAR (150) NOT NULL
    );

    CREATE TABLE IF NOT EXISTS innovation
    (
        name VARCHAR (50) NOT NULL PRIMARY KEY,
        source VARCHAR (25) NOT NULL,
        keywords VARCHAR (150) 
    );
EOSQL

psql -v ON_ERROR_STOP=1 -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c "\copy item FROM /seed/items.csv WITH (FORMAT CSV)"
psql -v ON_ERROR_STOP=1 -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c "\copy resource FROM /seed/resources.csv WITH (FORMAT CSV)"
psql -v ON_ERROR_STOP=1 -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c "\copy location FROM /seed/locations.csv WITH (FORMAT CSV)"
psql -v ON_ERROR_STOP=1 -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c "\copy innovation FROM /seed/innovations.csv WITH (FORMAT CSV)"