-- +goose Up
ALTER TABLE customers 
    ADD COLUMN is_email_verified BOOLEAN DEFAULT FALSE , 
    ADD COLUMN is_phone_number_verified BOOLEAN DEFAULT FALSE;

-- +goose Down