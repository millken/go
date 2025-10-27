package main

import (
	"fmt"
	"os"

	"github.com/millken/go/c2go"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: c2go <header-file> <output-go-file>")
		os.Exit(1)
	}

	headerFile := os.Args[1]
	outputFile := os.Args[2]

	// 解析头文件
	parsedHeader, err := c2go.ParseHeader(headerFile)
	if err != nil {
		fmt.Printf("Error parsing header: %v\n", err)
		os.Exit(1)
	}

	// 生成Go包装代码
	pkgName := "mylib"
	code, err := c2go.GenerateWrapper(pkgName, headerFile, parsedHeader.Functions)
	if err != nil {
		fmt.Printf("Error generating wrapper: %v\n", err)
		os.Exit(1)
	}

	// 写入文件
	if err := os.WriteFile(outputFile, []byte(code), 0644); err != nil {
		fmt.Printf("Error writing output: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully generated %s\n", outputFile)
}
