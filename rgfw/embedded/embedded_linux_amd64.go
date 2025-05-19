package embedded

import _ "embed"

const name = "libRGFW.so"

//go:embed linux_amd64/libRGFW.so
var lib []byte
