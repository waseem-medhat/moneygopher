-- +goose Up
CREATE TABLE accounts (
    id TEXT PRIMARY KEY NOT NULL,
    phone_number TEXT NOT NULL,
    balance_dollars INTEGER NOT NULL DEFAULT 0
);

-- +goose Down
DROP TABLE accounts;
