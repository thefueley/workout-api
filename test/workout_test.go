//go:build e2e
// +build e2e

package test

import (
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"testing"
)

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

func TestAddWorkout(t *testing.T) {
	client := resty.New()
	resp, err := client.R().
		SetBody(`{"date": "21 DEC 2021", "duration": "1 hr 0 mins", "exercise": [{"name": "deadlift", "weight": 175, "set": 1, "rep": 5}, {"name": "deadlift", "weight": 215, "set": 1, "rep": 5}, {"name": "deadlift", "weight": 235, "set": 1, "rep": 3}]}`).
		Post(BASE_URL + "/api/workout")

	assert.NoError(t, err)

	assert.Equal(t, 200, resp.StatusCode())
}

func TestUpdateWorkout(t *testing.T) {
	client := resty.New()
	resp, err := client.R().
		SetBody(`{"date": "32 DEC 1999", "duration": "0 hr 0 mins", "exercise": [{"id": 1, "name": "curl", "weight": 1, "set": 1, "rep": 1}, {"id": 2,"name": "curl", "weight": 1, "set": 1, "rep": 1}, {"id": 3,"name": "curl", "weight": 1, "set": 1, "rep": 1}]}`).
		Put(BASE_URL + "/api/workout/1")

	assert.NoError(t, err)

	assert.Equal(t, 200, resp.StatusCode())
}

func TestDeleteWorkout(t *testing.T) {
	client := resty.New()
	resp, err := client.R().Delete(BASE_URL + "/api/workout/6")

	assert.NoError(t, err)

	assert.Equal(t, 200, resp.StatusCode())
}
