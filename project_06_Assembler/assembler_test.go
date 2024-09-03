package assembler_test

import (
	"hack/project_06_Assembler"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFile(t *testing.T) {
	type testCase struct {
		targetFP string
		wantFP   string
	}
	testCases := []testCase{
		{
			targetFP: "add/Add.asm",
			wantFP:   "add/addCmp.hack",
		},
		{
			targetFP: "max/MaxL.asm",
			wantFP:   "max/maxCmp.hack",
		},
		{
			targetFP: "max/Max.asm",
			wantFP:   "max/maxCmp.hack",
		},
		{
			targetFP: "rect/RectL.asm",
			wantFP:   "rect/rectCmp.hack",
		},
		{
			targetFP: "rect/Rect.asm",
			wantFP:   "rect/rectCmp.hack",
		},
		{
			targetFP: "pong/PongL.asm",
			wantFP:   "pong/pongCmp.hack",
		},
		{
			targetFP: "pong/Pong.asm",
			wantFP:   "pong/pongCmp.hack",
		},
	}

	for _, tc := range testCases {
		want, err := os.ReadFile("test_programs/" + tc.wantFP)
		if err != nil {
			t.Fatal(tc.wantFP, err)
		}
		asm, err := assembler.NewAssembler("test_programs/" + tc.targetFP)
		if err != nil {
			t.Fatal(tc.wantFP, err)
		}

		got := asm.Assemble()
		assert.Equal(t, string(want), got)
	}
}

func TestAssembler(t *testing.T) {
	t.Run("Throws error if file could not be found", func(t *testing.T) {
		filepath := "invalid"
		_, err := assembler.NewAssembler(filepath)
		assert.Error(t, err)
	})

	t.Run("Test cmds dependent on rect/Rect.asm", func(t *testing.T) {
		asm, err := assembler.NewAssembler("test_programs/rect/Rect.asm")
		if err != nil {
			t.Fatal(err)
		}

		asm.FirstPass()
		type testCase struct {
			cmd  string
			want string
		}

		testCases := []testCase{
			{
				cmd:  "@counter",
				want: "0000000000010000",
			},
			{
				cmd:  "@address",
				want: "0000000000010001",
			},
		}

		for _, tc := range testCases {
			got := asm.TranslateCmd(tc.cmd)
			assert.Equal(t, tc.want, got)
		}
	})

	t.Run("Test cmds dependent on max/Max.asm", func(t *testing.T) {
		asm, err := assembler.NewAssembler("test_programs/max/Max.asm")
		if err != nil {
			t.Fatal(err)
		}

		asm.FirstPass()

		type testCase struct {
			cmd  string
			want string
		}

		testCases := []testCase{
			{
				cmd:  "@INFINITE_LOOP",
				want: "0000000000001110",
			},
			{
				cmd:  "@OUTPUT_D",
				want: "0000000000001100",
			},
			{
				cmd:  "@OUTPUT_FIRST",
				want: "0000000000001010",
			},
		}

		for _, tc := range testCases {
			got := asm.TranslateCmd(tc.cmd)
			assert.Equal(t, tc.want, got)
		}
	})

	t.Run("Test cmds that are independent of file", func(t *testing.T) {
		asm := assembler.Assembler{}
		type testCase struct {
			cmd  string
			want string
		}

		testCases := []testCase{
			{
				cmd:  "@R1",
				want: "0000000000000001",
			},
			{
				cmd:  "0;JMP",
				want: "1110101010000111",
			},
			{
				cmd:  "D;JGT",
				want: "1110001100000001",
			},
			{
				cmd:  "@10",
				want: "0000000000001010",
			},
			{
				cmd:  "D=D-M",
				want: "1111010011010000",
			},
			{
				cmd:  "D=M",
				want: "1111110000010000",
			},
			{
				cmd:  "M=D",
				want: "1110001100001000",
			},
			{
				cmd:  "D=D+A",
				want: "1110000010010000",
			},
			{
				cmd:  "D=A",
				want: "1110110000010000",
			},
			{
				cmd:  "@0",
				want: "0000000000000000",
			},
			{
				cmd:  "@1",
				want: "0000000000000001",
			},
			{
				cmd:  "@2",
				want: "0000000000000010",
			},
			{
				cmd:  "@3",
				want: "0000000000000011",
			},
			{
				cmd:  "@4",
				want: "0000000000000100",
			},
		}

		for _, tc := range testCases {
			got := asm.TranslateCmd(tc.cmd)
			assert.Equal(t, tc.want, got)
		}
	})
}
