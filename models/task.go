package models

import (
	"database/sql"
	"fmt"
)

type Task struct {
	ID          int
	UserID      int
	Name        string
	Description string
	Status      string
	Rate        int
	Deadline    int
	Earned      float64
}

func CreateTask(db *sql.DB, userID int, description string, rate int, deadline int) error {
	_, err := db.Exec("INSERT INTO tasks (user_id, description, rate, deadline, status, earned) VALUES ($1, $2, $3, $4, 'В ожидании', 0)", userID, description, rate, deadline)
	return err
}

func GetTasksByUser(db *sql.DB, userID int) ([]Task, error) {
	query := `
        SELECT tasks.id, tasks.user_id, users.name, tasks.description, tasks.status, tasks.rate, tasks.deadline, tasks.earned
        FROM tasks
        JOIN users ON tasks.user_id = users.id
        WHERE tasks.user_id = $1
    `
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.UserID, &t.Name, &t.Description, &t.Status, &t.Rate, &t.Deadline, &t.Earned); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func GetActiveTasks(db *sql.DB) ([]Task, error) {
	rows, err := db.Query("SELECT id, rate, earned, deadline FROM tasks WHERE status = 'В работе'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Rate, &t.Earned, &t.Deadline); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func GetFilteredTasks(db *sql.DB, userID int, sortBy string) ([]Task, error) {
	query := `
        SELECT tasks.id, tasks.user_id, users.name, tasks.description, tasks.status, tasks.rate, tasks.deadline, tasks.earned
        FROM tasks
        JOIN users ON tasks.user_id = users.id
    `

	if userID > 0 {
		query += fmt.Sprintf(" WHERE users.id = %d", userID)
	}

	if sortBy == "earned" {
		query += " ORDER BY tasks.earned DESC"
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.UserID, &t.Name, &t.Description, &t.Status, &t.Rate, &t.Deadline, &t.Earned); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func DeleteTask(db *sql.DB, taskID int) error {
	_, err := db.Exec("DELETE FROM tasks WHERE id = $1", taskID)
	return err
}

func UpdateTaskEarned(db *sql.DB, taskID int, earned float64, newDeadline int) error {
	_, err := db.Exec("UPDATE tasks SET earned = $1, deadline = $2 WHERE id = $3", earned, newDeadline, taskID)
	return err
}

func UpdateTaskStatus(db *sql.DB, taskID int, status string) error {
	_, err := db.Exec("UPDATE tasks SET status = $1 WHERE id = $2", status, taskID)
	return err
}
