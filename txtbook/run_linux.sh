export GOPATH=$(cd "$(dirname "$0")"; pwd)
export PATH=$PATH:$GOPATH/bin
go run main.go
