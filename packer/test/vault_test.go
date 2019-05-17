package test

import (
	"errors"
	"fmt"
	"github.com/hashicorp/vault/api"
	"math/rand"
	"os"
	"testing"
)

// TestVaultAvailable connects to vault and reads/writes some values
// Env vars:
//   VAULT_ADDRESS - Address to connect to
// Auth:
//   VAULT_TOKEN - token to connect to vault
//    or
//   VAULT_USERNAME - username for login
//   VAULT_PASSWORD - password for login
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

	if err = configureAuthToken(client); err != nil {
		t.Fatalf("Could not configure auth token: %s", err)
	}

	suffix := rand.Int()
	key := fmt.Sprintf("cubbyhole/TEST_VAULT_SECRETS_%d", suffix)

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

func configureAuthToken(client *api.Client) error {

	token := os.Getenv("VAULT_TOKEN")
	if token != "" {
		client.SetToken(token)
		return nil
	}

	username := os.Getenv("VAULT_USERNAME")
	password := os.Getenv("VAULT_PASSWORD")

	if username == "" || password == "" {
		return errors.New("missing token or username and password")
	}

	// to pass the password
	options := map[string]interface{}{
		"password": password,
	}

	// the login path for userpass. Configure the correct method
	// https://www.vaultproject.io/docs/auth/userpass.html
	path := fmt.Sprintf("auth/userpass/login/%s", username)

	// PUT call to get a token
	secret, err := client.Logical().Write(path, options)
	if err != nil {
		return err
	}

	client.SetToken(secret.Auth.ClientToken)
	return nil
}
