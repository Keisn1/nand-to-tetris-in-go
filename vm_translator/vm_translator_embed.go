package vmtrans

import (
	"embed"
)

//go:embed asm_codes/*.asm
var templateFiles embed.FS
