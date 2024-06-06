package symbolTable

type Symbol struct {
	name    string
	varType string
	kind    string
	index   int
}
type SymbolTable struct {
	symbols map[string]Symbol
	counts  map[string]int
}

func NewSymbolTable() SymbolTable {
	return SymbolTable{
		symbols: make(map[string]Symbol),
		counts:  make(map[string]int),
	}
}

func (st *SymbolTable) StartSubroutine() {
	st.symbols = make(map[string]Symbol)
	st.counts = make(map[string]int)
}

func (st *SymbolTable) Define(name, varType, kind string) {
	st.symbols[name] = Symbol{name: name, varType: varType, kind: kind, index: st.counts[kind]}
	st.counts[kind]++
}

func (st SymbolTable) VarCount(kind string) int {
	return st.counts[kind]
}

func (st SymbolTable) KindOf(name string) string {
	return st.symbols[name].kind
}

func (st SymbolTable) TypeOf(name string) string {
	return st.symbols[name].varType
}

func (st SymbolTable) IndexOf(name string) int {
	return st.symbols[name].index
}
