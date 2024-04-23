package vmtrans

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

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

const (
	C_ARITHMETIC = "C_ARITHMETIC"
	C_PUSH       = "C_PUSH"
	C_POP        = "C_POP"
	C_LABEL      = "C_LABEL"
	C_GOTO       = "C_GOTO"
	C_IF         = "C_IF"
	C_FUNCTION   = "C_FUNCTION"
	C_RETURN     = "C_RETURN"
	C_CALL       = "C_CALL"
)

var (
	cmdTable = map[string]string{
		"add":      C_ARITHMETIC,
		"sub":      C_ARITHMETIC,
		"neg":      C_ARITHMETIC,
		"eq":       C_ARITHMETIC,
		"get":      C_ARITHMETIC,
		"lt":       C_ARITHMETIC,
		"and":      C_ARITHMETIC,
		"or":       C_ARITHMETIC,
		"not":      C_ARITHMETIC,
		"push":     C_PUSH,
		"pop":      C_POP,
		"label":    C_LABEL,
		"goto":     C_GOTO,
		"if":       C_IF,
		"function": C_FUNCTION,
		"return":   C_RETURN,
		"call":     C_CALL,
	}
)

var (
	ErrNotAdvanced = errors.New("you need to advance before reading")
)

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
