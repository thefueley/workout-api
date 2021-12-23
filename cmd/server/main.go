package main

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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
	appInf := []string{"AppName:", app.Name, "AppVersion:", app.Version}
	log.Info().Strs("APP INFO: ", appInf).Msg("")
	log.Info().Msg("Setting up API")

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
		log.Error().Msg("Failed to setup server")
		return err
	}

	return nil
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	fmt.Println(`
	.-------------.
	| WORKOUT API |
	'-------------'
	`)

	app := App{
		Name:    "Workout API",
		Version: "1.0.0",
	}
	if err := app.Run(); err != nil {
		log.Error().Msg("Error starting REST API")
		log.Fatal().Msg(err.Error())
	}
}
