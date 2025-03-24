-- +goose Up
CREATE TABLE IF NOT EXISTS songs (
                                     id UUID PRIMARY KEY,
                                     group_id UUID NOT NULL,
                                     group_name VARCHAR(255),
    title VARCHAR(255) NOT NULL,
    release_date DATE NOT NULL,
    text TEXT NOT NULL,
    link TEXT NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (group_id) REFERENCES groups(id)
    );

-- +goose Down
DROP TABLE IF EXISTS songs;