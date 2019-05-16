package test

import (
	"fmt"
	"github.com/hashicorp/vault/api"
	"math/rand"
	"os"
	"testing"
)

func TestVaultAvailable(t *testing.T) {
	// Get a new client
	// set VAULT_ADDRESS=https://1.2.3.4:8222
	addr := os.Getenv("VAULT_ADDRESS")
	if addr == "" {
		t.Fatal("Missing value for environment variable VAULT_ADDRESS")
	}

	t.Logf("Using vault server at %s", addr)

	client, err := api.NewClient(&api.Config{
		Address: addr,
	})

	if err != nil {
		t.Fatalf("Error creating client %v", err)
	}

	suffix := rand.Int()
	key := fmt.Sprintf("TEST_VAULT_%d", suffix)

	// PUT a new KV pair
	secretData := map[string]interface{}{
		"value": key,
	}

	_, err = client.Logical().Write(key, secretData)
	if err != nil {
		t.Fatalf("Error writing key %v", err)
	}

	t.Logf("Vault Write: %s %v\n", key, secretData)

	// Lookup the pair
	secret, err := client.Logical().Read(key)
	if err != nil {
		t.Fatalf("Error reading key %v", err)
	}

	t.Logf("Vault Read: %s %v\n", key, secret.Data["value"])

	_, err = client.Logical().Delete(key)

	if err != nil {
		t.Fatalf("Error deleting key %v", err)
	}

	t.Logf("Vault Delete: %v\n", key)
}
