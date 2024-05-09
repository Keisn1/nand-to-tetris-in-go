package vmtrans

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
)

type CodeWriter struct {
	src         string
	out         *os.File
	templates   map[string]*template.Template
	compCounter int
	funcCounter int
}

func getFileName(fp string) string {
	if len(fp) == 0 {
		return fp
	}
	filename := filepath.Base(fp)
	extension := filepath.Ext(filename)
	filename = filename[:len(filename)-len(extension)]
	return strings.ToUpper(string(filename[0])) + filename[1:]
}

func NewCodeWriter(outPath, srcPath string) *CodeWriter {
	src := getFileName(srcPath)
	file, err := os.OpenFile(outPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	cw := &CodeWriter{
		src:       src,
		out:       file,
		templates: make(map[string]*template.Template),
	}
	cw.templates = loadTemplates("asm_codes")
	return cw
}

func (cw *CodeWriter) Write(cmdType, arg1, arg2 string) {
	var buf bytes.Buffer

	switch cmdType {
	case C_ARITHMETIC:
		cw.handleArithmeticCmd(cmdType, arg1, buf)
	case C_POP:
		cw.handleSegmentCmd(cmdType, arg1, arg2, buf)
	case C_PUSH:
		cw.handleSegmentCmd(cmdType, arg1, arg2, buf)
	case C_LABEL:
		cw.handleLabelCmd(cmdType, arg1, buf)
	case C_IF:
		cw.handleGotoCmd(cmdType, arg1, buf)
	case C_GOTO:
		cw.handleGotoCmd(cmdType, arg1, buf)
	case C_FUNCTION:
		cw.handleFunctionCmd(cmdType, arg1, arg2, buf)
	case C_RETURN:
		cw.handleReturn(cmdType, buf)
	case C_CALL:
		cw.handleCall(cmdType, arg1, arg2, buf)
	}
}

func (cw *CodeWriter) WriteBootStrap() {
	var buf bytes.Buffer
	cw.templates["boot"].Execute(&buf, map[string]interface{}{})
	cw.out.Write(buf.Bytes())
}

func (cw *CodeWriter) handleReturn(cmdType string, buf bytes.Buffer) {
	cw.templates[cmdType].Execute(&buf, map[string]interface{}{})
	cw.out.Write(buf.Bytes())
}

func (cw *CodeWriter) handleFunctionCmd(cmdType, arg1, arg2 string, buf bytes.Buffer) {
	n, err := strconv.Atoi(arg2)
	if err != nil {
		panic(err)
	}

	cw.templates[cmdType].Execute(&buf, map[string]interface{}{
		"functionName":  arg1,
		"nbrIterations": arg2,
		"Numbers":       make([]struct{}, n),
	})
	cw.out.Write(buf.Bytes())
}

func (cw *CodeWriter) handleGotoCmd(cmdType, arg1 string, buf bytes.Buffer) {
	cw.templates[cmdType].Execute(&buf, map[string]interface{}{
		"loopName": arg1,
	})
	cw.out.Write(buf.Bytes())
}

func (cw *CodeWriter) handleLabelCmd(cmdType, arg1 string, buf bytes.Buffer) {
	cw.templates[cmdType].Execute(&buf, map[string]interface{}{
		"label": arg1,
	})
	cw.out.Write(buf.Bytes())
}

func (cw *CodeWriter) handleCall(cmdType, arg1, arg2 string, buf bytes.Buffer) {
	cw.templates[cmdType].Execute(&buf, map[string]interface{}{
		"caller":           cw.src,
		"function_name":    arg1,
		"n_args":           arg2,
		"function_counter": cw.funcCounter,
	})
	cw.funcCounter++
	cw.out.Write(buf.Bytes())
}

func (cw *CodeWriter) handleSegmentCmd(cmdType, arg1, arg2 string, buf bytes.Buffer) {
	switch arg1 {
	case "static":
		cw.templates[cmdType+" "+arg1].Execute(&buf, map[string]interface{}{
			"x": arg2,
		})
		cw.out.Write(buf.Bytes())
		return
	case "pointer":
		segment := pointerToSegmentName[arg2]
		cw.templates[cmdType+" "+arg1].Execute(&buf, map[string]interface{}{
			"arg":                   arg2,
			"segment":               segment,
			"segment_register":      segmentRegisters[segment],
			"segment_register_name": segmentRegisterName[segment],
		})
		cw.out.Write(buf.Bytes())
		return
	case "temp":
		cw.templates[cmdType+" "+arg1].Execute(&buf, map[string]interface{}{"x": arg2})
		cw.out.Write(buf.Bytes())
	default: // "local" "argument" "this" "that"
		cw.templates[cmdType+" "+arg1].Execute(&buf, map[string]interface{}{
			"segment":               arg1,
			"segment_register_name": segmentRegisterName[arg1],
			"segment_register":      segmentRegisters[arg1],
			"x":                     arg2,
		})
		cw.out.Write(buf.Bytes())
		return
	}
}

func (cw *CodeWriter) handleArithmeticCmd(cmdType, arg1 string, buf bytes.Buffer) {
	ariType, ok := arithmeticType[arg1]
	if !ok {
		log.Fatalf("not a valid type of arithmetic command: cmdType %s, arg1 %s", cmdType, arg1)
	}

	switch ariType {
	case "Negation":
		cw.templates[cmdType+" "+arg1].Execute(&buf, map[string]interface{}{
			"negation":          arg1,
			"negation_operator": negationOperators[arg1],
		})
		cw.out.Write(buf.Bytes())
		return

	case "Comparison":
		cw.templates[cmdType+" "+arg1].Execute(&buf, map[string]interface{}{
			"comp":          arg1,
			"comp_operator": comparisonOperators[arg1],
			"comp_verbose":  comparisonVerbose[arg1],
			"comp_counter":  cw.compCounter,
		})
		cw.compCounter++
		cw.out.Write(buf.Bytes())
		return
	case "Calculation":
		cw.templates[cmdType+" "+arg1].Execute(&buf, map[string]interface{}{
			"calculation":        arg1,
			"calculation_symbol": calculationSymbols[arg1],
		})
		cw.out.Write(buf.Bytes())
		return
	}
}

func (cw CodeWriter) WriteNewline() {
	cw.out.Write([]byte{'\n'})
}

func (cw CodeWriter) CloseFile() {
	cw.out.Close()
}

func isNegation(cmdType, arg1 string) bool {
	return (cmdType == C_ARITHMETIC) && (arg1 == "neg" || arg1 == "not")
}

func isComparison(cmdType, arg1 string) bool {
	return (cmdType == C_ARITHMETIC) && (arg1 == "eq" || arg1 == "lt" || arg1 == "gt")
}

func isCalculation(cmdType, arg1 string) bool {
	return (cmdType == C_ARITHMETIC) && (arg1 == "add" || arg1 == "sub" || arg1 == "and" || arg1 == "or")
}
