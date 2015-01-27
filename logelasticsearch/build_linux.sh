#!/bin/sh
current_dir=$(cd "$(dirname "$0")"; pwd)
export GOPATH=$current_dir
export GOBIN=$GOPATH/bin
REVISION=`git rev-parse --short=5 HEAD`
go get -v
go fmt
go install


