package test

import (
	"os"
	"testing"
	"time"

	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/stretchr/testify/assert"
)

func TestEndToEndHttp(t *testing.T) {
	t.Parallel()

	url := os.Getenv("HTTP")
	instanceText := "Consul Agent"
	expectedText := "Consul Agent"

	maxRetries := 3
	timeBetweenRetries := 5 * time.Second

	// Verify that we get back a 200 OK with the expected consul response
	http_helper.HttpGetWithRetry(t, url, 200, instanceText, maxRetries, timeBetweenRetries)

	assert.Equal(t, instanceText, expectedText, "Consul is running as expected")
}

// func TestReturnsNode(t *testing.T) {
// 	"http://"+ url + "/v1/catalog/nodes"
// }

func TestNoRedirect(t *testing.T) {
	t.Parallel()

	url := os.Getenv("HTTP")
	instanceText := "Consul Agent"
	expectedText := "Consul Agent"
	notExpectedText := "redirect" // It defintely shouldn't be Consul Agent o.k

	maxRetries := 3
	timeBetweenRetries := 5 * time.Second

	http_helper.HttpGetWithRetry(t, url, 200, instanceText, maxRetries, timeBetweenRetries)

	assert.NotEqual(t, notExpectedText, expectedText, "they should not be equal")
}
