package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/thefueley/workout-api/internal/comment"
)

// Handler : stores pointer to your comments service
type Handler struct {
	Router  *mux.Router
	Service *comment.Service
}

// Response : store responses from API
type Response struct {
	Message string
	Error   string
}

// NewHandler : returns pointer to a Handler
func NewHandler(service *comment.Service) *Handler {
	return &Handler{
		Service: service,
	}
}

// SetupRoutes : sets up routes for app
func (h *Handler) SetupRoutes() {
	fmt.Println("Setting up Routes")
	h.Router = mux.NewRouter()

	h.Router.HandleFunc("/api/comment", h.GetAllComments).Methods("GET")
	h.Router.HandleFunc("/api/comment", h.PostComment).Methods("POST")
	h.Router.HandleFunc("/api/comment/{id}", h.GetComment).Methods("GET")
	h.Router.HandleFunc("/api/comment/{id}", h.UpdateComment).Methods("PUT")
	h.Router.HandleFunc("/api/comment/{id}", h.DeleteComment).Methods("DELETE")

	h.Router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		if err := sendOkResponse(w, Response{Message: "Hooray for me!"}); err != nil {
			panic(err)
		}
	})
}

// GetAllComments : get all comments
func (h *Handler) GetAllComments(w http.ResponseWriter, r *http.Request) {
	comments, err := h.Service.GetAllComments()

	if err != nil {
		sendErrorResponse(w, "Failed to retrieve all comments.\nPunt!", err)
		return
	}

	if err := sendOkResponse(w, comments); err != nil {
		panic(err)
	}
}

// PostComment : post a comment to db
func (h *Handler) PostComment(w http.ResponseWriter, r *http.Request) {
	var cmt comment.Comment
	if err := json.NewDecoder(r.Body).Decode(&cmt); err != nil {
		sendErrorResponse(w, "Failed to decode JSON Body.\nPunt!", err)
		return
	}

	cmt, err := h.Service.PostComment(cmt)

	if err != nil {
		sendErrorResponse(w, "Failed ot post comment.\nPunt!", err)
		return
	}

	if err := sendOkResponse(w, cmt); err != nil {
		panic(err)
	}
}

// GetComment : get comment by ID
func (h *Handler) GetComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	i, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		sendErrorResponse(w, "Error parsing UINT from ID string.\nPunt!", err)
		return
	}

	cmt, err := h.Service.GetComment(uint(i))
	if err != nil {
		sendErrorResponse(w, "Error retrieving comment by ID.\nPunt", err)
		return
	}

	if err := sendOkResponse(w, cmt); err != nil {
		panic(err)
	}
}

// UpdateComment : update comment in db
func (h *Handler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	var cmt comment.Comment
	if err := json.NewDecoder(r.Body).Decode(&cmt); err != nil {
		sendErrorResponse(w, "Failed to decode JSON Body.\nPunt!", err)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	commentID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		sendErrorResponse(w, "Error parsing UINT from ID string.\nPunt!", err)
		return
	}

	cmt, err = h.Service.UpdateComment(uint(commentID), cmt)

	if err != nil {
		sendErrorResponse(w, "Failed to update comment.\nPunt!", err)
		return
	}

	if err := sendOkResponse(w, cmt); err != nil {
		panic(err)
	}
}

// DeleteComment : delete comment from db
func (h *Handler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	commentID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		sendErrorResponse(w, "Error parsing UINT from ID string.\nPunt!", err)
		return
	}

	err = h.Service.DeleteComment(uint(commentID))
	if err != nil {
		sendErrorResponse(w, "Error deleting comment.\nPunt!", err)
		return
	}

	if err = sendOkResponse(w, Response{Message: "Poof, message deleted!"}); err != nil {
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
