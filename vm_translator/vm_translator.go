package vmtrans

import (
	"errors"
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
		return "", errors.New("you need to advance before reading")
	}
	return p.cmds[p.pos-1], nil
}
