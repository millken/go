package embedded

import _ "embed"

const name = "RGFW.dll"

//go:embed libRGFW_windows_amd64.dll
var lib []byte
