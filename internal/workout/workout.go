package workout

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/rs/zerolog/log"
)

// tableName
var tableName = os.Getenv("AZ_TABLE_NAME")

// Service : workout service
type Service struct {
	DB *aztables.ServiceClient
}

// Workout : struct
type Workout struct {
	PartitionKey string               `json:"PartitionKey"`
	RowKey       string               `json:"RowKey"`
	Date         aztables.EDMDateTime `json:"Date"` // 2020-01-26T00:00:00.000Z
	Weight       int32                `json:"Weight"`
	Reps         int32                `json:"Reps"`
	Sets         int32                `json:"Sets"`
	Warmup       bool                 `json:"Warmup"`
}

// WorkoutService : interface for workout service
type WorkoutService interface {
	GetWorkout(pKey string, rKey string) (Workout, error)
	AddWorkout(workout Workout) (int, error)
	UpdateWorkout(pkey string, rKey string, newWorkout Workout) (Workout, error)
	DeleteWorkout(pkey string, rKey string) error
	GetAllWorkouts() ([]Workout, error)
}

// NewService : new workout service
func NewService(db *aztables.ServiceClient) *Service {
	return &Service{
		DB: db,
	}
}

// GetWorkout : get workout
func (s *Service) GetWorkout(pKey string, rKey string) (Workout, error) {
	var workout Workout

	client := s.DB.NewClient(tableName)

	filter := fmt.Sprintf("PartitionKey eq '%s' and RowKey eq '%s'", pKey, rKey)
	options := &aztables.ListEntitiesOptions{
		Filter: &filter,
		Select: to.StringPtr("PartitionKey,RowKey,Date,Weight,Reps,Sets,Warmup"),
		// Top:    to.Int32Ptr(15),
	}

	pager := client.List(options) // pass in "nil" if you want to list all entities
	for pager.NextPage(context.TODO()) {
		resp := pager.PageResponse()
		log.Info().Msgf("I found: %v entities in this table.\n", len(resp.Entities))

		for _, entity := range resp.Entities {
			var myEntity Workout
			err := json.Unmarshal(entity, &myEntity)

			if err != nil {
				log.Error().Msg(err.Error())
			}
			workout = myEntity
		}
	}
	return workout, nil
}

// AddWorkout : add new workout
func (s *Service) AddWorkout(workout Workout) (Workout, error) {
	client := s.DB.NewClient(tableName)
	// check if entity exists first
	pKey := workout.PartitionKey
	rKey := workout.RowKey

	if _, err := s.GetWorkout(pKey, rKey); err != nil {
		log.Error().Msg(err.Error())
		return Workout{}, err
	}

	//
	myEntity := aztables.EDMEntity{
		Entity: aztables.Entity{
			PartitionKey: workout.PartitionKey,
			RowKey:       workout.RowKey,
		},
		Properties: map[string]interface{}{
			"Date":   workout.Date,
			"Weight": workout.Weight,
			"Reps":   workout.Reps,
			"Sets":   workout.Sets,
			"Warmup": workout.Warmup,
		},
	}

	marshalled, err := json.Marshal(myEntity)
	if err != nil {
		log.Error().Msg(err.Error())
	}

	resp, err := client.AddEntity(context.TODO(), marshalled, nil)
	if err != nil {
		log.Error().Msg(err.Error())
	}
	log.Info().Msgf("RawResponse Status: %+v", resp.RawResponse.Status)
	return workout, nil
}

// UpdateWorkout : update workout
func (s *Service) UpdateWorkout(pKey string, rKey string, newWorkout Workout) (Workout, error) {
	client := s.DB.NewClient(tableName)

	// check if entity exists first
	if _, err := s.GetWorkout(pKey, rKey); err != nil {
		log.Error().Msg(err.Error())
		return Workout{}, err
	}

	//
	myEntity := aztables.EDMEntity{
		Entity: aztables.Entity{
			PartitionKey: newWorkout.PartitionKey,
			RowKey:       newWorkout.RowKey,
		},
		Properties: map[string]interface{}{
			"Date":   newWorkout.Date,
			"Weight": newWorkout.Weight,
			"Reps":   newWorkout.Reps,
			"Sets":   newWorkout.Sets,
			"Warmup": newWorkout.Warmup,
		},
	}

	marshalled, err := json.Marshal(myEntity)
	if err != nil {
		log.Error().Msg(err.Error())
	}
	resp, err := client.UpdateEntity(context.TODO(), marshalled, nil)
	if err != nil {
		log.Error().Msg(err.Error())
	}
	log.Info().Msgf("RawResponse Status: %+v", resp.RawResponse.Status)
	return newWorkout, nil
}

// DeleteWorkout : delete workout
func (s *Service) DeleteWorkout(pKey string, rKey string) error {
	client := s.DB.NewClient(tableName)
	// check if entity exists first
	if _, err := s.GetWorkout(pKey, rKey); err != nil {
		log.Error().Msg(err.Error())
		return err
	}

	// entity does not exist; ok to delete
	resp, err := client.DeleteEntity(context.TODO(), pKey, rKey, nil)
	if err != nil {
		log.Error().Msg(err.Error())
	}
	log.Info().Msgf("RawResponse Status: %+v", resp.RawResponse.Status)

	return err
}

// GetAllWorkouts : get all workouts
func (s *Service) GetAllWorkouts() ([]Workout, error) {
	var workouts []Workout
	client := s.DB.NewClient(tableName)

	pager := client.List(nil) // pass in "nil" if you want to list all entities
	for pager.NextPage(context.TODO()) {
		resp := pager.PageResponse()
		log.Info().Msgf("I found: %v entities in this table.\n", len(resp.Entities))

		for _, entity := range resp.Entities {
			var myEntity Workout
			err := json.Unmarshal(entity, &myEntity)

			if err != nil {
				log.Error().Msg(err.Error())
			}
			workouts = append(workouts, myEntity)
		}
	}
	return workouts, nil
}
