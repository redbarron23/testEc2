package test

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"math/rand"
	"os"
	"testing"
)

func TestConsulKVAvailable(t *testing.T) {
	// Get a new client
	// set CONSUL_ADDRESS=1.2.3.4:8500
	addr := os.Getenv("CONSUL_ADDRESS")
	if addr == "" {
		t.Fatal("Missing value for environment variable CONSUL_ADDRESS")
	}

	t.Logf("Using consul server at %s", addr)

	client, err := api.NewClient(&api.Config{
		Address: addr,
	})
	if err != nil {
		t.Fatalf("Error creating client %v", err)
	}

	// Get a handle to the KV API
	kv := client.KV()

	suffix := rand.Int()
	key := fmt.Sprintf("TEST_CONSUL_KV_%d", suffix)

	// PUT a new KV pair
	p := &api.KVPair{Key: key, Value: []byte(fmt.Sprintf("%d", rand.Int()))}
	_, err = kv.Put(p, nil)
	if err != nil {
		t.Fatalf("Error putting key %v", err)
	}

	t.Logf("KV Put: %v with data length %d\n", p.Key, len(p.Value))

	// Lookup the pair
	pair, _, err := kv.Get(key, nil)
	if err != nil {
		t.Fatalf("Error getting key %v", err)
	}

	t.Logf("KV Get: %v %s\n", pair.Key, pair.Value)

	_, err = kv.Delete(p.Key, nil)

	if err != nil {
		t.Fatalf("Error deleting key %v", err)
	}

	t.Logf("KV Delete: %v\n", p.Key)
}
