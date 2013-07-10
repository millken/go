set GOPATH=%~dp0;
set path=%path%;$GOPATH/bin
go build  -ldflags "-s"
go install