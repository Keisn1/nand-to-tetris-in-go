package symbolTable_test

import (
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"hack/compiler/symbolTable"
	"hack/compiler/token"
)

func Test_SymbolTable(t *testing.T) {
	sTab := symbolTable.NewSymbolTable()
	wantVarCount := 0
	gotVarCount := sTab.VarCount(token.STATIC)
	assert.Equal(t, wantVarCount, gotVarCount)

	testCases := []struct {
		varName      string
		wantType     string
		wantKind     string
		wantIndex    int
		wantVarCount int
	}{
		// | name | type  | kind     | # |
		// |------+-------+----------+---|
		// | x1   | int   | static   | 0 |
		// | s1   | char  | field    | 0 |
		// | p1   | point | static   | 1 |
		// | a1   | int   | argument | 0 |
		// | l1   | char  | local    | 0 |
		{randomS(15), token.INT, symbolTable.STATIC, 0, 1},
		{randomS(15), token.CHAR, symbolTable.FIELD, 0, 1},
		{randomS(15), "Point", symbolTable.STATIC, 1, 2},
		{randomS(15), token.INT, symbolTable.ARG, 0, 1},
		{randomS(15), token.CHAR, symbolTable.VAR, 0, 1},
		{randomS(15), token.BOOLEAN, symbolTable.VAR, 1, 2},
	}

	for _, tc := range testCases {
		sTab.Define(tc.varName, tc.wantType, tc.wantKind)
		assert.Equal(t, tc.wantType, sTab.TypeOf(tc.varName))
		assert.Equal(t, tc.wantKind, sTab.KindOf(tc.varName))
		assert.Equal(t, tc.wantIndex, sTab.IndexOf(tc.varName))
		assert.Equal(t, tc.wantVarCount, sTab.VarCount(tc.wantKind))
	}

	sTab.StartSubroutine()
	for _, tc := range testCases[:3] {
		assert.Equal(t, tc.wantType, sTab.TypeOf(tc.varName))
		assert.Equal(t, tc.wantKind, sTab.KindOf(tc.varName))
		assert.Equal(t, tc.wantIndex, sTab.IndexOf(tc.varName))
	}
	assert.Equal(t, 2, sTab.VarCount(symbolTable.STATIC))
	assert.Equal(t, 1, sTab.VarCount(symbolTable.FIELD))
	assert.Equal(t, 0, sTab.VarCount(symbolTable.VAR))
	assert.Equal(t, 0, sTab.VarCount(symbolTable.ARG))
	// // | name | type  | kind     | # |
	// // |------+-------+----------+---|
	// // | x1   | int   | static   | 0 |
	// // | s1   | char  | field    | 0 |
	// // | p1   | point | static   | 1 |
}

func randomS(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	var result []byte
	for i := 0; i < length; i++ {
		randomIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result = append(result, charset[randomIndex.Int64()])
	}

	return string(result)
}
