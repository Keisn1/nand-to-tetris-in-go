package main

import (
	"os"
	"testing"
)

func TestMain(t *testing.T) {
	os.Args = []string{"add/Add.asm"}
	main()
	t.Fatal()
}
