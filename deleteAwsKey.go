package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func deleteAwsKey() {

	fmt.Println("Deleting Key....")

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-2")},
	)

	// Create an EC2 service client.
	svc := ec2.New(sess)

	input := &ec2.DeleteKeyPairInput{
		KeyName: aws.String(pairName),
	}

	result, err := svc.DeleteKeyPair(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(result)
}
