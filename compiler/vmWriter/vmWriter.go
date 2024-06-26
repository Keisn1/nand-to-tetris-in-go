package vmWriter

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	ADD   = "add"
	SUB   = "sub"
	NEG   = "neg"
	EQUAL = "eq"
	GT    = "gt"
	LT    = "lt"
	AND   = "and"
	OR    = "or"
	NOT   = "not"
)

const (
	ARG     = "argument"
	LOCAL   = "local"
	STATIC  = "static"
	THIS    = "this"
	THAT    = "that"
	POINTER = "pointer"
	TEMP    = "temp"
	CONST   = "constant"
)

const (
	NULL  = 0
	FALSE = 0
	TRUE  = 1 // needs to be negated inside the engine
)

type VmWriter struct {
	file *os.File
}

func NewVmWriter(fp string) VmWriter {
	file, err := os.Create(fp)
	if err != nil {
		log.Fatalln("NewVmWriter:", err)
	}
	return VmWriter{file: file}
}

func (vw VmWriter) GetFilename() string {
	base := filepath.Base(vw.file.Name())
	name := strings.TrimSuffix(base, filepath.Ext(base))
	return strings.ToUpper(name[0:1]) + name[1:]
}

func (vw VmWriter) WritePush(segment string, index int) {
	fmt.Fprintln(vw.file, "push", segment, index)
}

func (vw VmWriter) WritePop(segment string, index int) {
	fmt.Fprintln(vw.file, "pop", segment, index)
}

func (vw VmWriter) WriteArithmetic(cmd string) {
	fmt.Fprintln(vw.file, cmd)
}

func (vw VmWriter) WriteLabel(label string) {
	fmt.Fprintln(vw.file, "label", label)
}

func (vw VmWriter) WriteGoto(label string) {
	fmt.Fprintln(vw.file, "goto", label)
}

func (vw VmWriter) WriteIf(label string) {
	fmt.Fprintln(vw.file, "if-goto", label)
}

func (vw VmWriter) WriteCall(name string, nArgs int) {
	fmt.Fprintln(vw.file, "call", name, nArgs)
}

func (vw VmWriter) WriteFunction(name string, nLocals int) {
	fmt.Fprintln(vw.file, "function", name, nLocals)
}

func (vw VmWriter) WriteReturn() {
	fmt.Fprintln(vw.file, "return")
}

func (vw VmWriter) Close() {
	vw.file.Close()
}
