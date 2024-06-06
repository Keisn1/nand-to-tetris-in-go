package vmWriter

import (
	"fmt"
	"os"
)

const (
	ADD = "add"
	SUB = "sub"
	NEG = "neg"
	EQ  = "eq"
	GT  = "gt"
	LT  = "lt"
	AND = "and"
	OR  = "or"
	NOT = "not"
)

const (
	ARG     = "argument"
	LOCAL   = "local"
	STATIC  = "static"
	THIS    = "this"
	THAT    = "that"
	POINTER = "pointer"
	TEMP    = "temp"
	CONST   = "const"
)

type VmWriter struct {
	file *os.File
}

func NewVmWriter(fp string) (VmWriter, error) {
	file, err := os.Create(fp)
	if err != nil {
		return VmWriter{}, fmt.Errorf("NewVmWriter: %w", err)
	}
	return VmWriter{file: file}, nil
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
