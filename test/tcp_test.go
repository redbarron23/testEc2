package test

import (
	"fmt"
	"net"
	"os"
	"testing"
)

func TestEndToEndSsh(t *testing.T) {
	t.Parallel()

	if len(os.Args) != 2 {
        	fmt.Printf("needs ENVIRONMENT VARIABLES set\n")
        	fmt.Printf("i.e. usage: IP=localhost PORT=22 go test -v tcp_test.go\n")
        	os.Exit(1)
    	}

	ip := os.Getenv("IP")
	port := os.Getenv("PORT")
	host := fmt.Sprintf("%s:%s", ip, port)

        conn, err := net.Dial("tcp", host)

	if err != nil {
		t.Errorf("server isn't responding on port: %v", err)
		os.Exit(1)

	}


	defer conn.Close()

}
