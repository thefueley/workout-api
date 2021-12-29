package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/thefueley/workout-api/models"
)

// Handler : stores pointer to workout service
type Handler struct {
	Router  *mux.Router
	Service *models.Service
}

// Response : store responses from API
type Response struct {
	Message string
	Error   string
}

// NewHandler : returns pointer to a Handler
func NewHandler(service *models.Service) *Handler {
	return &Handler{
		Service: service,
	}
}

// LoggingMiddleware - a handy middleware function that logs out incoming requests
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ev := []string{"Method", r.Method, "Path", r.URL.Path}

		log.Info().Strs("API", ev).Msg("")
		next.ServeHTTP(w, r)
	})
}

// validateToken - validates an incoming jwt token
func validateToken(accessToken string) bool {
	// replace this by loading in a private RSA cert for more security
	var mySigningKey = []byte("missionimpossible")
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there was an error")
		}
		return mySigningKey, nil
	})

	if err != nil {
		return false
	}

	return token.Valid
}

// JWTAuth - a handy middleware function that will provide basic auth around specific endpoints
func JWTAuth(original func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info().Msg("jwt auth endpoint hit")
		authHeader := r.Header["Authorization"]
		if authHeader == nil {
			sendErrorResponse(w, "not authorized", errors.New("not authorized"))
			return
		}

		authHeaderParts := strings.Split(authHeader[0], " ")
		if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
			sendErrorResponse(w, "not authorized", errors.New("not authorized"))
			return
		}

		if validateToken(authHeaderParts[1]) {
			original(w, r)
		} else {
			sendErrorResponse(w, "not authorized", errors.New("not authorized"))
			return
		}
	}
}

// SetupRoutes : sets up routes for app
func (h *Handler) SetupRoutes() {
	log.Info().Msg("Setting up Routes")
	h.Router = mux.NewRouter()
	h.Router.Use(LoggingMiddleware)

	h.Router.HandleFunc("/api/workout/{pKey}/{rKey}", h.GetWorkout).Methods("GET")
	h.Router.HandleFunc("/api/workout", h.AddWorkout).Methods("POST")
	h.Router.HandleFunc("/api/workout/{pKey}/{rKey}", h.UpdateWorkout).Methods("PUT")
	h.Router.HandleFunc("/api/workout/{pKey}/{rKey}", h.DeleteWorkout).Methods("DELETE")
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
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err := json.NewEncoder(w).Encode(Response{Message: message, Error: err.Error()}); err != nil {
		log.Error().Msg(err.Error())
	}
}
