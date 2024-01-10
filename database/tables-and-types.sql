/*
 * Copyright (c) 2024.  Jan Eike Suchard. All rights reserved
 * SPDX-License-Identifier: MIT
 */

-- name: 01-create-schema
CREATE SCHEMA IF NOT EXISTS cinema_management;

-- name: 02-table-users
CREATE TABLE IF NOT EXISTS cinema_management.users
(
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    external_id text UNIQUE NOT NULL,
    name        text        NOT NULL,
    active      bool             DEFAULT FALSE
);

-- name: 03-create-logging-table
CREATE TABLE IF NOT EXISTS cinema_management.event_log
(
    id      serial PRIMARY KEY,
    user_id uuid        NOT NULL REFERENCES cinema_management.users (id) MATCH SIMPLE ON UPDATE CASCADE ON DELETE RESTRICT,
    at      timestamptz NOT NULL DEFAULT NOW(),
    event   text        NOT NULL
);

-- name: 04-table-items
CREATE TABLE IF NOT EXISTS cinema_management.items
(
    id           uuid         NOT NULL PRIMARY KEY UNIQUE DEFAULT gen_random_uuid(),
    name         varchar(255) NOT NULL UNIQUE,
    price        numeric                                  DEFAULT 0.00,
    icon         text         NOT NULL,
    issue_ticket bool                                     DEFAULT FALSE
);

-- name: 05-table-registers
CREATE TABLE IF NOT EXISTS cinema_management.registers
(
    id          uuid NOT NULL PRIMARY KEY UNIQUE DEFAULT gen_random_uuid(),
    name        text NOT NULL UNIQUE,
    description text
);

-- name: 06-table-sales
CREATE TABLE IF NOT EXISTS cinema_management.sales
(
    id              uuid        NOT NULL PRIMARY KEY UNIQUE DEFAULT gen_random_uuid(),
    at              timestamptz NOT NULL                    DEFAULT NOW(),
    amount          numeric     NOT NULL                    DEFAULT 0.00,
    items           uuid[]                                  DEFAULT NULL,
    custom_items    jsonb                                   DEFAULT NULL,
    card_payment_id text                                    DEFAULT NULL,
    refunded        bool                                    DEFAULT FALSE
);

-- name: 07-table-transactions
CREATE TABLE IF NOT EXISTS cinema_management.transactions
(
    id          uuid        NOT NULL PRIMARY KEY UNIQUE DEFAULT gen_random_uuid(),
    at          timestamptz NOT NULL                    DEFAULT NOW(),
    amount      numeric     NOT NULL                    DEFAULT 0.00,
    by          uuid        NOT NULL REFERENCES cinema_management.users (id) MATCH SIMPLE ON UPDATE CASCADE ON DELETE RESTRICT,
    title       text        NOT NULL,
    register    uuid                                    DEFAULT NULL REFERENCES cinema_management.registers (id) MATCH SIMPLE ON UPDATE CASCADE ON DELETE RESTRICT,
    description text,
    sale_id     uuid                                    DEFAULT NULL REFERENCES cinema_management.sales (id) MATCH SIMPLE ON UPDATE CASCADE ON DELETE RESTRICT
);

-- name: 08-type-screening
DO
$$
    BEGIN
        CREATE TYPE screening_time AS
        (
            id                 uuid,
            at                 timestamptz,
            available_seats    int,
            allow_reservations bool
        );
    EXCEPTION
        WHEN DUPLICATE_OBJECT THEN NULL;
    END
$$;


-- name: 09-table-movies
CREATE TABLE IF NOT EXISTS cinema_management.movies
(
    id                     uuid NOT NULL PRIMARY KEY UNIQUE DEFAULT gen_random_uuid(),
    title                  text NOT NULL,
    original_title         text NOT NULL,
    description            text,
    screening_times        screening_time[],
    audio_language         text NOT NULL,
    subtitle_language      text                             DEFAULT NULL,
    duration               int                              DEFAULT 0,
    additional_information jsonb                            DEFAULT NULL
);

-- name: 10-table-tickets
CREATE TABLE IF NOT EXISTS cinema_management.tickets
(
    id             serial PRIMARY KEY,
    external_id    uuid        NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    issued_at      timestamptz NOT NULL        DEFAULT NOW(),
    movie          uuid                        DEFAULT NULL REFERENCES cinema_management.movies (id) MATCH SIMPLE ON UPDATE CASCADE ON DELETE CASCADE,
    screening_time uuid                        DEFAULT NULL
);

-- name: 11-table-reservations
CREATE TABLE IF NOT EXISTS cinema_management.reservations
(
    id             uuid PRIMARY KEY UNIQUE DEFAULT gen_random_uuid(),
    movie          uuid                    DEFAULT NULL REFERENCES cinema_management.movies (id) MATCH SIMPLE ON UPDATE CASCADE ON DELETE RESTRICT,
    screening_time uuid                    DEFAULT NULL,
    at             timestamptz NOT NULL    DEFAULT NOW(),
    first_name     text        NOT NULL,
    last_name      text        NOT NULL,
    email_address  text        NOT NULL,
    tickets        int         NOT NULL
)