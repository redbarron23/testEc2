package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func DestroyInstance(svc *ec2.EC2, ec2Inst string) {

	response, err := svc.TerminateInstances(&ec2.TerminateInstancesInput{
		InstanceIds: []*string{aws.String(ec2Inst)},
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(response)
}
