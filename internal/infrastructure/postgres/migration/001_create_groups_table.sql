-- +goose Up
CREATE TABLE IF NOT EXISTS groups (
                                      id UUID PRIMARY KEY,
                                      name VARCHAR(255) NOT NULL
    );

-- +goose Down
DROP TABLE IF EXISTS groups;