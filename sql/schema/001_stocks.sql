-- +goose Up
CREATE TABLE stocks (
    id UUID PRIMARY KEY,
    stock_symbol VARCHAR(10) NOT NULL,
    valuation NUMERIC(18, 4)
);


-- +goose Down
DROP TABLE stocks