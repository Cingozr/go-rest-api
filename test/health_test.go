//go:build e2e
// +build e2e

package test

import (
	"fmt"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func TestHealthEndPoint(t *testing.T) {
	fmt.Println("Running E2E test for health check endpoint")

	client := resty.New()

	rest, err := client.R().Get(BASE_URL + "/api/health")
	if err != nil {
		t.Fail()
	}

	assert.Equal(t, 200, rest.StatusCode())
}
