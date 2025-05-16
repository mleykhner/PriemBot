package models

import "time"

type Article struct {
	ID        uint
	Title     string
	Subtitle  string
	ImageURL  string
	CreatedAt time.Time
}
