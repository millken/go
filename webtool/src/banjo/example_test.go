package banjo_test

import (
        "github.com/bmatsuo/banjo"
        "os"
)

// This example is initialized by parsing a view named "hello". The view is
// rendered by a context that has been assigned a value for "Name".
func Example() {
        banjo.Parse("hello", "Hello, {{.Name}}!\n")

        context := banjo.NewContext()
        context.Set("Name", "Banjo")
        context.Render(os.Stdout, "hello")

        // Output:
        // Hello, Banjo!
}
