# Synopsis
Launches Ec2 Instance given parametres of subnet and ami and runs a series of tests to test Shared Services.
This is intended to be compiled from source

** Disclaimer
runs only for RedHAT Linux at the moment

## Linux rpms

```sudo yum -y install git wget gcc```

## Go Requirements

## Latest version of Go

Run  ```go version``` to check version

currently tested go1.11.6

set GOPATH
i.e.  ```export GOPATH=$HOME/go```

Install GO dep
https://github.com/golang/dep

```dep version```
to verify
currently tested with 0.5.1

## AWS Requirements

- AWS Credentials and Config

- ssh key

- Permissions to Launch Ec2

## Dependencies

- setting up GOHOME

```mkdir -p ~/go```

```export GOPATH=~/go```

```export GOBIN=$GOPATH/bin```

```export PATH=$PATH:/usr/local/go/bin:$GOBIN```

```mkdir -p ~/go/src ~/go/bin```

- godep

```curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh```

```cd $GOHOME && mkdir -p bitbucket/testEc2/test```

```dep init```

```dep ensure --add github.com/aws/aws-sdk-go```

```dep ensure -add github.com/gruntwork-io/terratest/modules/aws```

## Run without Building

```GOCACHE=off go run *.go```

## Building

```go build -o testEc2Instance```

## Running

```./testEc2Instance -ip "172.31.21.211" -ami "ami-0188c0c5eddd2d032"```

## Running Tests Individually outside of Targeted ec2

```cd ./test```

```HTTP=http://10.10.10.10:8500 go test -v http_test.go```

```IP=10.10.10.10 PORT=22 go test -v tcp_test.go```

```BUCKET_NAME=urbucketname go test -v s3_test.go```

## Potential Errata

- Wrong subnet for VPC

    Could not create instance InvalidParameterValue: Address 172.31.32.2 does not fall within the subnet's address range
    status code: 400, request id: 6c739663-3a0d-43ff-b9cb-356c9484df36

- Security group is not created or defined

    Could not create instance InvalidParameterValue: Value () for parameter groupId is invalid. The value cannot be empty
    status code: 400, request id: 4aaf3d16-c366-467a-b115-bef66a855106

- IP is already in use

    Could not create instance InvalidIPAddress.InUse: Address 172.31.22.128 is in use.
    status code: 400, request id: 4957ed3b-2bfa-434b-97de-01e843c11bf5

## Misc

