package test

import (
	"os"
	"testing"
	"time"

	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
)

func TestEndToEndHttp(t *testing.T) {
	t.Parallel()

	url := os.Getenv("HTTP")
	instanceText := "Consul Agent"

	maxRetries := 3
	timeBetweenRetries := 5 * time.Second

	// Verify that we get back a 200 OK with the expected consul response
	http_helper.HttpGetWithRetry(t, url, 200, instanceText, maxRetries, timeBetweenRetries)

}
