##
# From Nand to Tetris
#
# @file
# @version 0.1

build-assembler:
	go build -o cmd/hack_assembler/hack_assembler cmd/hack_assembler/main.go

install-assembler:
	@go install ./cmd/hack_assembler/

# end
