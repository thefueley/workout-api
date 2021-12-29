package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/thefueley/workout-api/models"
)

var hmacSigningKey = []byte(os.Getenv("WORKOUT_API_TOKEN"))

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

// createToken : create JWT
func CreateToken() {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"PartitionKey": "Squat",
		"RowKey":       "1313",
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(hmacSigningKey)

	log.Info().Msgf("tokenString: %s Error: %v", tokenString, err)
}

// validateToken : validates an incoming jwt token
func validateToken(accessToken string) bool {
	//myKey := []byte(hmacSigningKey)
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return hmacSigningKey, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		log.Info().Msgf("claims[%v], claims[%v]", claims["PartitionKey"], claims["RowKey"])
	} else {
		log.Error().Msg(err.Error())
	}
	return token.Valid
}

// JWTAuth - a handy middleware function that will provide basic auth around specific endpoints
func JWTAuth(original func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info().Msg("jwt auth endpoint hit")
		authHeader := r.Header["Authorization"]

		if authHeader == nil {
			log.Error().Msg("authHeader is nil")
			sendErrorResponse(w, "not authorized", errors.New("not authorized"))
			return
		}

		authHeaderParts := strings.Split(authHeader[0], " ")

		if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
			log.Error().Msg("authHeaderParts != 2 or does not contain bearer")
			sendErrorResponse(w, "not authorized", errors.New("not authorized"))
			return
		}

		if validateToken(authHeaderParts[1]) {
			original(w, r)
		} else {
			log.Error().Msg("Token is invalid")
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

	h.Router.HandleFunc("/api/workout/{pKey}/{rKey}", JWTAuth(h.GetWorkout)).Methods("GET")
	h.Router.HandleFunc("/api/workout", JWTAuth(h.AddWorkout)).Methods("POST")
	h.Router.HandleFunc("/api/workout/{pKey}/{rKey}", JWTAuth(h.UpdateWorkout)).Methods("PUT")
	h.Router.HandleFunc("/api/workout/{pKey}/{rKey}", JWTAuth(h.DeleteWorkout)).Methods("DELETE")
	h.Router.HandleFunc("/api/workout", JWTAuth(h.GetAllWorkouts)).Methods("GET")

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
