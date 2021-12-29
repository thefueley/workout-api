package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/thefueley/workout-api/models"
)

// GetWorkout : get workout
func (h *Handler) GetWorkout(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pKey := vars["pKey"]
	rKey := vars["rKey"]

	workout, err := h.Service.GetWorkout(pKey, rKey)
	if err != nil {
		sendErrorResponse(w, "Error retrieving workout by pKey and rKey.\nPunt", err)
		return
	}

	if err := sendOkResponse(w, workout); err != nil {
		panic(err)
	}
}

// AddWorkout : add a workout
func (h *Handler) AddWorkout(w http.ResponseWriter, r *http.Request) {
	var unmarshaledWorkout models.Workout
	if err := json.NewDecoder(r.Body).Decode(&unmarshaledWorkout); err != nil {
		sendErrorResponse(w, "Failed to decode JSON Body.\nPunt!", err)
		return
	}

	op, err := h.Service.AddWorkout(unmarshaledWorkout)

	if err != nil {
		sendErrorResponse(w, "Error adding workout.\nPunt!", err)
		return
	}

	if err := sendOkResponse(w, op); err != nil {
		panic(err)
	}
}

// UpdateWorkout : update workout
func (h *Handler) UpdateWorkout(w http.ResponseWriter, r *http.Request) {
	var unmarshalledWorkout models.Workout

	if err := json.NewDecoder(r.Body).Decode(&unmarshalledWorkout); err != nil {
		sendErrorResponse(w, "Failed to decode JSON Body.\nPunt!", err)
		return
	}

	vars := mux.Vars(r)
	pKey := vars["pKey"]
	rKey := vars["rKey"]

	_, err := h.Service.GetWorkout(pKey, rKey)
	if err != nil {
		sendErrorResponse(w, "Error retrieving workout by pKey and rKey.\nPunt", err)
		return
	}

	op, err := h.Service.UpdateWorkout(pKey, rKey, unmarshalledWorkout)

	if err != nil {
		sendErrorResponse(w, "Error updating workout.\nPunt!", err)
		return
	}

	if err := sendOkResponse(w, op); err != nil {
		panic(err)
	}
}

// DeleteWorkout : delete workout
func (h *Handler) DeleteWorkout(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pKey := vars["pKey"]
	rKey := vars["rKey"]

	err := h.Service.DeleteWorkout(pKey, rKey)
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
