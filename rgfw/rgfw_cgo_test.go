//go:build !purego
// +build !purego

package rgfw

import (
	"testing"
)

func TestSetClassName(t *testing.T) {
	SetClassName("Test Window")
	win := CreateWindow("RGFW Example Window", Rect{X: 500, Y: 500, W: 100, H: 100}, WindowCenter)
	if win == nil {
		t.Fatal("Failed to create window")
	}
	win.MakeCurrent()
}
