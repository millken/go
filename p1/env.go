//https://coderwall.com/p/kjuyqw
/*
export KEY=abc=cdef;KEY2=tester;
go build env.go
./env
unset KEY2
echo $KEY2
*/
package main

import (
	"fmt"
	"os"
	"strings"
	"runtime"
)

var _ENV map[string]string

func init() {
	getenvironment := func(data []string, getkeyval func(item string) (key, val string)) map[string]string {
		items := make(map[string]string)
		for _, item := range data {
			key, val := getkeyval(item)
			items[key] = val
		}
		return items
	}
	_ENV = getenvironment(os.Environ(), func(item string) (key, val string) {
		splits := strings.Split(item, "=")
		key = splits[0]
		val = strings.Join(splits[1:], "=")
		return
	})
}

func main() {

	fmt.Println(_ENV["KEY"])
	fmt.Println(_ENV["KEY2"])
	fmt.Println(runtime.GOOS)
}
