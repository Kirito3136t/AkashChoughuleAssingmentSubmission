-- +goose Up

CREATE TABLE users (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL ,
    isReferral BOOLEAN DEFAULT FALSE,
    referralUserId UUID REFERENCES users(id)
);

-- +goose Down
DROP TABLE users;
