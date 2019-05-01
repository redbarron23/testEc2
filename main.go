package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

// TBD
// key - local to jenkins

// These are static
const (
	TestInstance = "1.2.3.4" // change to private ip of newly instantiated host
	User         = "ec2-user"
	//GoVer  = "go1.12.4.linux-amd64.tar.gz"
)

func main() {

	// if len(os.Args) != 2 {
	// 	fmt.Println("testEc2Instance -ip '172.31.21.211' -ami 'ami-0188c0c5eddd2d032'")
	// }

	for i, arg := range os.Args {
		// print index and value
		fmt.Println("item", i, "is", arg)
	}

	// flags
	var IPAddr *string
	IPAddr = flag.String("ip", "", "IP Address")

	var Ami *string
	Ami = flag.String("ami", "", "AMI for your region")

	flag.Parse()

	// args := flag.Args()
	// if len(args) != 2 {
	// 	log.Fatal("Missing Flags -ip and -ami, `testEc2Instance -h for more information`")
	// }

	// create a new instance
	instance := createInstance(*IPAddr, *Ami)
	fmt.Println(instance)

	// need to add waiters
	time.Sleep(30 * time.Second)

	// Connect to Instance
	conn, err := Connect(TestInstance, User)
	if err != nil {
		log.Fatal(err)
	}

	output, err := conn.SendCommands("date",
		"sudo whoami",
		"ps -aux",
		//"sudo yum -y install wget unzip git",
		//"wget https://dl.google.com/go/go1.12.4.linux-amd64.tar.gz",
		//"sudo tar -C /usr/local -xvzf go1.12.4.linux-amd64.tar.gz",
		// "git clone "some bitbucket repo",
		// "GOHOME=/home/ec2-user/bit-bucket; cd ~/bit-bucket/test; go test -v"
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(output))

	defer destroyInstance(instance)

}
