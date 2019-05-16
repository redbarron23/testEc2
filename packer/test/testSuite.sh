#!/bin/bash

TESTPATH="/home/ec2-user/go/src/bitbucket/testEc2/test"

# we need this fail jenkins
set -e -u
rc=$?

# run all the tests
#==================
# run all tests
# go test

# run each test  need to put these in Jenkins Vars
REGION=eu-west-1 BUCKET_NAME=ctp-dev-lseg go test -v $TESTPATH/s3_test.go
HTTP=http://172.31.22.132:8500 go test -v $TESTPATH/http_test.go
IP=172.31.22.132 PORT=22 go test -v $TESTPATH/tcp_test.go
CONSUL_ADDRESS=1.2.3.4:8500 go test -v $TESTPATH/consul_test.go
VAULT_ADDRESS=https://1.2.3.4:8222 go test -v $TESTPATH/vault_test.go

# check err code
if [ $rc -ne 0 ]; then
  echo "testing failed" >&2
  exit $rc
fi
