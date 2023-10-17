#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "campaign" <<-EOSQL
    CREATE TABLE IF NOT EXISTS settlement
    (
        id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
        owner uuid NOT NULL,
        name character varying(50) NOT NULL,
        survivalLimit integer,
        departingSurvival integer,
        collectiveCognition integer,
        currentYear integer,
    );
EOSQL