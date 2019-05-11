#!/bin/bash

# we need this fail jenkins
set -e -u
rc=$?

# run all the tests
#==================
# run all tests
# go test

# run each test  need to put these in Jenkins Vars
REGION=eu-west-1 BUCKET_NAME=ctp-dev-lseg go test -v s3_test.go
HTTP=http://172.31.22.132:8500 go test -v http_test.go
IP=172.31.22.132 PORT=22 go test -v tcp_test.go

# check err code
if [ $rc -ne 0 ]; then
  echo "testing failed" >&2
  exit $rc
fi