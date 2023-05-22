CREATE TABLE IF NOT EXISTS movies
(
    id         SERIAL PRIMARY KEY,
    uuid       UUID                  DEFAULT gen_random_uuid() UNIQUE,
    created_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    title      VARCHAR(255) NOT NULL,
    duration   INTEGER      NOT NULL,
    year       INTEGER      NOT NULL,
    synopsis   TEXT
);

comment on column movies.duration is 'minutes';

CREATE TABLE IF NOT EXISTS genres
(
    id         SERIAL PRIMARY KEY,
    uuid       UUID                  DEFAULT gen_random_uuid() UNIQUE,
    created_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    name       VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS movie_genres
(
    movie_id INTEGER REFERENCES movies (id),
    genre_id INTEGER REFERENCES genres (id),
    PRIMARY KEY (movie_id, genre_id)
);