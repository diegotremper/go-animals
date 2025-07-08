CREATE TABLE IF NOT EXISTS animals (
    id BIGSERIAL NOT NULL primary key,
    name VARCHAR NOT NULL,
    age INT default 0,
    description TEXT
);
