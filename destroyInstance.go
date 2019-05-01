package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func destroyInstance(ec2Inst string) {

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-2")},
	)

	// Create EC2 service client
	svc := ec2.New(sess)

	response, err := svc.TerminateInstances(&ec2.TerminateInstancesInput{
		InstanceIds: []*string{aws.String(ec2Inst)},
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(response)
}
