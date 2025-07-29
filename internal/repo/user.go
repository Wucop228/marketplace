package repo

import (
	"database/sql"
	"github.com/Wucop228/marketplace/internal/models"
	_ "github.com/lib/pq"
)

func UserExists(db *sql.DB, username string) (bool, error) {
	query := "SELECT COUNT(*) FROM users WHERE username=$1"
	var count int
	err := db.QueryRow(query, username).Scan(&count)
	if err != nil {
		return false, err
	}
	if count != 0 {
		return true, nil
	}

	return false, nil
}

func CreateUser(db *sql.DB, user *models.User) error {
	query := "INSERT INTO users (id, username, role, password_hash) VALUES (DEFAULT, $1, $2, $3)"
	_, err := db.Exec(query, user.Username, "user", user.PasswordHash)
	if err != nil {
		return err
	}

	var id int64
	err = db.QueryRow("SELECT id FROM users WHERE username=$1", user.Username).Scan(&id)
	user.ID = id
	return err
}

func LoginUser(db *sql.DB, username string) (models.User, error) {
	query := "SELECT id, username, role, password_hash FROM users WHERE username=$1"
	var user models.User
	err := db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Role, &user.PasswordHash)
	return user, err
}

func GetUserById(db *sql.DB, id int64) (string, error) {
	query := "SELECT username FROM users WHERE id=$1"
	var username string
	err := db.QueryRow(query, id).Scan(&username)
	return username, err
}
