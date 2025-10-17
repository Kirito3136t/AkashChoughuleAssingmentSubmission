-- +goose Up

CREATE TABLE users (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL ,
    password TEXT NOT NULL
);

-- +goose Down
DROP TABLE users;
