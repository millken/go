package picasso

/*
#cgo CFLAGS: -I${SRCDIR}
#cgo darwin CXXFLAGS:  -std=c++11 -I${SRCDIR} -mmacosx-version-min=15.4
#cgo LDFLAGS:  -lpicasso2_sw -lstdc++
#cgo darwin LDFLAGS: -mmacosx-version-min=15.4 -L${SRCDIR}/_libs/darwin/arm64  -framework CoreGraphics -framework CoreText
#include "include/picasso.h"
#include <stdlib.h>
#include <stddef.h>
#include <stdint.h>
#include "cgo_helpers.h"

*/
import "C"

const (
	True  Bool = 1
	False Bool = 0
)
