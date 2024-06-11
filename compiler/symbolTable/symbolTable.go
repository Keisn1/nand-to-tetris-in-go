package symbolTable

type Symbol struct {
	name    string
	varType string
	kind    string
	index   int
}
type SymbolTable struct {
	classLevel            map[string]Symbol
	countsClassLevel      map[string]int
	subroutineLevel       map[string]Symbol
	countsSubroutineLevel map[string]int
}

func NewSymbolTable() SymbolTable {
	return SymbolTable{
		classLevel:            make(map[string]Symbol),
		subroutineLevel:       make(map[string]Symbol),
		countsClassLevel:      make(map[string]int),
		countsSubroutineLevel: make(map[string]int),
	}
}

func (st *SymbolTable) StartSubroutine() {
	st.subroutineLevel = make(map[string]Symbol)
	st.countsSubroutineLevel = make(map[string]int)
}

func (st *SymbolTable) Define(name, varType, kind string) {
	if kind == STATIC || kind == FIELD {
		st.classLevel[name] = Symbol{name: name, varType: varType, kind: kind, index: st.countsClassLevel[kind]}
		st.countsClassLevel[kind]++
		return
	}
	st.subroutineLevel[name] = Symbol{name: name, varType: varType, kind: kind, index: st.countsSubroutineLevel[kind]}
	st.countsSubroutineLevel[kind]++
}

func (st SymbolTable) VarCount(kind string) int {
	if kind == STATIC || kind == FIELD {
		return st.countsClassLevel[kind]
	}
	return st.countsSubroutineLevel[kind]
}

func (st SymbolTable) KindOf(name string) string {
	if _, ok := st.subroutineLevel[name]; ok {
		return st.subroutineLevel[name].kind
	}
	return st.classLevel[name].kind
}

func (st SymbolTable) TypeOf(name string) string {
	if _, ok := st.subroutineLevel[name]; ok {
		return st.subroutineLevel[name].varType
	}
	return st.classLevel[name].varType
}

func (st SymbolTable) IndexOf(name string) int {
	if _, ok := st.subroutineLevel[name]; ok {
		return st.subroutineLevel[name].index
	}
	return st.classLevel[name].index
}
