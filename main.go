package main

import (
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"log"
	"os"
	"testing"
)

const (
	User        = "ec2-user"
	keyName     = "test"
	pemFileName = "test.pem"
)

var cwd, _ = os.Getwd()

// flags
var ami = flag.String("ami", "", "AMI for your region")
var region = flag.String("region", os.Getenv("AWS_REGION"), "AWS Region. Default env AWS_REGION")
var terraformState = flag.String("state-folder", cwd, "Folder path containing the terraform.tfstate file. Default CWD")
var terraformOptions *terraform.Options

// needed to pass into terratest functions. It is used for logging on errors
var tStub = &testing.T{}

func main() {

	flag.Parse()

	if *ami == "" {
		log.Fatal("Usage: testEc2Instance -ami ami-abcdef123 [-region ipAddress] [-state-folder /a/path] -h for more info")
	}

	terraformOptions = &terraform.Options{
		TerraformDir: *terraformState,
	}

	terraform.Init(tStub, terraformOptions)

	// Initialize a session in eu-west-2 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(*region)},
	)

	// Create an EC2 service client.
	svc := ec2.New(sess)

	// create key
	CreateAwsKey(svc, keyName, pemFileName)
	defer DeleteAwsKey(svc)

	securityGroupId := terraform.OutputRequired(tStub, terraformOptions, "instance_sg_id")
	subnetId := terraform.OutputRequired(tStub, terraformOptions, "instance_subnet_id")

	// create a new instance
	instance, ip := CreateInstance(svc, *ami, subnetId, securityGroupId)
	defer DestroyInstance(svc, instance)
	//TestInstance := *IPAddr // change to private ip of newly instantiated host
	fmt.Println(instance)
	fmt.Println("Starting Testing....")

	// Connect to Instance
	conn, err := Connect(ip, User, pemFileName)
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

}
