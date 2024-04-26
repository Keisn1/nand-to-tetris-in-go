package vmtrans

import (
	"bytes"
	"errors"
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

type Parser struct {
	cmds []string
	pos  int
}

func NewParser(fp string) (*Parser, error) {
	file, _ := os.ReadFile(fp)
	cmds := strings.Split(string(file), "\n")
	if len(cmds) == 1 && cmds[0] == "" {
		return &Parser{}, nil
	}
	return &Parser{cmds: cmds}, nil
}

func (p *Parser) Arg2() (string, error) {
	if cmdType, _ := p.CommandType(); cmdType == C_ARITHMETIC {
		return "", errors.New("should not be called with that command type")
	}

	cmd, err := p.CurrentCmd()
	if err != nil {
		return "", fmt.Errorf("commandType: %w", err)
	}
	args := strings.Split(cmd, " ")
	return args[2], nil
}

func (p *Parser) Arg1() (string, error) {
	cmd, err := p.CurrentCmd()
	if err != nil {
		return "", fmt.Errorf("commandType: %w", err)
	}

	if cmdType, _ := p.CommandType(); cmdType == C_ARITHMETIC {
		return cmd, nil
	}

	args := strings.Split(cmd, " ")
	return args[1], nil
}

func (p *Parser) CommandType() (string, error) {
	cmd, err := p.CurrentCmd()
	if err != nil {
		return "", fmt.Errorf("commandType: %w", err)
	}
	args := strings.Split(cmd, " ")
	return cmdTable[args[0]], nil
}

func (p *Parser) HasMoreCommands() bool {
	return p.pos < len(p.cmds) && len(p.cmds) != 0
}

func (p *Parser) Advance() {
	if p.pos < len(p.cmds) {
		p.pos++
	}
}

func (p *Parser) CurrentCmd() (string, error) {
	if p.pos < 1 {
		return "", fmt.Errorf("currentCmd: %w", ErrNotAdvanced)
	}
	return p.cmds[p.pos-1], nil
}
