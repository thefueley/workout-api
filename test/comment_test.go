//go:build e2e
// +build e2e

package test

import (
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetAllComments(t *testing.T) {
	client := resty.New()
	resp, err := client.R().Get(BASE_URL + "/api/comment")

	if err != nil {
		t.Fail()
	}

	assert.Equal(t, 200, resp.StatusCode())
}

func TestPostComment(t *testing.T) {
	client := resty.New()
	resp, err := client.R().
		SetBody(`{"slug": "/e2e", "author": "666", "body": "e2e test post"}`).
		Post(BASE_URL + "/api/comment")

	assert.NoError(t, err)

	assert.Equal(t, 200, resp.StatusCode())
}
