package test

import (
	"fmt"
	"net"
	"os"
	"testing"
)

func TestEndToEndSsh(t *testing.T) {
	t.Parallel()

	// need to make args
	IP := os.Getenv("IP")
	port := 22
	host := fmt.Sprintf("%s:%d", IP, port)

	// test ssh for consul
	conn, err := net.Dial("tcp", host)
	if err != nil {
		t.Errorf("myserver isn't responding on SSH port: %v", err)
	}

	conn.Close()

}
