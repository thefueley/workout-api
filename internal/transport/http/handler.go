package http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/thefueley/workout-api/internal/workout"
)

// Handler : stores pointer to workout service
type Handler struct {
	Router  *mux.Router
	Service *workout.Service
}

// Response : store responses from API
type Response struct {
	Message string
	Error   string
}

// NewHandler : returns pointer to a Handler
func NewHandler(service *workout.Service) *Handler {
	return &Handler{
		Service: service,
	}
}

// LoggingMiddleware - a handy middleware function that logs out incoming requests
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(
			log.Fields{
				"Method": r.Method,
				"Path":   r.URL.Path,
			}).
			Info("handled request")
		next.ServeHTTP(w, r)
	})
}

// SetupRoutes : sets up routes for app
func (h *Handler) SetupRoutes() {
	log.Info("Setting up Routes")
	h.Router = mux.NewRouter()
	h.Router.Use(LoggingMiddleware)

	h.Router.HandleFunc("/api/workout/{id}", h.GetWorkout).Methods("GET")
	h.Router.HandleFunc("/api/workout", h.AddWorkout).Methods("POST")
	h.Router.HandleFunc("/api/workout/{id}", h.UpdateWorkout).Methods("PUT")
	h.Router.HandleFunc("/api/workout/{id}", h.DeleteWorkout).Methods("DELETE")
	h.Router.HandleFunc("/api/workout", h.GetAllWorkouts).Methods("GET")

	h.Router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		if err := sendOkResponse(w, Response{Message: "Hooray for me!"}); err != nil {
			panic(err)
		}
	})
}

// sendOkResponse : send ok response
func sendOkResponse(w http.ResponseWriter, resp interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(resp)
}

// sendErrorResponse : send error response
func sendErrorResponse(w http.ResponseWriter, message string, err error) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusInternalServerError)

	if err := json.NewEncoder(w).Encode(Response{Message: message, Error: err.Error()}); err != nil {
		panic(err)
	}
}
