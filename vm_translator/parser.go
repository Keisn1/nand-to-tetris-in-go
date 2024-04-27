package vmtrans

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Parser struct {
	cmds   []string
	pos    int
	curCmd string
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
	if p.pos >= len(p.cmds) {
		return
	}
	p.curCmd = trimSpace(removeRepeatedWhitespaces(cutComment(p.cmds[p.pos])))
	p.pos++

	if isEmptyLine(p.curCmd) {
		p.Advance()
	}
	if isComment(p.curCmd) {
		p.Advance()
	}
}

func (p *Parser) CurrentCmd() (string, error) {
	if p.pos == 0 {
		return "", fmt.Errorf("currentCmd: %w", ErrNotAdvanced)
	}
	return p.curCmd, nil
}

func trimSpace(l string) string { return strings.TrimSpace(l) }

func cutComment(l string) string { return strings.Split(l, "//")[0] }

func isEmptyLine(l string) bool { return len(strings.TrimSpace(l)) == 0 }

func isComment(l string) bool { return strings.HasPrefix(l, "//") }

func removeRepeatedWhitespaces(input string) string {
	re := regexp.MustCompile(`\s+`)
	output := re.ReplaceAllString(input, " ")
	return strings.TrimSpace(output)
}
