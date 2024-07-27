package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"taskapi/logger"
	"taskapi/models"

	"github.com/gorilla/mux"
)

func GetUsersHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		offset, _ := strconv.Atoi(query.Get("offset"))
		limit, _ := strconv.Atoi(query.Get("limit"))

		users, err := models.GetAllUsers(db, offset, limit)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(users)

		logger.Debug("Retrieved users: %+v", users)
		logger.Info("The list of users has been successfully received")
	}
}

func GetUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		user, err := models.GetUserByID(db, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(user)

		logger.Debug("Retrieved user: %+v", user)
		logger.Info("The user was successfully received")
	}
}

func CreateUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Name           string `json:"name"`
			PassportNumber string `json:"passportNumber"`
		}
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		parts := strings.Split(input.PassportNumber, " ")
		if len(parts) != 2 {
			http.Error(w, "Invalid passport format", http.StatusBadRequest)
			return
		}

		name := input.Name
		passportSeries := parts[0]
		passportNumber := parts[1]

		err = models.CreateUser(db, name, passportSeries, passportNumber)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)

		logger.Debug("Created user: %+v", input)
		logger.Info("The user was successfully created")
	}
}

func UpdateUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		var input struct {
			Name           string `json:"name"`
			PassportNumber string `json:"passportNumber"`
		}
		err = json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		parts := strings.Split(input.PassportNumber, " ")
		if len(parts) != 2 {
			http.Error(w, "Invalid passport format", http.StatusBadRequest)
			return
		}

		passportSeries := parts[0]
		passportNumber := parts[1]

		err = models.UpdateUser(db, userID, input.Name, passportSeries, passportNumber)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)

		logger.Debug("Updated user ID %d with data: %+v", userID, input)
		logger.Info("The user's information has been successfully updated")
	}
}

func DeleteUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		err = models.DeleteUser(db, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)

		logger.Debug("Deleted user ID: %d", userID)
		logger.Info("The user was successfully deleted")
	}
}
