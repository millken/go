//go:build darwin && SOKOL_METAL

package gfx

/*
#cgo LDFLAGS: -framework Metal

#define SOKOL_GFX_IMPL
#define SOKOL_METAL
#include "sokol_gfx.h"
*/
import "C"
