#!/bin/sh
#go get -v
export GOPATH=$(cd "$(dirname "$0")"; pwd)
REVISION=`git rev-parse --short=5 HEAD`
go build mvc
go build logger
go install webtool
