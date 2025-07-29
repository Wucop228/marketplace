package repo

import (
	"database/sql"
	"time"
)

func CreateProduct(db *sql.DB, price float64, header string, text string, username string, image_url string) error {
	query := "INSERT INTO products (time, price, header, text, username, image_url) values ($1, $2, $3, $4, $5, $6)"
	_, err := db.Exec(query, time.Now(), price, header, text, username, image_url)

	return err
}
