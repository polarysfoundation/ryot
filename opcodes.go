package ryot

const (
	OP_FUNC     = 0x01
	OP_PUSH     = 0x02
	OP_LOAD_ARG = 0x03
	OP_ADD      = 0x10
	OP_SUB      = 0x11
	OP_MUL      = 0x12
	OP_DIV      = 0x13
	OP_MOD      = 0x14
	OP_EQ       = 0x20
	OP_NEQ      = 0x21
	OP_LT       = 0x22
	OP_LTE      = 0x23
	OP_GT       = 0x24
	OP_GTE      = 0x25
	OP_RETURN   = 0xFF
	OP_SLOAD    = 0x30
	OP_SSTORE   = 0x31
)
