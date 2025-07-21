package models

import "time"

type Product struct {
	ID       int64
	Name     string
	Time     time.Time
	Price    float64
	Text     string
	Username string
}
