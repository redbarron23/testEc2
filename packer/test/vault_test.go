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
	client, err := api.NewClient(&api.Config{
		Address: os.Getenv("VAULT_ADDRESS"),
	})

	if err != nil {
		panic(err)
	}

	suffix := rand.Int()
	key := fmt.Sprintf("TEST_VAULT_%d", suffix)

	// PUT a new KV pair
	secretData := map[string]interface{}{
		"value": key,
	}

	_, err = client.Logical().Write(key, secretData)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Vault Write: %s %v\n", key, secretData)

	// Lookup the pair
	secret, err := client.Logical().Read(key)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Vault Read: %s %v\n", key, secret.Data["value"])

	_, err = client.Logical().Delete(key)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Vault Delete: %v\n", key)
}
