package assembler

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const varOffset = 16

func main() {
	target := os.Args[1]
	Assemble(target)
}

func Assemble(target string) {
	parser, _ := NewAssembler(target)
	code := parser.Assemble()
	FileSave(code, "add/Add.hack")
}

func FileSave(text, targetFP string) {
	os.WriteFile("add/Add.hack", []byte(text), 0644)
}

type Assembler struct {
	raw        string
	cmds       []string
	varCounter int
}

func removeWhiteSpace(line string) string {
	line = strings.TrimSpace(line)
	splits := strings.Split(line, " ")
	var elems []string
	for _, s := range splits {
		elems = append(elems, strings.TrimSpace(s))
	}
	return strings.Join(elems, "")
}

func NewAssembler(fp string) (*Assembler, error) {
	content, err := os.ReadFile(fp)
	if err != nil {
		return nil, fmt.Errorf("newAssembler: file found: %w", err)
	}
	return &Assembler{raw: string(content)}, nil
}

func cleanLine(line string) string {
	line = removeWhiteSpace(line)
	return removeComments(line)
}

func (asm *Assembler) FirstPass() {
	cmdCounter := 0
	for _, line := range strings.Split(string(asm.raw), "\n") {
		line = cleanLine(line)
		if len(line) == 0 {
			continue
		}
		if isLabelSymbol(line) {
			addToSymbolTable(line, cmdCounter)
			continue
		}
		asm.cmds = append(asm.cmds, line)
		cmdCounter++
	}
}

func (asm *Assembler) Assemble() string {
	var hackCmds []string
	for _, cmd := range asm.cmds {
		hackCmds = append(hackCmds, asm.TranslateCmd(cmd))
	}
	return strings.Join(hackCmds, "\n") + "\n"
}

func (p *Assembler) parseAInstruction(cmd string) string {
	val := strings.Split(cmd, "@")[1]

	address, ok := symbolTable[val]
	var err error
	if !ok {
		address, err = strconv.Atoi(val)
		if err != nil {
			symbolTable[val] = varOffset + p.varCounter
			p.varCounter++
			address = symbolTable[val]
		}
	}

	binaryStr := strconv.FormatInt(int64(address), 2)
	missingZeros := 16 - len(binaryStr)
	for range missingZeros {
		binaryStr = "0" + binaryStr
	}
	return binaryStr
}

func (p *Assembler) jump(cmd string) string {
	if strings.Contains(cmd, ";") {
		jumpPart := strings.Split(cmd, ";")[1]
		return jumpTable[jumpPart]
	}
	return "000"
}

func (p *Assembler) dest(cmd string) string {
	if strings.Contains(cmd, "=") {
		destPart := strings.Split(cmd, "=")[0]
		return destTable[destPart]
	}
	return "000"
}

func (p *Assembler) comp(cmd string) string {
	if strings.Contains(cmd, "=") {
		compString := strings.Split(cmd, "=")[1]
		if strings.Contains(compString, "M") {
			a := "1"
			return a + compTableWithM[compString]
		} else {
			a := "0"
			return a + compTableWithOutM[compString]
		}
	}
	compString := strings.Split(cmd, ";")[0]
	if strings.Contains(compString, "M") {
		a := "1"
		return a + compTableWithM[compString]
	} else {
		a := "0"
		return a + compTableWithOutM[compString]
	}
}

func (asm *Assembler) parseCInstruction(cmd string) string {
	return "111" + asm.comp(cmd) + asm.dest(cmd) + asm.jump(cmd)
}

func (asm *Assembler) TranslateCmd(cmd string) string {
	switch cmd[0] {
	case '@':
		return asm.parseAInstruction(cmd)
	default:
		return asm.parseCInstruction(cmd)
	}
}

func addToSymbolTable(cmd string, cmdCounter int) {
	symbol := cmd[1 : len(cmd)-1]
	symbolTable[symbol] = cmdCounter
}

func isLabelSymbol(cmd string) bool {
	return cmd[0] == '('
}

func removeComments(line string) string {
	line = strings.Split(line, "//")[0]
	return line
}

// func (ci CInstruction) jumpCode(cmd string) string {
// 	if strings.Contains(cmd, ";") {
// 		jumpPart := strings.Split(cmd, ";")[1]
// 		return jumpTable[jumpPart]
// 	}
// 	return "000"
// }

// func (ci CInstruction) destCode(cmd string) string {
// 	if strings.Contains(cmd, "=") {
// 		destPart := strings.Split(cmd, "=")[0]
// 		return destTable[destPart]
// 	}
// 	return "000"
// }

// func (ci CInstruction) compCode(cmd string) string {
// 	if strings.Contains(cmd, "=") {
// 		compString := strings.Split(cmd, "=")[1]
// 		if strings.Contains(compString, "M") {
// 			a := "1"
// 			return a + compTableWithM[compString]
// 		} else {
// 			a := "0"
// 			return a + compTableWithOutM[compString]
// 		}
// 	}
// 	compString := strings.Split(cmd, ";")[0]
// 	if strings.Contains(compString, "M") {
// 		a := "1"
// 		return a + compTableWithM[compString]
// 	} else {
// 		a := "0"
// 		return a + compTableWithOutM[compString]
// 	}
// }

// type AInstruction struct {
// 	instruction string
// 	address     string
// }

// type CInstruction struct {
// 	instruction string
// 	comp        string
// 	dest        string
// 	jump        string
// }
