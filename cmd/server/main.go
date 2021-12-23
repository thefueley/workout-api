package main

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/thefueley/workout-api/internal/database"
	transportHTTP "github.com/thefueley/workout-api/internal/transport/http"
	"github.com/thefueley/workout-api/internal/workout"
)

// App : app info
type App struct {
	Name    string
	Version string
}

// Run : Sets up app
func (app *App) Run() error {
	log.SetFormatter(&log.JSONFormatter{})
	log.WithFields(
		log.Fields{
			"AppName":    app.Name,
			"AppVersion": app.Version,
		}).Info("Setting up API")

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
		log.Error("Failed to setup server")
		return err
	}

	return nil
}

func main() {
	log.Info(`
	.-------------.
	| WORKOUT API |
	'-------------'
	`)

	app := App{
		Name:    "Workout API",
		Version: "1.0.0",
	}
	if err := app.Run(); err != nil {
		log.Error("Error starting REST API")
		log.Fatal(err)
	}
}
