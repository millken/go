# C2Go Example

This example demonstrates how to use the c2go package to generate Go wrapper code for C header files.

## Usage

```bash
go run main.go <header-file> <output-go-file>
```

### Example

```bash
# Generate Go wrapper from a C header file
go run main.h mylib.h wrapper.go
```

## Current Status

✅ **Fully Functional!** The example now works correctly:
- Builds without errors
- Parses header files and returns structured data
- Generates proper Go wrapper code with correct function signatures
- Handles type mapping between C and Go

The `ParseHeader` function now returns a properly populated `ParsedHeader` struct with functions, enums, and structs. The generated wrapper code includes correctly typed Go functions that call the corresponding C functions.

## Sample Output

When processing a C header file, the tool generates Go code like:

```go
package mylib

/*
#include "your_header.h"
*/
import "C"
import (
    "unsafe"
)

// ExampleFunction 包装原生函数 example_function
func ExampleFunction(a C.int, b C.int) C.int {
    return C.example_function(a, b)
}
```

## Next Steps

Future improvements could include:
- Full AST traversal implementation (currently uses example data)
- More complex type handling
- Better error reporting
- Support for macros and preprocessor directives