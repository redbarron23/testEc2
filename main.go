package main

import (
	"flag"
	"fmt"
	"log"
	"time"
)

// User for ssh
const (
	User = "ec2-user"
)

var keyName string
var pemFileName string

func init() {
	keyName = InitKey()
	pemFileName = InitKeyName()
}

// InitKey persist key throughout
func InitKey() string {
	return "test2"
}

// InitKeyName persist keyname throughout
func InitKeyName() string {
	return "test2.pem"
}

func main() {

	// flags
	var IPAddr *string
	IPAddr = flag.String("ip", "", "IP Address")

	var Ami *string
	Ami = flag.String("ami", "", "AMI for your region")

	var Subnet *string
	Subnet = flag.String("subnet", "", "Subnet for your VPC")

	var SgId *string
	Subnet = flag.String("sgid", "", "Security Group for your VPC")

	flag.Parse()

	if flag.NFlag() != 4 {
		log.Fatal("Usage: testEc2Instance [-ip ipAddress] [-ami ami-abcdef123] [-subnet subnet-6e7f829e] [-sgid sg-0b44d9223c8071f5b] -h for more info")
	}

	// create key
	createAwsKey()

	// create a new instance
	instance := createInstance(*IPAddr, *Ami, *Subnet, *SgId)
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
		"cp /home/ec2-user/testSuite.sh /home/ec2-user/go/src/bitbucket/testEc2/test/",
		"chmod +x /home/ec2-user/go/src/bitbucket/testEc2/test/testSuite.sh",
		"/home/ec2-user/go/src/bitbucket/testEc2/test/testSuite.sh",
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(output))

	deleteAwsKey()

	defer destroyInstance(instance)

}
