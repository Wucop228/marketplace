package models

import "time"

type Product struct {
	ID       int64
	Time     time.Time
	Price    float64
	Header   string
	Text     string
	Username string
	ImageURL string
}
