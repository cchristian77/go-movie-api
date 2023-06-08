CREATE TABLE IF NOT EXISTS users
(
    id                  SERIAL PRIMARY KEY,
    uuid                UUID                         DEFAULT gen_random_uuid() UNIQUE,
    created_at          TIMESTAMP           NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP           NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at          TIMESTAMP,
    username            VARCHAR(255) UNIQUE NOT NULL,
    email               VARCHAR(255) UNIQUE NOT NULL,
    full_name           VARCHAR(255)        NOT NULL,
    password            TEXT                NOT NULL,
    is_admin            BOOLEAN             NOT NULL DEFAULT false,
    is_email_verified   bool                NOT NULL DEFAULT false,
    password_changed_at timestamptz         NOT NULL DEFAULT '0001-01-01'
);

CREATE TABLE IF NOT EXISTS sessions
(
    id            UUID PRIMARY KEY,
    user_id       INTEGER REFERENCES users (id) NOT NULL,
    refresh_token TEXT                          NOT NULL,
    user_agent    VARCHAR(255)                  NOT NULL,
    client_ip     VARCHAR(255)                  NOT NULL,
    is_blocked    BOOLEAN                       NOT NULL DEFAULT false,
    expires_at    TIMESTAMPTZ                   NOT NULL,
    created_at    TIMESTAMPTZ                   NOT NULL DEFAULT (now())
);