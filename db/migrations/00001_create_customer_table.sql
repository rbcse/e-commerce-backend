-- +goose Up
CREATE TABLE customers (
    customer_id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    phone_number VARCHAR(14) NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE IF EXISTS customers;