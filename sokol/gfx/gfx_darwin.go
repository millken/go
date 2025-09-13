//go:build darwin

package gfx

// #cgo CFLAGS: -x objective-c
// #cgo LDFLAGS: -framework Foundation
import "C"
