package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func createInstance(ip string, ami string, subnet string, sgid string) (MyinstanceID string) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-2")},
	)

	// Create EC2 service client
	svc := ec2.New(sess)

	// Specify the details of the instance that you want to create.
	runResult, err := svc.RunInstances(&ec2.RunInstancesInput{
		// An Amazon Linux AMI ID for t2.micro instances in the eu-west-2 region
		ImageId:          aws.String(ami), // ami-0188c0c5eddd2d032
		InstanceType:     aws.String("t2.micro"),
		MinCount:         aws.Int64(1),
		MaxCount:         aws.Int64(1),
		PrivateIpAddress: aws.String(ip),     // TBD: pull from terraform output
		SubnetId:         aws.String(subnet), // TBD: pull from terraform output
		KeyName:          aws.String(keyName),
		// SecurityGroups:   aws.StringSlice([]string{"terratest"}),
		SecurityGroupIds: []*string{
			aws.String(sgid),
		},
	})

	if err != nil {
		fmt.Println("Could not create instance", err)
		return
	}

	fmt.Println("Created instance", *runResult.Instances[0].InstanceId)

	MyinstanceID = *runResult.Instances[0].InstanceId

	instanceSlice := aws.StringSlice([]string{MyinstanceID})

	describeInstancesInput := &ec2.DescribeInstancesInput{
		InstanceIds: instanceSlice,
	}

	if err := svc.WaitUntilInstanceRunning(describeInstancesInput); err != nil {
		panic(err)
	}

	// Add tags to the created instance
	_, errtag := svc.CreateTags(&ec2.CreateTagsInput{
		Resources: []*string{runResult.Instances[0].InstanceId},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String("FT-Testing-E2E"),
			},
		},
	})

	if errtag != nil {
		log.Println("Could not create tags for instance", runResult.Instances[0].InstanceId, errtag)
		return
	}

	fmt.Println("Successfully tagged instance")

	return string(MyinstanceID)

}
