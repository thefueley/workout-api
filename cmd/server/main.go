package main

import (
	"fmt"
	"net/http"

	"github.com/thefueley/workout-api/internal/comment"
	"github.com/thefueley/workout-api/internal/database"
	transportHTTP "github.com/thefueley/workout-api/internal/transport/http"
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

	commentService := comment.NewService(db)

	handler := transportHTTP.NewHandler(commentService)
	handler.SetupRoutes()

	if err := http.ListenAndServe(":8080", handler.Router); err != nil {
		fmt.Println("Failed to setup server")
		return err
	}

	return nil
}

func main() {
	fmt.Println(`
	_______ _______ __   __ __   __ _______ __    _ _______    _______ _______ ___ 
	|       |       |  |_|  |  |_|  |       |  |  | |       |  |       |       |   |
	|       |   _   |       |       |    ___|   |_| |_     _|  |   _   |    _  |   |
	|      _|  | |  |       |       |   |___|       | |   |    |  |_|  |   |_| |   |
	|     | |  |_|  |       |       |    ___|  _    | |   |    |       |    ___|   |
	|     |_|       | ||_|| | ||_|| |   |___| | |   | |   |    |   _   |   |   |   |
	|_______|_______|_|   |_|_|   |_|_______|_|  |__| |___|    |__| |__|___|   |___|
	
	`)

	app := App{}
	if err := app.Run(); err != nil {
		fmt.Println("Error starting REST API")
		fmt.Println(err)
	}
}
