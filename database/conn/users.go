package conn

import (
	"database/sql"
	"errors"

	"github.com/A-Victory/blog/models"
)

func (db *DB) SaveUser(data models.User) (int, error) {

	query := "INSERT INTO Users (username, email, password) VALUES (?, ?, ?)"

	result, err := db.Conn.DB.Exec(query, data.Username, data.Email, data.Password)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (db *DB) GetUser(identifierType string, value interface{}) (models.User, error) {
	var query string

	switch identifierType {
	case "id":
		query = "SELECT id, username, email, password FROM users WHERE id = ?"
	case "username":
		query = "SELECT id, username, email, password FROM users WHERE username = ?"
	case "email":
		query = "SELECT id, username, email, password FROM users WHERE email = ?"
	default:
		return models.User{}, errors.New("invalid identifier type")
	}
	user := models.User{}

	err := db.Conn.DB.QueryRow(query, value).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, sql.ErrNoRows
		}
		return models.User{}, err
	}
	return user, nil

}
