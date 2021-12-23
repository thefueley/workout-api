package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/thefueley/workout-api/internal/workout"
)

// GetWorkout : get workout
func (h *Handler) GetWorkout(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	i, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		sendErrorResponse(w, "Error parsing UINT from ID string.\nPunt!", err)
		return
	}

	cmt, err := h.Service.GetWorkout(h.Service.DB, uint(i))
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

	work, err = h.Service.UpdateWorkout(uint(workoutID), work)

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
	workouts, err := h.Service.GetAllWorkouts(h.Service.DB)

	if err != nil {
		sendErrorResponse(w, "Error retrieving all workouts.\nPunt!", err)
		return
	}

	if err := sendOkResponse(w, workouts); err != nil {
		panic(err)
	}
}
