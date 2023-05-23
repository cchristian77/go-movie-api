CREATE TABLE IF NOT EXISTS ratings
(
    id         SERIAL PRIMARY KEY,
    uuid       UUID               DEFAULT gen_random_uuid() UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    movie_id   INTEGER   NOT NULL REFERENCES movies (id),
    rating     NUMERIC(3, 2) CHECK ( rating > 0 ),
    comment    TEXT
);

