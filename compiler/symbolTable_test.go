package compiler_test

import (
	"crypto/rand"
	"hack/compiler"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SymbolTable(t *testing.T) {
	sTab := compiler.NewSymbolTable()
	wantVarCount := 0
	gotVarCount := sTab.VarCount(compiler.STATIC)
	assert.Equal(t, wantVarCount, gotVarCount)

	testCases := []struct {
		varName      string
		wantType     string
		wantKind     string
		wantIndex    int
		wantVarCount int
	}{
		{randomS(15), compiler.INT, compiler.STATIC, 0, 1},
		{randomS(15), compiler.CHAR, compiler.FIELD, 0, 1},
		{randomS(15), "Point", compiler.STATIC, 1, 2},
	}

	for _, tc := range testCases {
		sTab.Define(tc.varName, tc.wantType, tc.wantKind)
		assert.Equal(t, tc.wantType, sTab.TypeOf(tc.varName))
		assert.Equal(t, tc.wantKind, sTab.KindOf(tc.varName))
		assert.Equal(t, tc.wantIndex, sTab.IndexOf(tc.varName))
		assert.Equal(t, tc.wantVarCount, sTab.VarCount(tc.wantKind))
	}
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
