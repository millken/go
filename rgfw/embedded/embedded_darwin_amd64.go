package embedded

import _ "embed"

const name = "libRGFW.dylib"

//go:embed libRGFW_darwin_amd64.dylib
var lib []byte
