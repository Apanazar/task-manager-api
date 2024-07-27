package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"taskapi/handlers"
	"taskapi/models"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // добавьте этот импорт
	httpSwagger "github.com/swaggo/http-swagger"
)

//	@title			Task API
//	@version		1.0
//	@description	This is a sample server for managing tasks.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

// @host		localhost:8080
// @BasePath	/
// @schemes	http
func main() {
	err := godotenv.Load("./config/config.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	apiPort := os.Getenv("API_PORT")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")

	db, err := sql.Open("postgres", fmt.Sprintf("dbname=%s user=%s password=%s sslmode=disable", dbName, dbUser, dbPassword))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := runMigrations(db); err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()

	//	@Summary		Get all users
	//	@Description	Get all users with pagination
	//	@Tags			Users
	//	@Accept			json
	//	@Produce		json
	//	@Param			offset	query		int	false	"Offset"	default(0)
	//	@Param			limit	query		int	false	"Limit"		default(100)
	//	@Success		200		{array}		models.User
	//	@Failure		400		{object}	http.Response
	//	@Failure		500		{object}	http.Response
	//	@Router			/api/users [get]
	router.HandleFunc("/api/users", handlers.GetUsersHandler(db)).Methods("GET")

	//	@Summary		Get a user by ID
	//	@Description	Get a user by ID
	//	@Tags			Users
	//	@Accept			json
	//	@Produce		json
	//	@Param			id	path		int	true	"User ID"
	//	@Success		200	{object}	models.User
	//	@Failure		400	{object}	http.Response
	//	@Failure		404	{object}	http.Response
	//	@Failure		500	{object}	http.Response
	//	@Router			/api/users/{id} [get]
	router.HandleFunc("/api/users/{id:[0-9]+}", handlers.GetUserHandler(db)).Methods("GET")

	//	@Summary		Create a new user
	//	@Description	Create a new user
	//	@Tags			Users
	//	@Accept			json
	//	@Produce		json
	//	@Param			user	body		models.User	true	"User"
	//	@Success		201		{object}	models.User
	//	@Failure		400		{object}	http.Response
	//	@Failure		500		{object}	http.Response
	//	@Router			/api/users [post]
	router.HandleFunc("/api/users", handlers.CreateUserHandler(db)).Methods("POST")

	//	@Summary		Update a user
	//	@Description	Update a user
	//	@Tags			Users
	//	@Accept			json
	//	@Produce		json
	//	@Param			id		path		int			true	"User ID"
	//	@Param			user	body		models.User	true	"User"
	//	@Success		200		{object}	models.User
	//	@Failure		400		{object}	http.Response
	//	@Failure		404		{object}	http.Response
	//	@Failure		500		{object}	http.Response
	//	@Router			/api/users/{id} [put]
	router.HandleFunc("/api/users/{id:[0-9]+}", handlers.UpdateUserHandler(db)).Methods("PUT")

	//	@Summary		Delete a user
	//	@Description	Delete a user
	//	@Tags			Users
	//	@Accept			json
	//	@Produce		json
	//	@Param			id	path	int	true	"User ID"
	//	@Success		204
	//	@Failure		400	{object}	http.Response
	//	@Failure		404	{object}	http.Response
	//	@Failure		500	{object}	http.Response
	//	@Router			/api/users/{id} [delete]
	router.HandleFunc("/api/users/{id:[0-9]+}", handlers.DeleteUserHandler(db)).Methods("DELETE")

	//	@Summary		Get all tasks
	//	@Description	Get all tasks
	//	@Tags			Tasks
	//	@Accept			json
	//	@Produce		json
	//	@Success		200	{array}		models.Task
	//	@Failure		400	{object}	http.Response
	//	@Failure		500	{object}	http.Response
	//	@Router			/api/tasks [get]
	router.HandleFunc("/api/tasks", handlers.GetTasksHandler(db)).Methods("GET")

	//	@Summary		Create a new task
	//	@Description	Create a new task
	//	@Tags			Tasks
	//	@Accept			json
	//	@Produce		json
	//	@Param			task	body		models.Task	true	"Task"
	//	@Success		201		{object}	models.Task
	//	@Failure		400		{object}	http.Response
	//	@Failure		500		{object}	http.Response
	//	@Router			/api/tasks [post]
	router.HandleFunc("/api/tasks", handlers.CreateTaskHandler(db)).Methods("POST")

	//	@Summary		Start a task
	//	@Description	Start a task
	//	@Tags			Tasks
	//	@Accept			json
	//	@Produce		json
	//	@Param			id	path		int	true	"Task ID"
	//	@Success		200	{object}	models.Task
	//	@Failure		400	{object}	http.Response
	//	@Failure		404	{object}	http.Response
	//	@Failure		500	{object}	http.Response
	//	@Router			/api/tasks/{id}/start [post]
	router.HandleFunc("/api/tasks/{id:[0-9]+}/start", handlers.StartTaskHandler(db)).Methods("POST")

	//	@Summary		Stop a task
	//	@Description	Stop a task
	//	@Tags			Tasks
	//	@Accept			json
	//	@Produce		json
	//	@Param			id	path		int	true	"Task ID"
	//	@Success		200	{object}	models.Task
	//	@Failure		400	{object}	http.Response
	//	@Failure		404	{object}	http.Response
	//	@Failure		500	{object}	http.Response
	//	@Router			/api/tasks/{id}/stop [post]
	router.HandleFunc("/api/tasks/{id:[0-9]+}/stop", handlers.StopTaskHandler(db)).Methods("POST")

	//	@Summary		Delete a task
	//	@Description	Delete a task
	//	@Tags			Tasks
	//	@Accept			json
	//	@Produce		json
	//	@Param			id	path	int	true	"Task ID"
	//	@Success		204
	//	@Failure		400	{object}	http.Response
	//	@Failure		404	{object}	http.Response
	//	@Failure		500	{object}	http.Response
	//	@Router			/api/tasks/{id} [delete]
	router.HandleFunc("/api/tasks/{id:[0-9]+}", handlers.DeleteTaskHandler(db)).Methods("DELETE")

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	go func() {
		for {
			updateEarnedMoney(db)
			time.Sleep(1 * time.Minute)
		}
	}()

	log.Println("The API server is running")
	http.ListenAndServe(":"+apiPort, router)
}

func runMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("could not create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("could not run up migrations: %w", err)
	}

	return nil
}

func updateEarnedMoney(db *sql.DB) {
	tasks, err := models.GetActiveTasks(db)
	if err != nil {
		log.Println("Error getting active tasks:", err)
		return
	}

	for _, task := range tasks {
		ratePerMinute := float64(task.Rate / 60.0)
		newEarned := task.Earned + ratePerMinute

		formattedEarned, err := strconv.ParseFloat(fmt.Sprintf("%.2f", newEarned), 64)
		if err != nil {
			log.Println("Error formatting earned value:", err)
			continue
		}

		if task.Deadline > 0 {
			task.Deadline--
			err := models.UpdateTaskEarned(db, task.ID, formattedEarned, task.Deadline)
			if err != nil {
				log.Println("Error updating task earned:", err)
			}
		}
	}
}
