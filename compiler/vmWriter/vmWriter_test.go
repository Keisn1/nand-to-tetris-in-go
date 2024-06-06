package vmWriter_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"hack/compiler/vmWriter"

	"github.com/stretchr/testify/assert"
)

func Test_VmWriter(t *testing.T) {
	dir := t.TempDir()

	out := filepath.Join(dir, "test.vm")
	vw, err := vmWriter.NewVmWriter(out)
	assert.NoError(t, err)

	want := ""

	segment := vmWriter.CONST
	index := 1
	vw.WritePush(segment, index)
	want += fmt.Sprintln("push", segment, index)

	segment = vmWriter.ARG
	index = 2
	vw.WritePop(segment, index)
	want += fmt.Sprintln("pop", segment, index)

	command := vmWriter.ADD
	vw.WriteArithmetic(command)
	want += fmt.Sprintln(vmWriter.ADD)

	label := "newLabel"
	vw.WriteLabel(label)
	want += fmt.Sprintln("label", label)

	gotoLabel := "gotoLabel"
	vw.WriteGoto(gotoLabel)
	want += fmt.Sprintln("goto", gotoLabel)

	ifLabel := "ifLabel"
	vw.WriteIf(ifLabel)
	want += fmt.Sprintln("if-goto", ifLabel)

	callName := "Foo.bar"
	nArgs := 10
	vw.WriteCall(callName, nArgs)
	want += fmt.Sprintln("call", callName, nArgs)

	funName := "Foo.bar"
	nLocals := 5
	vw.WriteFunction(funName, nLocals)
	want += fmt.Sprintln("function", funName, nLocals)

	vw.WriteReturn()
	want += fmt.Sprintln("return")

	vw.Close()
	assert.FileExists(t, out)

	got, err := os.ReadFile(out)

	assert.NoError(t, err)
	assert.Equal(t, want, string(got))
}
