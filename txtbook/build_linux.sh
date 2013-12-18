#!/bin/sh
#go get -v
export GOPATH=$(cd "$(dirname "$0")"; pwd)
REVISION=`git rev-parse --short=5 HEAD`
echo $REVISION > REVISION
go build -ldflags "-s -X main.gitVersion $REVISION" -v
go install