-- +goose Up
CREATE DATABASE IF NOT EXISTS goa;

-- +goose Down
DROP DATABASE goa;
