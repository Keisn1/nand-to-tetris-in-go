package vmtrans

import (
	"errors"
	"fmt"
	"log"
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
	file, err := os.ReadFile(fp)
	if err != nil {
		log.Fatal(err)
	}
	cmds := strings.Split(string(file), "\n")
	if len(cmds) == 1 && cmds[0] == "" {
		return &Parser{}, nil
	}
	return &Parser{cmds: cmds}, nil
}

func (p *Parser) Arg2() (string, error) {
	if cmdType, _ := p.CommandType(); cmdType == C_ARITHMETIC {
		return "", errors.New("arithmetic cmd does not have arg2")
	}

	if cmdType, _ := p.CommandType(); cmdType == C_LABEL {
		return "", errors.New("label cmd does not have arg2")
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

	cmdType, _ := p.CommandType()
	if cmdType == C_ARITHMETIC {
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
	if len(p.cmds) == 0 {
		return false
	}
	if p.pos >= len(p.cmds) {
		return false
	}

	if isEmptyLine(p.cmds[p.pos]) {
		p.pos++
		return p.HasMoreCommands()
	}

	return true
}

func (p *Parser) Advance() {
	if p.pos >= len(p.cmds) {
		return
	}

	cmd := p.cmds[p.pos]
	p.pos++

	if isEmptyLine(cmd) {
		p.Advance()
		return
	}
	if isComment(cmd) {
		p.Advance()
		return
	}

	p.curCmd = cmd
}

func (p *Parser) CurrentCmd() (string, error) {
	if p.pos == 0 {
		return "", fmt.Errorf("currentCmd: %w", ErrNotAdvanced)
	}
	return trimSpace(removeRepeatedWhitespaces(cutComment(p.curCmd))), nil
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
