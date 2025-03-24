package models

import "github.com/google/uuid"

type Verse struct {
	Id          uuid.UUID `json:"id"`
	SongId      uuid.UUID `json:"song_id"`
	VerseNumber int       `json:"verse_number"`
	Text        string    `json:"text"`
}
