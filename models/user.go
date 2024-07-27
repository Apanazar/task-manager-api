package models

import (
	"database/sql"
)

type User struct {
	ID             int
	Name           string
	PassportSeries string
	PassportNumber string
}

func GetAllUsers(db *sql.DB, offset int, limit int) ([]User, error) {
	rows, err := db.Query("SELECT id, name, passport_series, passport_number FROM users LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.PassportSeries, &u.PassportNumber); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func GetUserByID(db *sql.DB, userID int) (*User, error) {
	row := db.QueryRow("SELECT id, name, passport_series, passport_number FROM users WHERE id = $1", userID)

	var u User
	if err := row.Scan(&u.ID, &u.Name, &u.PassportSeries, &u.PassportNumber); err != nil {
		return nil, err
	}
	return &u, nil
}

func CreateUser(db *sql.DB, name, series, number string) error {
	_, err := db.Exec("INSERT INTO users (name, passport_series, passport_number) VALUES ($1, $2, $3)", name, series, number)
	return err
}

func UpdateUser(db *sql.DB, userID int, name, series, number string) error {
	_, err := db.Exec("UPDATE users SET name = $1, passport_series = $2, passport_number = $3 WHERE id = $4", name, series, number, userID)
	return err
}

func DeleteUser(db *sql.DB, userID int) error {
	_, err := db.Exec("DELETE FROM users WHERE id = $1", userID)
	return err
}
