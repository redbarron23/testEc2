package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// createAwsKey
func createAwsKey() {
	fmt.Println("Creating Key....")

	//pairName := "test2"
	pemFileName := pairName + ".pem"

	// Initialize a session in eu-west-2 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-2")},
	)

	// Create an EC2 service client.
	svc := ec2.New(sess)

	// Creates a new key pair with the given name
	result, err := svc.CreateKeyPair(&ec2.CreateKeyPairInput{
		KeyName: aws.String(pairName),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == "InvalidKeyPair.Duplicate" {
			exitErrorf("Keypair %q already exists.", pairName)
		}
		exitErrorf("Unable to create key pair: %s, %v.", pairName, err)
	}

	// lets write to a pem file and STDOUT
	f, err := os.Create(pemFileName)
	check(err)
	defer f.Close()

	fmt.Printf("Created key pair %q %s\n",
		*result.KeyName, *result.KeyFingerprint)

	w := bufio.NewWriter(f)
	_, err = fmt.Fprintf(w, "%v\n", *result.KeyMaterial)

	check(err)
	w.Flush()

	err = os.Chmod(pemFileName, 0600)
	if err != nil {
		log.Println(err)
	}

}

// check
func check(err error) {
	if err != nil {
		panic(err)
	}
}

// exitErrorf
func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
