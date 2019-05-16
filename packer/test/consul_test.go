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
	client, err := api.NewClient(&api.Config{
		Address: os.Getenv("CONSUL_ADDRESS"),
	})
	if err != nil {
		panic(err)
	}

	// Get a handle to the KV API
	kv := client.KV()

	suffix := rand.Int()
	key := fmt.Sprintf("TEST_CONSUL_KV_%d", suffix)

	// PUT a new KV pair
	p := &api.KVPair{Key: key, Value: []byte("1000")}
	_, err = kv.Put(p, nil)
	if err != nil {
		panic(err)
	}

	fmt.Printf("KV Put: %v %s\n", p.Key, p.Value)

	// Lookup the pair
	pair, _, err := kv.Get(key, nil)
	if err != nil {
		panic(err)
	}

	fmt.Printf("KV Get: %v %s\n", pair.Key, pair.Value)

	_, err = kv.Delete(p.Key, nil)

	if err != nil {
		panic(err)
	}

	fmt.Printf("KV Delete: %v\n", p.Key)
}
