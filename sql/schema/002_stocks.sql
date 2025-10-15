-- +goose Up

-- Enable pgcrypto for UUID generation
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Create the table
CREATE TABLE IF NOT EXISTS stocks (
    id UUID PRIMARY KEY,
    stock_symbol VARCHAR(10) NOT NULL,
    valuation NUMERIC(18, 4)
);

-- Prefill with default stocks
INSERT INTO stocks (id, stock_symbol, valuation) VALUES
(gen_random_uuid(), 'RELIANCE', 2600.50),
(gen_random_uuid(), 'TCS', 3300.75),
(gen_random_uuid(), 'INFY', 1500.25),
(gen_random_uuid(), 'ITC', 430.00);

-- +goose Down

DROP TABLE IF EXISTS stocks;
