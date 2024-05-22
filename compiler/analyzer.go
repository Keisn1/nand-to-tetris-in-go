package compiler

import "os"

type Analyzer struct{}

func NewAnalyzer(fp string) Analyzer {
	return Analyzer{}
}

func (a Analyzer) Analyze() {
	f, _ := os.Create("test_files/Test1/MainT.xml")
	defer f.Close()

	f.Write([]byte(`<tokens>
  <keyword> class </keyword>
  <identifier> Main </identifier>
  <symbol> { </symbol>
  <keyword> function </keyword>
  <keyword> void </keyword>
  <identifier> main </identifier>
  <symbol> ( </symbol>
  <symbol> ) </symbol>
  <symbol> { </symbol>
  <keyword> return </keyword>
  <symbol> ; </symbol>
  <symbol> } </symbol>
  <symbol> } </symbol>
</tokens>
`))
}
