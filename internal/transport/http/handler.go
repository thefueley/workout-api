package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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

// SetupRoutes : sets up routes for app
func (h *Handler) SetupRoutes() {
	fmt.Println("Setting up Routes")
	h.Router = mux.NewRouter()

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

// GetWorkout : get workout
func (h *Handler) GetWorkout(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	i, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		sendErrorResponse(w, "Error parsing UINT from ID string.\nPunt!", err)
		return
	}

	cmt, err := h.Service.GetWorkout(uint(i))
	if err != nil {
		sendErrorResponse(w, "Error retrieving workout by ID.\nPunt", err)
		return
	}

	if err := sendOkResponse(w, cmt); err != nil {
		panic(err)
	}
}

// AddWorkout : add a workout
func (h *Handler) AddWorkout(w http.ResponseWriter, r *http.Request) {
	var work workout.Workout
	if err := json.NewDecoder(r.Body).Decode(&work); err != nil {
		sendErrorResponse(w, "Failed to decode JSON Body.\nPunt!", err)
		return
	}

	cmt, err := h.Service.AddWorkout(work)

	if err != nil {
		sendErrorResponse(w, "Error adding workout.\nPunt!", err)
		return
	}

	if err := sendOkResponse(w, cmt); err != nil {
		panic(err)
	}
}

// UpdateWorkout : update workout
func (h *Handler) UpdateWorkout(w http.ResponseWriter, r *http.Request) {
	var work workout.Workout
	if err := json.NewDecoder(r.Body).Decode(&work); err != nil {
		sendErrorResponse(w, "Failed to decode JSON Body.\nPunt!", err)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	workoutID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		sendErrorResponse(w, "Error parsing UINT from ID string.\nPunt!", err)
		return
	}

	work, err = h.Service.UpdateComment(uint(workoutID), work)

	if err != nil {
		sendErrorResponse(w, "Error updating workout.\nPunt!", err)
		return
	}

	if err := sendOkResponse(w, work); err != nil {
		panic(err)
	}
}

// DeleteWorkout : delete workout
func (h *Handler) DeleteWorkout(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	workoutID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		sendErrorResponse(w, "Error parsing UINT from ID string.\nPunt!", err)
		return
	}

	err = h.Service.DeleteWorkout(uint(workoutID))
	if err != nil {
		sendErrorResponse(w, "Error deleting workout.\nPunt!", err)
		return
	}

	if err = sendOkResponse(w, Response{Message: "Poof! Workout deleted."}); err != nil {
		panic(err)
	}
}

// GetAllWorkouts : get all workouts
func (h *Handler) GetAllWorkouts(w http.ResponseWriter, r *http.Request) {
	workouts, err := h.Service.GetAllWorkouts()

	if err != nil {
		sendErrorResponse(w, "Error retrieving all workouts.\nPunt!", err)
		return
	}

	if err := sendOkResponse(w, workouts); err != nil {
		panic(err)
	}
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
