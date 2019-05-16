package test

import (
	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestDatadogRouting(t *testing.T) {
	t.Parallel()

	proxy := os.Getenv("http_proxy")

	url := "https://api.datadoghq.com/api/v1/dashboard"

	t.Logf("Testing proxy connection http_proxy='%s' to datadog %s", proxy, url)

	expected := "{\"errors\": [\"API key required\"]}"

	// Verify that we get back a 200 OK with the expected consul response
	status, resp := http_helper.HttpGet(t, url)

	if status != 403 {
		t.Fatalf("Datadog url %s returned unexpected status %d", url, status)
	}

	if resp != expected {
		t.Fatalf("Datadog url %s returned unexpected response %s", url, resp)
	}

	t.Logf("Datadog url %s returned the correct status '%d' and response: %s", url, status, resp)
}

func TestDatadogAgentStatus(t *testing.T) {
	t.Parallel()

	t.Log("Checking datadog agent's status")

	// TODO: determine if we need sudo or not
	cmd := exec.Command("/bin/sh", "-c", "sudo datadog-agent status")

	stdoutStderr, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("Command %s failed with error %v:\n %s",
			strings.Join(cmd.Args, " "), err, stdoutStderr)
	}

	t.Logf("Command %s finished with output: %s",
		strings.Join(cmd.Args, " "), stdoutStderr)

}
