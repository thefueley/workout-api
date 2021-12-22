package main

import (
	"fmt"
	"net/http"

	"github.com/thefueley/workout-api/internal/database"
	transportHTTP "github.com/thefueley/workout-api/internal/transport/http"
	"github.com/thefueley/workout-api/internal/workout"
)

// db connections
type App struct{}

// Run : Sets up app
func (app *App) Run() error {
	fmt.Println("Setting up APP")

	var err error

	db, err := database.NewDatabase()
	if err != nil {
		return err
	}

	err = database.MigrateDB(db)
	if err != nil {
		return err
	}

	workoutService := workout.NewService(db)

	handler := transportHTTP.NewHandler(workoutService)
	handler.SetupRoutes()

	if err := http.ListenAndServe(":8080", handler.Router); err != nil {
		fmt.Println("Failed to setup server")
		return err
	}

	return nil
}

func main() {
	fmt.Println(`
	.-------------.
	| WORKOUT API |
	'-------------'
	`)

	app := App{}
	if err := app.Run(); err != nil {
		fmt.Println("Error starting REST API")
		fmt.Println(err)
	}
}
