package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"taskapi/models"

	"taskapi/logger"

	"github.com/gorilla/mux"
)

func CreateTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Debug("Received request to create task")
		var task models.Task
		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			logger.Info("Error decoding JSON: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		taskRate := int(task.Rate)
		logger.Debug("Decoded task: %+v", task)

		err = models.CreateTask(db, task.UserID, task.Description, taskRate, task.Deadline)
		if err != nil {
			logger.Info("Error creating task: %v", err)
			if err.Error() == "user does not exist" {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)

		logger.Info("Task has been successfully created: %+v", task)
	}
}

func GetTasksHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Debug("Received request to get tasks")
		query := r.URL.Query()
		userIDStr := query.Get("user_id")
		var userID int
		if userIDStr != "" {
			var err error
			userID, err = strconv.Atoi(userIDStr)
			if err != nil {
				logger.Info("Invalid user ID: %v", err)
				http.Error(w, "Invalid user ID", http.StatusBadRequest)
				return
			}
		}

		sortBy := query.Get("sort_by")
		logger.Debug("Query parameters - user_id: %d, sort_by: %s", userID, sortBy)

		tasks, err := models.GetFilteredTasks(db, userID, sortBy)
		if err != nil {
			logger.Info("Error getting tasks: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(tasks)

		logger.Debug("Retrieved tasks: %+v", tasks)
		logger.Info("Tasks have been successfully received")
	}
}

func DeleteTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Debug("Received request to delete task")
		if r.Method == "DELETE" {
			vars := mux.Vars(r)
			taskID, err := strconv.Atoi(vars["id"])
			if err != nil {
				logger.Info("Invalid task ID: %v", err)
				http.Error(w, "Invalid task ID", http.StatusBadRequest)
				return
			}
			err = models.DeleteTask(db, taskID)
			if err != nil {
				logger.Info("Error deleting task: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)

			logger.Info("Task has been successfully deleted: TaskID %d", taskID)
		} else {
			logger.Info("Invalid request method")
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	}
}

func StartTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Debug("Received request to start task")
		vars := mux.Vars(r)
		taskID, err := strconv.Atoi(vars["id"])
		if err != nil {
			logger.Info("Invalid task ID: %v", err)
			http.Error(w, "Invalid task ID", http.StatusBadRequest)
			return
		}

		err = models.UpdateTaskStatus(db, taskID, "В работе")
		if err != nil {
			logger.Info("Error starting task: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)

		logger.Info("Task has been successfully started: TaskID %d", taskID)
	}
}

func StopTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Debug("Received request to stop task")
		vars := mux.Vars(r)
		taskID, err := strconv.Atoi(vars["id"])
		if err != nil {
			logger.Info("Invalid task ID: %v", err)
			http.Error(w, "Invalid task ID", http.StatusBadRequest)
			return
		}

		err = models.UpdateTaskStatus(db, taskID, "Завершена")
		if err != nil {
			logger.Info("Error stopping task: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)

		logger.Info("Task has been successfully stopped: TaskID %d", taskID)
	}
}
