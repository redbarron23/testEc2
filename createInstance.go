package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// CreateInstance creates an ec2 instance for testing generated cloud infrastructure
func CreateInstance(svc *ec2.EC2, ami string, subnetId string, securityGroupId string) (instanceId string, instanceIp string) {
	// Specify the details of the instance that you want to create.
	runResult, err := svc.RunInstances(&ec2.RunInstancesInput{
		// An Amazon Linux AMI ID for t2.micro instances in the eu-west-2 region
		ImageId:          aws.String(ami), // ami-0188c0c5eddd2d032
		InstanceType:     aws.String("t2.micro"),
		MinCount:         aws.Int64(1),
		MaxCount:         aws.Int64(1),
		KeyName:          aws.String(keyName),
		SubnetId:         aws.String(subnetId),
		SecurityGroupIds: aws.StringSlice([]string{securityGroupId}),
	})

	if err != nil {
		log.Fatal("Could not create instance", err)
	}

	instanceId = *runResult.Instances[0].InstanceId

	fmt.Println("Created instance ", instanceId, "with ip", instanceIp)

	instanceSlice := aws.StringSlice([]string{instanceId})

	describeInstancesInput := &ec2.DescribeInstancesInput{
		InstanceIds: instanceSlice,
	}

	if err := svc.WaitUntilInstanceRunning(describeInstancesInput); err != nil {
		panic(err)
	}

	describeResult, err := svc.DescribeInstances(describeInstancesInput)

	if err != nil {
		panic(err)
	}

	// use PrivateIpAddress if possible?
	instanceIp = *describeResult.Reservations[0].Instances[0].PublicIpAddress

	// Add tags to the created instance
	_, errtag := svc.CreateTags(&ec2.CreateTagsInput{
		Resources: []*string{&instanceId},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String("GO-SDK-Instance"),
			},
		},
	})

	if errtag != nil {
		log.Println("Could not create tags for instance", instanceId, errtag)
		return
	}

	fmt.Println("Successfully tagged instance")

	return string(instanceId), instanceIp

}
