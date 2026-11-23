CMPQ SP, 0x10(R14)
JBE 0x49e18c
PUSHQ BP
MOVQ SP, BP
SUBQ $0x38, SP
LEAQ 0xb46b(IP), DX
MOVQ DX, 0x28(SP)
LEAQ 0x311af(IP), DX
MOVQ DX, 0x30(SP)
MOVQ os.Stdout(SB), BX
LEAQ go:itab.*os.File,io.Writer(SB), AX
LEAQ 0x28(SP), CX
MOVL $0x1, DI
MOVQ DI, SI
CALL fmt.Fprintln(SB)
ADDQ $0x38, SP
POPQ BP
RET
CALL runtime.morestack_noctxt.abi0(SB)
JMP main.main(SB)
