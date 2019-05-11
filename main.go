package main

import (
	"flag"
	"fmt"
	"log"
	"time"
)

const (
	User    = "ec2-user"
	Key     = "test"
	KeyFile = "test.pem"
)

func main() {

	// flags
	var IPAddr *string
	IPAddr = flag.String("ip", "", "IP Address")

	var Ami *string
	Ami = flag.String("ami", "", "AMI for your region")

	flag.Parse()

	if flag.NFlag() != 2 {
		log.Fatal("Usage: testEc2Instance [-ip ipAddress] [-ami ami-abcdef123] -h for more info")
	}

	// create a new instance
	instance := createInstance(*IPAddr, *Ami)
	//TestInstance := *IPAddr // change to private ip of newly instantiated host
	fmt.Println(instance)
	fmt.Println("Starting Testing....")

	// need to add waiters
	time.Sleep(30 * time.Second)

	// Connect to Instance
	conn, err := Connect(*IPAddr, User)
	if err != nil {
		log.Fatal(err)
	}

	output, err := conn.SendCommands("date",
		"/usr/local/go/bin/go version",
		"chmod +x ./test/testSuite.sh",
		"./test/testSuite.sh",
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(output))

	defer destroyInstance(instance)

}
