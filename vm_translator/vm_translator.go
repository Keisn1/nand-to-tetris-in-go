package vmtrans

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type CodeWriter struct {
	filename     string
	f            *os.File
	templates    map[string]*template.Template
	equalCounter int
}

func getFileName(fp string) string {
	filename := filepath.Base(fp)
	extension := filepath.Ext(filename)
	return filename[:len(filename)-len(extension)]
}

func NewCodeWriter(fp string) *CodeWriter {
	filename := getFileName(fp)
	file, _ := os.OpenFile(fp, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	cw := &CodeWriter{
		filename:  filename,
		f:         file,
		templates: make(map[string]*template.Template),
	}
	cw.templates = loadTemplates("asm_codes")
	return cw
}

func (cw *CodeWriter) handleArithmetic(cmdType, arg1 string) {
	var buf bytes.Buffer
	if arg1 == "neg" || arg1 == "not" {
		cw.templates[cmdType+" "+arg1].Execute(&buf, map[string]interface{}{
			"negation":          arg1,
			"negation_operator": negationOperators[arg1],
		})
		cw.f.Write(buf.Bytes())
		return
	}
	if arg1 == "eq" || arg1 == "lt" || arg1 == "gt" {
		cw.templates[cmdType+" "+arg1].Execute(&buf, map[string]interface{}{
			"comp":          arg1,
			"comp_operator": comparisonOperators[arg1],
			"comp_verbose":  comparisonVerbose[arg1],
			"counter":       cw.equalCounter,
		})
		cw.equalCounter++
		cw.f.Write(buf.Bytes())
		return
	} else {
		cw.templates[cmdType+" "+arg1].Execute(&buf, map[string]interface{}{
			"calculation":        arg1,
			"calculation_symbol": calculationSymbols[arg1],
		})
		cw.f.Write(buf.Bytes())
		return
	}
}

func (cw *CodeWriter) WriteArithmetic(cmdType, arg1, arg2 string) {
	var buf bytes.Buffer
	// TODO: refactor to switch, extract functions

	// switch cmdType {
	// case C_ARITHMETIC:
	// 	cw.handleArithmetic(cmdType, arg1)
	// }
	if cmdType == C_ARITHMETIC {
		if arg1 == "neg" || arg1 == "not" {
			cw.templates[cmdType+" "+arg1].Execute(&buf, map[string]interface{}{
				"negation":          arg1,
				"negation_operator": negationOperators[arg1],
			})
			cw.f.Write(buf.Bytes())
			return
		}
		if arg1 == "eq" || arg1 == "lt" || arg1 == "gt" {
			cw.templates[cmdType+" "+arg1].Execute(&buf, map[string]interface{}{
				"comp":          arg1,
				"comp_operator": comparisonOperators[arg1],
				"comp_verbose":  comparisonVerbose[arg1],
				"counter":       cw.equalCounter,
			})
			cw.equalCounter++
			cw.f.Write(buf.Bytes())
			return
		} else {
			cw.templates[cmdType+" "+arg1].Execute(&buf, map[string]interface{}{
				"calculation":        arg1,
				"calculation_symbol": calculationSymbols[arg1],
			})
			cw.f.Write(buf.Bytes())
			return

		}
	}

	if (cmdType == C_PUSH || cmdType == C_POP) && isGeneralSegment(arg1) {
		cw.templates[cmdType+" "+arg1].Execute(&buf, map[string]interface{}{
			"segment":               arg1,
			"segment_register_name": segmentRegisterName[arg1],
			"segment_register":      segmentRegisters[arg1],
			"x":                     arg2,
		})
		cw.f.Write(buf.Bytes())
		return
	}

	if (cmdType == C_PUSH || cmdType == C_POP) && arg1 == "static" {
		cw.templates[cmdType+" "+arg1].Execute(&buf, map[string]interface{}{
			"filename": strings.ToUpper(cw.filename),
			"x":        arg2,
		})
		cw.f.Write(buf.Bytes())
		return
	}

	if (cmdType == C_PUSH || cmdType == C_POP) && arg1 == "pointer" {
		segment := pointerToSegmentName[arg2]
		cw.templates[cmdType+" "+arg1].Execute(&buf, map[string]interface{}{
			"arg":                   arg2,
			"segment":               segment,
			"segment_register":      segmentRegisters[segment],
			"segment_register_name": segmentRegisterName[segment],
		})
		cw.f.Write(buf.Bytes())
		return
	}

	// temp is what is left
	cw.templates[cmdType+" "+arg1].Execute(&buf, map[string]interface{}{"x": arg2})
	cw.f.Write(buf.Bytes())
}

func (cw CodeWriter) WriteNewline() {
	cw.f.Write([]byte{'\n'})
}

func (cw CodeWriter) CloseFile() {
	cw.f.Close()
}
