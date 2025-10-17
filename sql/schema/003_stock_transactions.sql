-- +goose Up
CREATE TABLE stock_transactions (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id)  ON DELETE CASCADE,
    stock_id UUID NOT NULL REFERENCES stocks(id) ON DELETE CASCADE,
    type VARCHAR(10) NOT NULL,      
    quantity NUMERIC(18,6) NOT NULL,           
    price NUMERIC(18,4) NOT NULL,   
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE stock_transactions;