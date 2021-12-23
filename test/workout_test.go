//go:build e2e
// +build e2e

package test

import (
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddWorkout(t *testing.T) {
	client := resty.New()
	resp, err := client.R().
		SetBody(`{"id":1,"date":"unknown","exercises":[{"id":1,"name":"n/a","weight":1,"sets":1,"reps":1,"warmup":true},{"id":2,"name":"n/a","weight":1,"sets":1,"reps":1,"warmup":true},{"id":3,"name":"n/a","weight":1,"sets":1,"reps":1,"warmup":true},{"id":4,"name":"n/a","weight":1,"sets":1,"reps":1,"warmup":false},{"id":5,"name":"n/a","weight":1,"sets":1,"reps":1,"warmup":false},{"id":6,"name":"n/a","weight":1,"sets":1,"reps":1,"warmup":false},{"id":7,"name":"n/a","weight":1,"sets":1,"reps":1,"warmup":false}]}`).
		Post(BASE_URL + "/api/workout")

	assert.NoError(t, err)

	assert.Equal(t, 200, resp.StatusCode())
}

func TestGetWorkout(t *testing.T) {
	client := resty.New()
	resp, err := client.R().Get(BASE_URL + "/api/workout/1")

	if err != nil {
		t.Fail()
	}

	assert.Equal(t, 200, resp.StatusCode())
}

func TestGetAllWorkouts(t *testing.T) {
	client := resty.New()
	resp, err := client.R().Get(BASE_URL + "/api/workout")

	if err != nil {
		t.Fail()
	}

	assert.Equal(t, 200, resp.StatusCode())
}

func TestUpdateWorkout(t *testing.T) {
	client := resty.New()
	resp, err := client.R().
		SetBody(`{"id":1,"date":"unknown","exercises":[{"id":1,"name":"n/a","weight":111,"sets":111,"reps":111,"warmup":true},{"id":2,"name":"n/a","weight":111,"sets":111,"reps":111,"warmup":true},{"id":3,"name":"n/a","weight":111,"sets":111,"reps":111,"warmup":true},{"id":4,"name":"n/a","weight":111,"sets":111,"reps":111,"warmup":true},{"id":5,"name":"n/a","weight":111,"sets":111,"reps":111,"warmup":true},{"id":6,"name":"n/a","weight":111,"sets":111,"reps":111,"warmup":true},{"id":7,"name":"n/a","weight":111,"sets":111,"reps":111,"warmup":true}]}`).
		Put(BASE_URL + "/api/workout/1")

	assert.NoError(t, err)

	assert.Equal(t, 200, resp.StatusCode())
}

func TestDeleteWorkout(t *testing.T) {
	client := resty.New()
	resp, err := client.R().Delete(BASE_URL + "/api/workout/1")

	assert.NoError(t, err)

	assert.Equal(t, 200, resp.StatusCode())
}
