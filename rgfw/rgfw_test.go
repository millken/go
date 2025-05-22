//go:build !cgo

package rgfw

import (
	"testing"
)

func TestXX(t *testing.T) {
	x := Rect{X: 500, Y: 500, W: 100, H: 100}
	CreateWindow("RGFW Example Window", x, WindowCenter)
}
