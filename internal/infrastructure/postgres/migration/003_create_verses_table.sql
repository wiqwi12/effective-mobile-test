-- +goose Up
CREATE TABLE IF NOT EXISTS verses (
                                      id UUID PRIMARY KEY,
                                      song_id UUID NOT NULL,
                                      verse_number INTEGER NOT NULL,
                                      text TEXT NOT NULL,
                                      FOREIGN KEY (song_id) REFERENCES songs(id) ON DELETE CASCADE
    );

-- +goose Down
DROP TABLE IF EXISTS verses;