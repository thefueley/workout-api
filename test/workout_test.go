//go:build e2e
// +build e2e

package test

import (
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"testing"
)

var addWOBody = `{"PartitionKey":"Fake","RowKey":"1999","Date":"2020-01-26T00:00:00.000Z", "Weight": 185, "Reps": 3, "Sets": 1, "Warmup": true}`
var updateWOBody = `{"PartitionKey":"Fake","RowKey":"1999","Date":"2020-01-26T00:00:00.000Z", "Weight": 195, "Reps": 3, "Sets": 1, "Warmup": true}`

func TestAddWorkout(t *testing.T) {
	client := resty.New()
	resp, err := client.R().
		SetBody(addWOBody).
		Post(BASE_URL + "/api/workout")

	assert.NoError(t, err)

	assert.Equal(t, 200, resp.StatusCode())
}

func TestGetWorkout(t *testing.T) {
	client := resty.New()
	resp, err := client.R().Get(BASE_URL + "/api/workout/Fake/1999")

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
		SetBody(updateWOBody).
		Put(BASE_URL + "/api/workout/Fake/1999")

	assert.NoError(t, err)

	assert.Equal(t, 200, resp.StatusCode())
}

func TestDeleteWorkout(t *testing.T) {
	client := resty.New()
	resp, err := client.R().Delete(BASE_URL + "/api/workout/Fake/1999")

	assert.NoError(t, err)

	assert.Equal(t, 200, resp.StatusCode())
}
