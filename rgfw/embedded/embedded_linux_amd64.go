package embedded

import _ "embed"

const name = "libRGFW.so"

//go:embed libRGFW_linux_amd64.so
var lib []byte
