package main

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	transportHTTP "github.com/thefueley/workout-api/controllers"
	"github.com/thefueley/workout-api/models"
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

	db, err := models.NewDatabase()
	if err != nil {
		return err
	}

	err = models.MigrateDB(db)
	if err != nil {
		return err
	}

	workoutService := models.NewService(db)

	handler := transportHTTP.NewHandler(workoutService)
	handler.SetupRoutes()

	// debug : print token
	transportHTTP.CreateToken()

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
