package vmtrans

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"
)

type CodeWriter struct {
	f         *os.File
	templates map[string]*template.Template
}

func NewCodeWriter(fp string) *CodeWriter {
	file, _ := os.OpenFile(fp, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	cw := &CodeWriter{f: file, templates: make(map[string]*template.Template)}
	cw.templates = loadTemplates("asm_codes")
	fmt.Println(cw.templates)
	return cw
}

func (cw CodeWriter) CloseFile() {
	cw.f.Close()
}

func (cw *CodeWriter) WriteArithmetic(cmdType, arg1, arg2 string) {
	var buf bytes.Buffer

	if cmdType == C_ARITHMETIC {
		cw.templates[cmdType+" "+arg1].Execute(&buf, map[string]interface{}{})
		cw.f.Write(buf.Bytes())
		return
	}

	if cmdType == C_PUSH && isGeneralSegment(arg1) {
		cw.templates[cmdType+" "+arg1].Execute(&buf, map[string]interface{}{
			"segment":               arg1,
			"segment_register_name": segmentRegisterName[arg1],
			"segment_register":      segmentRegisters[arg1],
			"x":                     strings.TrimSpace(arg2),
		})
		cw.f.Write(buf.Bytes())
		return
	}

	cw.templates[cmdType+" "+arg1].Execute(&buf, map[string]interface{}{"x": strings.TrimSpace(arg2)})
	cw.f.Write(buf.Bytes())
}
