package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// CreateAwsKey
func CreateAwsKey(svc *ec2.EC2, keyName string, pemFileName string) {
	// Creates a new key pair with the given name
	result, err := svc.CreateKeyPair(&ec2.CreateKeyPairInput{
		KeyName: aws.String(keyName),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == "InvalidKeyPair.Duplicate" {
			exitErrorf("Keypair %q already exists.", keyName)
		}
		exitErrorf("Unable to create key pair: %s, %v.", keyName, err)
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
	if err := w.Flush(); err != nil {
		log.Println("Could not flush", err)
	}

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
