set GOPATH=%~dp0;
go get github.com/millken/go-whois
go build logger 
go build mvc 
go build banjo
go install webtool