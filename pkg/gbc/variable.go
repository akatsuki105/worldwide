package gbc

const (
	OAM uint16 = 0xfe00
)

const (
	INS_ADC = iota
	INS_AND
	INS_ADD
	INS_CP
	INS_DEC
	INS_INC
	INS_OR
	INS_SBC
	INS_SUB
	INS_XOR
	INS_BIT
	INS_RES
	INS_SET
	INS_SWAP
	INS_RL
	INS_RLA
	INS_RLC
	INS_RLCA
	INS_RR
	INS_RRA
	INS_RRC
	INS_RRCA
	INS_SLA
	INS_SRA
	INS_SRL
	INS_LD
	INS_CALL
	INS_JP
	INS_JR
	INS_RET
	INS_RETI
	INS_RST
	INS_POP
	INS_PUSH
	INS_CCF
	INS_CPL
	INS_DAA
	INS_DI
	INS_EI
	INS_HALT
	INS_NOP
	INS_SCF
	INS_STOP
	INS_PREFIX
	INS_NONE
	INS_LDH
)

const (
	OP_NONE = iota
	OP_BC
	OP_d16
	OP_BC_PAREN
	OP_A
	OP_AF
	OP_B
	OP_d8
	OP_a16_PAREN
	OP_SP
	OP_SP_PLUS_r8
	OP_HL
	OP_C
	OP_C_PAREN
	OP_DE
	OP_DE_PAREN
	OP_D
	OP_r8
	OP_a8_PAREN
	OP_a16
	OP_E
	OP_NZ
	OP_NC
	OP_H
	OP_L
	OP_Z
	OP_HLPLUS_PAREN
	OP_HLMINUS_PAREN
	OP_HL_PAREN
)

type Opcode struct {
	Ins      int
	Operand1 int
	Operand2 int
	Cycle1   int
	Cycle2   int // cond is false
	Handler  func(*GBC, int, int)
}

var nilOpcode = Opcode{Ins: INS_NONE}

var opcodes [256]Opcode = [256]Opcode{
	/* 0x0x */ {INS_NOP, 0, 0, 1, 1, nop}, {INS_LD, BC, 0, 3, 3, ld16i}, {INS_LD, BC, A, 2, 2, ldm16r}, {INS_INC, BC, 0, 2, 2, inc16}, {INS_INC, B, 0, 1, 1, inc8}, {INS_DEC, B, 0, 1, 1, dec8}, {INS_LD, B, 0, 2, 2, ld8i}, {INS_RLCA, 0, 0, 1, 1, rlca}, {INS_LD, OP_a16_PAREN, OP_SP, 0, 0, op0x08}, {INS_ADD, HL, BC, 2, 2, addHL}, {INS_LD, A, BC, 2, 2, ld8m}, {INS_DEC, BC, 0, 2, 2, dec16}, {INS_INC, C, 0, 1, 1, inc8}, {INS_DEC, C, 0, 1, 1, dec8}, {INS_LD, C, 0, 2, 2, ld8i}, {INS_RRCA, 0, 0, 1, 1, rrca},
	/* 0x1x */ {INS_STOP, 0, 0, 1, 1, stop}, {INS_LD, DE, 0, 3, 3, ld16i}, {INS_LD, DE, A, 2, 2, ldm16r}, {INS_INC, DE, 0, 2, 2, inc16}, {INS_INC, D, 0, 1, 1, inc8}, {INS_DEC, D, 0, 1, 1, dec8}, {INS_LD, D, 0, 2, 2, ld8i}, {INS_RLA, 0, 0, 1, 1, rla}, {INS_JR, 0, 0, 0, 0, jr}, {INS_ADD, HL, DE, 2, 2, addHL}, {INS_LD, A, DE, 2, 2, ld8m}, {INS_DEC, DE, 0, 2, 2, dec16}, {INS_INC, E, 0, 1, 1, inc8}, {INS_DEC, E, 0, 1, 1, dec8}, {INS_LD, E, 0, 2, 2, ld8i}, {INS_RRA, 0, 0, 1, 1, rra},
	/* 0x2x */ {INS_JR, flagZ, 0, 0, 0, jrncc}, {INS_LD, HL, 0, 3, 3, ld16i}, {INS_LD, HLI, A, 2, 2, ldm16r}, {INS_INC, HL, 0, 2, 2, inc16}, {INS_INC, H, 0, 1, 1, inc8}, {INS_DEC, H, 0, 1, 1, dec8}, {INS_LD, H, 0, 2, 2, ld8i}, {INS_DAA, 0, 0, 1, 1, daa}, {INS_JR, flagZ, 0, 0, 0, jrcc}, {INS_ADD, HL, HL, 2, 2, addHL}, {INS_LD, A, HLI, 2, 2, ld8m}, {INS_DEC, HL, 0, 2, 2, dec16}, {INS_INC, L, 0, 1, 1, inc8}, {INS_DEC, L, 0, 1, 1, dec8}, {INS_LD, L, 0, 2, 2, ld8i}, {INS_CPL, 0, 0, 1, 1, cpl},
	/* 0x3x */ {INS_JR, flagC, 0, 0, 0, jrncc}, {INS_LD, SP, 0, 3, 3, ld16i}, {INS_LD, HLD, A, 2, 2, ldm16r}, {INS_INC, SP, 0, 2, 2, inc16}, {INS_INC, 0, 0, 0, 0, incHL}, {INS_DEC, 0, 0, 0, 0, decHL}, {INS_LD, OP_HL_PAREN, OP_d8, 0, 0, op0x36}, {INS_SCF, 0, 0, 1, 1, scf}, {INS_JR, flagC, 0, 0, 0, jrcc}, {INS_ADD, HL, SP, 2, 2, addHL}, {INS_LD, A, HLD, 2, 2, ld8m}, {INS_DEC, SP, 0, 2, 2, dec16}, {INS_INC, A, 0, 1, 1, inc8}, {INS_DEC, A, 0, 1, 1, dec8}, {INS_LD, A, 0, 2, 2, ld8i}, {INS_CCF, 0, 0, 1, 1, ccf},
	/* 0x4x */ {INS_LD, B, B, 1, 1, ld8r}, {INS_LD, B, C, 1, 1, ld8r}, {INS_LD, B, D, 1, 1, ld8r}, {INS_LD, B, E, 1, 1, ld8r}, {INS_LD, B, H, 1, 1, ld8r}, {INS_LD, B, L, 1, 1, ld8r}, {INS_LD, B, HL, 2, 2, ld8m}, {INS_LD, B, A, 1, 1, ld8r}, {INS_LD, C, B, 1, 1, ld8r}, {INS_LD, C, C, 1, 1, ld8r}, {INS_LD, C, D, 1, 1, ld8r}, {INS_LD, C, E, 1, 1, ld8r}, {INS_LD, C, H, 1, 1, ld8r}, {INS_LD, C, L, 1, 1, ld8r}, {INS_LD, C, HL, 2, 2, ld8m}, {INS_LD, C, A, 1, 1, ld8r},
	/* 0x5x */ {INS_LD, D, B, 1, 1, ld8r}, {INS_LD, D, C, 1, 1, ld8r}, {INS_LD, D, D, 1, 1, ld8r}, {INS_LD, D, E, 1, 1, ld8r}, {INS_LD, D, H, 1, 1, ld8r}, {INS_LD, D, L, 1, 1, ld8r}, {INS_LD, D, HL, 2, 2, ld8m}, {INS_LD, D, A, 1, 1, ld8r}, {INS_LD, E, B, 1, 1, ld8r}, {INS_LD, E, C, 1, 1, ld8r}, {INS_LD, E, D, 1, 1, ld8r}, {INS_LD, E, E, 1, 1, ld8r}, {INS_LD, E, H, 1, 1, ld8r}, {INS_LD, E, L, 1, 1, ld8r}, {INS_LD, E, HL, 2, 2, ld8m}, {INS_LD, E, A, 1, 1, ld8r},
	/* 0x6x */ {INS_LD, H, B, 1, 1, ld8r}, {INS_LD, H, C, 1, 1, ld8r}, {INS_LD, H, D, 1, 1, ld8r}, {INS_LD, H, E, 1, 1, ld8r}, {INS_LD, H, H, 1, 1, ld8r}, {INS_LD, H, L, 1, 1, ld8r}, {INS_LD, H, HL, 2, 2, ld8m}, {INS_LD, H, A, 1, 1, ld8r}, {INS_LD, L, B, 1, 1, ld8r}, {INS_LD, L, C, 1, 1, ld8r}, {INS_LD, L, D, 1, 1, ld8r}, {INS_LD, L, E, 1, 1, ld8r}, {INS_LD, L, H, 1, 1, ld8r}, {INS_LD, L, L, 1, 1, ld8r}, {INS_LD, L, HL, 2, 2, ld8m}, {INS_LD, L, A, 1, 1, ld8r},
	/* 0x7x */ {INS_LD, 0, B, 2, 2, ldHLR8}, {INS_LD, 0, C, 2, 2, ldHLR8}, {INS_LD, 0, D, 2, 2, ldHLR8}, {INS_LD, 0, E, 2, 2, ldHLR8}, {INS_LD, 0, H, 2, 2, ldHLR8}, {INS_LD, 0, L, 2, 2, ldHLR8}, {INS_HALT, 0, 0, 1, 1, halt}, {INS_LD, 0, A, 2, 2, ldHLR8}, {INS_LD, A, B, 1, 1, ld8r}, {INS_LD, A, C, 1, 1, ld8r}, {INS_LD, A, D, 1, 1, ld8r}, {INS_LD, A, E, 1, 1, ld8r}, {INS_LD, A, H, 1, 1, ld8r}, {INS_LD, A, L, 1, 1, ld8r}, {INS_LD, A, HL, 2, 2, ld8m}, {INS_LD, A, A, 1, 1, ld8r},
	/* 0x8x */ {INS_ADD, 0, B, 1, 1, add8}, {INS_ADD, 0, C, 1, 1, add8}, {INS_ADD, 0, D, 1, 1, add8}, {INS_ADD, 0, E, 1, 1, add8}, {INS_ADD, 0, H, 1, 1, add8}, {INS_ADD, 0, L, 1, 1, add8}, {INS_ADD, OP_A, OP_HL_PAREN, 2, 2, addaHL}, {INS_ADD, 0, A, 1, 1, add8}, {INS_ADC, 0, B, 1, 1, adc8}, {INS_ADC, 0, C, 1, 1, adc8}, {INS_ADC, 0, D, 1, 1, adc8}, {INS_ADC, 0, E, 1, 1, adc8}, {INS_ADC, 0, H, 1, 1, adc8}, {INS_ADC, 0, L, 1, 1, adc8}, {INS_ADC, 0, 0, 2, 2, adcaHL}, {INS_ADC, 0, A, 1, 1, adc8},
	/* 0x9x */ {INS_SUB, 0, B, 1, 1, sub8}, {INS_SUB, 0, C, 1, 1, sub8}, {INS_SUB, 0, D, 1, 1, sub8}, {INS_SUB, 0, E, 1, 1, sub8}, {INS_SUB, 0, H, 1, 1, sub8}, {INS_SUB, 0, L, 1, 1, sub8}, {INS_SUB, 0, 0, 2, 2, subaHL}, {INS_SUB, 0, A, 1, 1, sub8}, {INS_SBC, 0, B, 1, 1, sbc8}, {INS_SBC, 0, C, 1, 1, sbc8}, {INS_SBC, 0, D, 1, 1, sbc8}, {INS_SBC, 0, E, 1, 1, sbc8}, {INS_SBC, 0, H, 1, 1, sbc8}, {INS_SBC, 0, L, 1, 1, sbc8}, {INS_SBC, OP_A, OP_HL_PAREN, 2, 2, nil}, {INS_SBC, 0, A, 1, 1, sbc8},
	/* 0xax */ {INS_AND, A, B, 1, 1, and8}, {INS_AND, A, C, 1, 1, and8}, {INS_AND, A, D, 1, 1, and8}, {INS_AND, A, E, 1, 1, and8}, {INS_AND, A, H, 1, 1, and8}, {INS_AND, A, L, 1, 1, and8}, {INS_AND, OP_HL_PAREN, OP_NONE, 2, 2, op0xa6}, {INS_AND, A, A, 1, 1, and8}, {INS_XOR, 0, B, 1, 1, xor8}, {INS_XOR, 0, C, 1, 1, xor8}, {INS_XOR, 0, D, 1, 1, xor8}, {INS_XOR, 0, E, 1, 1, xor8}, {INS_XOR, 0, H, 1, 1, xor8}, {INS_XOR, 0, L, 1, 1, xor8}, {INS_XOR, OP_HL_PAREN, OP_NONE, 2, 2, xoraHL}, {INS_XOR, 0, A, 1, 1, xor8},
	/* 0xbx */ {INS_OR, A, B, 1, 1, orR8}, {INS_OR, A, C, 1, 1, orR8}, {INS_OR, A, D, 1, 1, orR8}, {INS_OR, A, E, 1, 1, orR8}, {INS_OR, A, H, 1, 1, orR8}, {INS_OR, A, L, 1, 1, orR8}, {INS_OR, OP_HL_PAREN, OP_NONE, 2, 2, oraHL}, {INS_OR, A, A, 1, 1, orR8}, {INS_CP, 0, B, 1, 1, cp}, {INS_CP, 0, C, 1, 1, cp}, {INS_CP, 0, D, 1, 1, cp}, {INS_CP, 0, E, 1, 1, cp}, {INS_CP, 0, H, 1, 1, cp}, {INS_CP, 0, L, 1, 1, cp}, {INS_CP, OP_HL_PAREN, OP_NONE, 2, 2, cpaHL}, {INS_CP, 0, A, 1, 1, cp},
	/* 0xcx */ {INS_RET, flagZ, 0, 0, 0, retncc}, {INS_POP, C, B, 0, 0, pop}, {INS_JP, flagZ, 0, 0, 0, jpncc}, {INS_JP, 0, 0, 0, 0, jp}, {INS_CALL, flagZ, 0, 0, 0, callncc}, {INS_PUSH, B, C, 0, 0, push}, {INS_ADD, OP_A, OP_d8, 2, 2, addu8}, {INS_RST, 0x00, 0, 4, 4, rst}, {INS_RET, flagZ, 0, 0, 0, retcc}, {INS_RET, 0, 0, 4, 4, ret}, {INS_JP, flagZ, 0, 0, 0, jpcc}, {INS_PREFIX, 0, 0, 0, 0, prefixCB}, {INS_CALL, flagZ, 0, 0, 0, callcc}, {INS_CALL, 0, 0, 0, 0, call}, {INS_ADC, 0, 0, 2, 2, adcu8}, {INS_RST, 0x08, 0, 4, 4, rst},
	/* 0xdx */ {INS_RET, flagC, 0, 0, 0, retncc}, {INS_POP, E, D, 0, 0, pop}, {INS_JP, flagC, 0, 0, 0, jpncc}, nilOpcode, {INS_CALL, flagC, 0, 0, 0, callncc}, {INS_PUSH, D, E, 0, 0, push}, {INS_SUB, 0, 0, 2, 2, subu8}, {INS_RST, 0x10, 0, 4, 4, rst}, {INS_RET, flagC, 0, 0, 0, retcc}, {INS_RETI, 0, 0, 4, 4, reti}, {INS_JP, flagC, 0, 0, 0, jpcc}, nilOpcode, {INS_CALL, flagC, 0, 0, 0, callcc}, nilOpcode, {INS_SBC, OP_A, OP_d8, 2, 2, nil}, {INS_RST, 0x18, 0, 4, 4, rst},
	/* 0xex */ {INS_LDH, OP_a8_PAREN, OP_A, 0, 0, op0xe0}, {INS_POP, L, H, 0, 0, pop}, {INS_LD, OP_C_PAREN, OP_A, 2, 2, op0xe2}, nilOpcode, nilOpcode, {INS_PUSH, H, L, 0, 0, push}, {INS_AND, OP_d8, OP_NONE, 2, 2, andu8}, {INS_RST, 0x20, 0, 4, 4, rst}, {INS_ADD, 0, 0, 4, 4, addSPi8}, {INS_JP, 0, 0, 0, 0, jpHL}, {INS_LD, OP_a16_PAREN, OP_A, 0, 0, op0xea}, nilOpcode, nilOpcode, nilOpcode, {INS_XOR, OP_d8, OP_NONE, 2, 2, xoru8}, {INS_RST, 0x28, 0, 4, 4, rst},
	/* 0xfx */ {INS_LDH, OP_A, OP_a8_PAREN, 0, 0, op0xf0}, {INS_POP, 0, 0, 0, 0, popAF}, {INS_LD, OP_A, OP_C_PAREN, 2, 2, op0xf2}, {INS_DI, 0, 0, 1, 1, di}, nilOpcode, {INS_PUSH, A, F, 0, 0, pushAF}, {INS_OR, OP_d8, OP_NONE, 2, 2, oru8}, {INS_RST, 0x30, 0, 4, 4, rst}, {INS_LD, OP_HL, OP_SP_PLUS_r8, 3, 3, op0xf8}, {INS_LD, OP_SP, OP_HL, 2, 2, op0xf9}, {INS_LD, OP_A, OP_a16_PAREN, 0, 0, op0xfa}, {INS_EI, 0, 0, 1, 1, ei}, nilOpcode, nilOpcode, {INS_CP, OP_d8, OP_NONE, 2, 2, cpu8}, {INS_RST, 0x38, 0, 4, 4, rst},
}

// issue #10
// cycle0 opcode(0x36, 0xe0, 0xea, 0xf0, 0xfa) increments cycle in execution
var prefixCBs [256]Opcode = [256]Opcode{
	/* 0x0x */ {INS_RLC, B, 0, 2, 2, rlc}, {INS_RLC, C, 0, 2, 2, rlc}, {INS_RLC, D, 0, 2, 2, rlc}, {INS_RLC, E, 0, 2, 2, rlc}, {INS_RLC, H, 0, 2, 2, rlc}, {INS_RLC, L, 0, 2, 2, rlc}, {INS_RLC, 0, 0, 0, 0, rlcHL}, {INS_RLC, A, 0, 2, 2, rlc}, {INS_RRC, B, 0, 2, 2, rrc}, {INS_RRC, C, 0, 2, 2, rrc}, {INS_RRC, D, 0, 2, 2, rrc}, {INS_RRC, E, 0, 2, 2, rrc}, {INS_RRC, H, 0, 2, 2, rrc}, {INS_RRC, L, 0, 2, 2, rrc}, {INS_RRC, 0, 0, 0, 0, rrcHL}, {INS_RRC, A, 0, 2, 2, rrc},
	/* 0x1x */ {INS_RL, 0, B, 2, 2, rl}, {INS_RL, 0, C, 2, 2, rl}, {INS_RL, 0, D, 2, 2, rl}, {INS_RL, 0, E, 2, 2, rl}, {INS_RL, 0, H, 2, 2, rl}, {INS_RL, 0, L, 2, 2, rl}, {INS_RL, 0, 0, 0, 0, rlHL}, {INS_RL, 0, A, 2, 2, rl}, {INS_RR, B, 0, 2, 2, rr}, {INS_RR, C, 0, 2, 2, rr}, {INS_RR, D, 0, 2, 2, rr}, {INS_RR, E, 0, 2, 2, rr}, {INS_RR, H, 0, 2, 2, rr}, {INS_RR, L, 0, 2, 2, rr}, {INS_RR, 0, 0, 0, 0, rrHL}, {INS_RR, A, 0, 2, 2, rr},
	/* 0x2x */ {INS_SLA, B, 0, 2, 2, sla}, {INS_SLA, C, 0, 2, 2, sla}, {INS_SLA, D, 0, 2, 2, sla}, {INS_SLA, E, 0, 2, 2, sla}, {INS_SLA, H, 0, 2, 2, sla}, {INS_SLA, L, 0, 2, 2, sla}, {INS_SLA, 0, 0, 0, 0, slaHL}, {INS_SLA, A, 0, 2, 2, sla}, {INS_SRA, B, 0, 2, 2, sra}, {INS_SRA, C, 0, 2, 2, sra}, {INS_SRA, D, 0, 2, 2, sra}, {INS_SRA, E, 0, 2, 2, sra}, {INS_SRA, H, 0, 2, 2, sra}, {INS_SRA, L, 0, 2, 2, sra}, {INS_SRA, 0, 0, 0, 0, sraHL}, {INS_SRA, A, 0, 2, 2, sra},
	/* 0x3x */ {INS_SWAP, 0, B, 2, 2, swap}, {INS_SWAP, 0, C, 2, 2, swap}, {INS_SWAP, 0, D, 2, 2, swap}, {INS_SWAP, 0, E, 2, 2, swap}, {INS_SWAP, 0, H, 2, 2, swap}, {INS_SWAP, 0, L, 2, 2, swap}, {INS_SWAP, 0, 0, 0, 0, swapHL}, {INS_SWAP, 0, A, 2, 2, swap}, {INS_SRL, B, 0, 2, 2, srl}, {INS_SRL, C, 0, 2, 2, srl}, {INS_SRL, D, 0, 2, 2, srl}, {INS_SRL, E, 0, 2, 2, srl}, {INS_SRL, H, 0, 2, 2, srl}, {INS_SRL, L, 0, 2, 2, srl}, {INS_SRL, 0, 0, 0, 0, srlHL}, {INS_SRL, A, 0, 2, 2, srl},

	/* 0x4x */ {INS_BIT, 0, B, 2, 2, bit}, {INS_BIT, 0, C, 2, 2, bit}, {INS_BIT, 0, D, 2, 2, bit}, {INS_BIT, 0, E, 2, 2, bit}, {INS_BIT, 0, H, 2, 2, bit}, {INS_BIT, 0, L, 2, 2, bit}, {INS_BIT, 0, 0, 3, 3, bitHL}, {INS_BIT, 0, A, 2, 2, bit}, {INS_BIT, 1, B, 2, 2, bit}, {INS_BIT, 1, C, 2, 2, bit}, {INS_BIT, 1, D, 2, 2, bit}, {INS_BIT, 1, E, 2, 2, bit}, {INS_BIT, 1, H, 2, 2, bit}, {INS_BIT, 1, L, 2, 2, bit}, {INS_BIT, 1, 0, 3, 3, bitHL}, {INS_BIT, 1, A, 2, 2, bit},
	/* 0x5x */ {INS_BIT, 2, B, 2, 2, bit}, {INS_BIT, 2, C, 2, 2, bit}, {INS_BIT, 2, D, 2, 2, bit}, {INS_BIT, 2, E, 2, 2, bit}, {INS_BIT, 2, H, 2, 2, bit}, {INS_BIT, 2, L, 2, 2, bit}, {INS_BIT, 2, 0, 3, 3, bitHL}, {INS_BIT, 2, A, 2, 2, bit}, {INS_BIT, 3, B, 2, 2, bit}, {INS_BIT, 3, C, 2, 2, bit}, {INS_BIT, 3, D, 2, 2, bit}, {INS_BIT, 3, E, 2, 2, bit}, {INS_BIT, 3, H, 2, 2, bit}, {INS_BIT, 3, L, 2, 2, bit}, {INS_BIT, 3, 0, 3, 3, bitHL}, {INS_BIT, 3, A, 2, 2, bit},
	/* 0x6x */ {INS_BIT, 4, B, 2, 2, bit}, {INS_BIT, 4, C, 2, 2, bit}, {INS_BIT, 4, D, 2, 2, bit}, {INS_BIT, 4, E, 2, 2, bit}, {INS_BIT, 4, H, 2, 2, bit}, {INS_BIT, 4, L, 2, 2, bit}, {INS_BIT, 4, 0, 3, 3, bitHL}, {INS_BIT, 4, A, 2, 2, bit}, {INS_BIT, 5, B, 2, 2, bit}, {INS_BIT, 5, C, 2, 2, bit}, {INS_BIT, 5, D, 2, 2, bit}, {INS_BIT, 5, E, 2, 2, bit}, {INS_BIT, 5, H, 2, 2, bit}, {INS_BIT, 5, L, 2, 2, bit}, {INS_BIT, 5, 0, 3, 3, bitHL}, {INS_BIT, 5, A, 2, 2, bit},
	/* 0x7x */ {INS_BIT, 6, B, 2, 2, bit}, {INS_BIT, 6, C, 2, 2, bit}, {INS_BIT, 6, D, 2, 2, bit}, {INS_BIT, 6, E, 2, 2, bit}, {INS_BIT, 6, H, 2, 2, bit}, {INS_BIT, 6, L, 2, 2, bit}, {INS_BIT, 6, 0, 3, 3, bitHL}, {INS_BIT, 6, A, 2, 2, bit}, {INS_BIT, 7, B, 2, 2, bit}, {INS_BIT, 7, C, 2, 2, bit}, {INS_BIT, 7, D, 2, 2, bit}, {INS_BIT, 7, E, 2, 2, bit}, {INS_BIT, 7, H, 2, 2, bit}, {INS_BIT, 7, L, 2, 2, bit}, {INS_BIT, 7, 0, 3, 3, bitHL}, {INS_BIT, 7, A, 2, 2, bit},

	/* 0x8x */ {INS_RES, 0, B, 2, 2, res}, {INS_RES, 0, C, 2, 2, res}, {INS_RES, 0, D, 2, 2, res}, {INS_RES, 0, E, 2, 2, res}, {INS_RES, 0, H, 2, 2, res}, {INS_RES, 0, L, 2, 2, res}, {INS_RES, 0, 0, 0, 0, resHL}, {INS_RES, 0, A, 2, 2, res}, {INS_RES, 1, B, 2, 2, res}, {INS_RES, 1, C, 2, 2, res}, {INS_RES, 1, D, 2, 2, res}, {INS_RES, 1, E, 2, 2, res}, {INS_RES, 1, H, 2, 2, res}, {INS_RES, 1, L, 2, 2, res}, {INS_RES, 1, 0, 0, 0, resHL}, {INS_RES, 1, A, 2, 2, res},
	/* 0x9x */ {INS_RES, 2, B, 2, 2, res}, {INS_RES, 2, C, 2, 2, res}, {INS_RES, 2, D, 2, 2, res}, {INS_RES, 2, E, 2, 2, res}, {INS_RES, 2, H, 2, 2, res}, {INS_RES, 2, L, 2, 2, res}, {INS_RES, 2, 0, 0, 0, resHL}, {INS_RES, 2, A, 2, 2, res}, {INS_RES, 3, B, 2, 2, res}, {INS_RES, 3, C, 2, 2, res}, {INS_RES, 3, D, 2, 2, res}, {INS_RES, 3, E, 2, 2, res}, {INS_RES, 3, H, 2, 2, res}, {INS_RES, 3, L, 2, 2, res}, {INS_RES, 3, 0, 0, 0, resHL}, {INS_RES, 3, A, 2, 2, res},
	/* 0xax */ {INS_RES, 4, B, 2, 2, res}, {INS_RES, 4, C, 2, 2, res}, {INS_RES, 4, D, 2, 2, res}, {INS_RES, 4, E, 2, 2, res}, {INS_RES, 4, H, 2, 2, res}, {INS_RES, 4, L, 2, 2, res}, {INS_RES, 4, 0, 0, 0, resHL}, {INS_RES, 4, A, 2, 2, res}, {INS_RES, 5, B, 2, 2, res}, {INS_RES, 5, C, 2, 2, res}, {INS_RES, 5, D, 2, 2, res}, {INS_RES, 5, E, 2, 2, res}, {INS_RES, 5, H, 2, 2, res}, {INS_RES, 5, L, 2, 2, res}, {INS_RES, 5, 0, 0, 0, resHL}, {INS_RES, 5, A, 2, 2, res},
	/* 0xbx */ {INS_RES, 6, B, 2, 2, res}, {INS_RES, 6, C, 2, 2, res}, {INS_RES, 6, D, 2, 2, res}, {INS_RES, 6, E, 2, 2, res}, {INS_RES, 6, H, 2, 2, res}, {INS_RES, 6, L, 2, 2, res}, {INS_RES, 6, 0, 0, 0, resHL}, {INS_RES, 6, A, 2, 2, res}, {INS_RES, 7, B, 2, 2, res}, {INS_RES, 7, C, 2, 2, res}, {INS_RES, 7, D, 2, 2, res}, {INS_RES, 7, E, 2, 2, res}, {INS_RES, 7, H, 2, 2, res}, {INS_RES, 7, L, 2, 2, res}, {INS_RES, 7, 0, 0, 0, resHL}, {INS_RES, 7, A, 2, 2, res},

	/* 0xcx */ {INS_SET, 0, B, 2, 2, set}, {INS_SET, 0, C, 2, 2, set}, {INS_SET, 0, D, 2, 2, set}, {INS_SET, 0, E, 2, 2, set}, {INS_SET, 0, H, 2, 2, set}, {INS_SET, 0, L, 2, 2, set}, {INS_SET, 0, 0, 0, 0, setHL}, {INS_SET, 0, A, 2, 2, set}, {INS_SET, 1, B, 2, 2, set}, {INS_SET, 1, C, 2, 2, set}, {INS_SET, 1, D, 2, 2, set}, {INS_SET, 1, E, 2, 2, set}, {INS_SET, 1, H, 2, 2, set}, {INS_SET, 1, L, 2, 2, set}, {INS_SET, 1, 0, 0, 0, setHL}, {INS_SET, 1, A, 2, 2, set},
	/* 0xdx */ {INS_SET, 2, B, 2, 2, set}, {INS_SET, 2, C, 2, 2, set}, {INS_SET, 2, D, 2, 2, set}, {INS_SET, 2, E, 2, 2, set}, {INS_SET, 2, H, 2, 2, set}, {INS_SET, 2, L, 2, 2, set}, {INS_SET, 2, 0, 0, 0, setHL}, {INS_SET, 2, A, 2, 2, set}, {INS_SET, 3, B, 2, 2, set}, {INS_SET, 3, C, 2, 2, set}, {INS_SET, 3, D, 2, 2, set}, {INS_SET, 3, E, 2, 2, set}, {INS_SET, 3, H, 2, 2, set}, {INS_SET, 3, L, 2, 2, set}, {INS_SET, 3, 0, 0, 0, setHL}, {INS_SET, 3, A, 2, 2, set},
	/* 0xex */ {INS_SET, 4, B, 2, 2, set}, {INS_SET, 4, C, 2, 2, set}, {INS_SET, 4, D, 2, 2, set}, {INS_SET, 4, E, 2, 2, set}, {INS_SET, 4, H, 2, 2, set}, {INS_SET, 4, L, 2, 2, set}, {INS_SET, 4, 0, 0, 0, setHL}, {INS_SET, 4, A, 2, 2, set}, {INS_SET, 5, B, 2, 2, set}, {INS_SET, 5, C, 2, 2, set}, {INS_SET, 5, D, 2, 2, set}, {INS_SET, 5, E, 2, 2, set}, {INS_SET, 5, H, 2, 2, set}, {INS_SET, 5, L, 2, 2, set}, {INS_SET, 5, 0, 0, 0, setHL}, {INS_SET, 5, A, 2, 2, set},
	/* 0xfx */ {INS_SET, 6, B, 2, 2, set}, {INS_SET, 6, C, 2, 2, set}, {INS_SET, 6, D, 2, 2, set}, {INS_SET, 6, E, 2, 2, set}, {INS_SET, 6, H, 2, 2, set}, {INS_SET, 6, L, 2, 2, set}, {INS_SET, 6, 0, 0, 0, setHL}, {INS_SET, 6, A, 2, 2, set}, {INS_SET, 7, B, 2, 2, set}, {INS_SET, 7, C, 2, 2, set}, {INS_SET, 7, D, 2, 2, set}, {INS_SET, 7, E, 2, 2, set}, {INS_SET, 7, H, 2, 2, set}, {INS_SET, 7, L, 2, 2, set}, {INS_SET, 7, 0, 0, 0, setHL}, {INS_SET, 7, A, 2, 2, set},
}

var (
	icon []byte = []byte{137, 80, 78, 71, 13, 10, 26, 10, 0, 0, 0, 13, 73, 72, 68, 82, 0, 0, 0, 32, 0, 0, 0, 32, 8, 6, 0, 0, 0, 115, 122, 122, 244, 0, 0, 0, 4, 115, 66, 73, 84, 8, 8, 8, 8, 124, 8, 100, 136, 0, 0, 0, 9, 112, 72, 89, 115, 0, 0, 14, 196, 0, 0, 14, 196, 1, 149, 43, 14, 27, 0, 0, 6, 9, 73, 68, 65, 84, 88, 133, 205,
		151, 89, 108, 19, 87, 20, 134, 255, 153, 241, 120, 137, 215, 196, 78, 216, 2, 33, 16, 32, 16, 168, 237, 76, 10, 73, 139, 128, 182, 168, 5, 26, 162, 96, 182, 22, 84, 250, 64, 27, 182, 7, 164, 46, 82, 43, 81, 241, 90, 30, 90, 218, 74, 21, 130, 7, 30, 104, 33, 15, 109, 165, 182, 44, 37, 108,
		21, 107, 33, 177, 147, 16, 2, 134, 176, 56, 118, 2, 222, 98, 59, 241, 62, 99, 103, 166, 15, 1, 7, 39, 142, 13, 18, 168, 61, 210, 72, 115, 239, 57, 62, 255, 119, 175, 206, 185, 115, 13, 252, 199, 70, 60, 239, 15, 74, 74, 74, 164, 186, 162, 137, 13, 154, 130, 252, 15, 102, 76, 159, 54, 143, 231, 57, 220, 127, 224, 184, 209, 31, 24, 248, 169, 207, 243, 232, 128, 221, 110, 143, 191, 12, 80, 0, 0, 195, 84, 87, 214, 214, 173, 234, 58, 241, 215, 113, 33, 145, 136, 10, 130, 192, 9, 44, 219, 43, 68, 34, 221, 194, 159, 71, 27, 133, 21, 117, 117, 93, 12, 83, 93, 249, 82, 196, 13,
		85, 213, 75, 183, 110, 223, 17, 238, 239, 239, 19, 4, 129, 75, 61, 44, 219, 155, 122, 60, 30, 171, 208, 176, 181, 33, 108, 168, 170, 94, 250, 66, 197, 25, 166, 186, 114, 203, 182, 29, 225, 120, 60, 148, 38, 62, 18, 128, 101, 123, 133, 96, 240, 129, 208, 176, 245, 227, 240, 11, 219, 9,
		189, 94, 175, 169, 55, 173, 181, 141, 92, 249, 88, 0, 44, 219, 43, 120, 189, 86, 161, 206, 100, 178, 233, 245, 122, 77, 174, 252, 100, 174, 0, 90, 166, 216, 187, 251, 171, 47, 166, 170, 213, 42, 132, 184, 48, 188, 177, 190, 156, 208, 42, 149, 18, 187, 190, 220, 62, 149, 150, 41, 246,
		230, 138, 165, 178, 57, 25, 102, 193, 66, 147, 169, 238, 251, 181, 107, 234, 9, 0, 144, 80, 98, 200, 233, 60, 132, 185, 48, 154, 236, 103, 112, 202, 113, 6, 23, 29, 221, 168, 158, 80, 6, 130, 24, 106, 168, 67, 183, 46, 98, 128, 141, 98, 193, 140, 217, 240, 249, 250, 245, 193, 129, 248, 89, 167, 243, 161, 99, 44, 141, 172, 59, 160, 208, 168, 191, 222, 182, 101, 115, 90, 171, 198, 146, 49, 236, 239, 60, 136, 107, 238, 22, 184, 163, 30, 152, 221, 54, 196, 7, 19, 41, 191, 35, 232, 131, 39, 26, 4, 0, 52, 124, 180, 158, 80, 168, 149, 123, 178, 105, 140, 9, 96, 168, 170, 94, 186, 110, 77, 253, 235, 98, 133, 24, 167, 28, 103, 145, 228, 147, 0, 128, 107, 174, 22, 4, 216, 64, 42, 142, 38, 41, 72, 41, 58, 53, 222, 85, 93, 143, 53, 51, 23, 0, 0, 52, 26, 21, 86, 155, 150, 189, 150, 173, 43, 198, 4, 200, 147, 201, 62, 89, 191, 206, 132, 88, 50, 6, 219, 64, 55, 18,
		252, 208, 42, 157, 17, 87, 122, 28, 45, 78, 109, 255, 211, 198, 11, 60, 66, 92, 12, 107, 87, 47, 67, 158, 76, 250, 233, 115, 1, 24, 141, 53, 37, 111, 44, 94, 248, 142, 182, 160, 0, 58, 153, 22, 91, 230, 109, 134, 76, 36, 3, 0, 168, 196, 170, 180, 88, 110, 48, 153, 49, 241, 31, 247, 44, 216, 249, 247, 33, 40, 213, 10, 44, 94, 52, 255, 109, 163, 177, 166, 228, 153, 1, 8, 74, 216, 184, 178, 118, 121, 70, 223, 252, 241, 12, 68, 132, 40, 53, 142, 37, 19, 96, 159, 170, 129, 39, 86, 51, 113, 38, 62, 172, 88, 4, 154, 18, 97, 197, 242, 37, 36, 65, 9, 27, 159, 25, 64, 167, 211, 153, 24, 198, 144, 54, 199, 14, 114, 0, 128, 66, 89, 33, 54, 205, 222, 0, 157, 84, 11, 0, 80, 137, 165, 160, 136, 209, 105, 198, 203, 213, 88, 92, 60, 27, 0, 80, 105, 172, 128, 78, 151, 111, 202, 164, 37, 26, 57, 97, 48, 24, 10, 25, 198, 96, 20, 137, 210, 93, 87, 93, 205, 88, 60, 105, 33, 0, 96, 154, 186, 20, 59, 13, 59, 16, 73, 70, 32, 36, 7, 32, 34, 135, 186, 249, 124, 175, 21, 22, 183, 13, 250, 194, 41, 120, 115, 114, 69, 170, 54, 68, 34, 10, 70, 99, 133, 209, 227, 114, 22, 182, 183, 183, 123, 179, 238, 0, 73, 75, 23, 50, 149, 250, 81, 243, 79, 196, 159, 24, 65, 16, 80, 208, 10, 72, 30, 119, 192, 169, 238, 14, 28, 232, 56, 7, 139, 219, 134, 131, 157, 231, 209, 120, 251, 74, 90, 188, 209, 80, 65, 146, 180, 52, 61, 73, 38, 0, 66, 128, 177, 188, 124, 214, 200, 233, 172, 230, 141, 6, 113, 100, 132, 224, 9, 219, 117, 120, 31, 159, 7, 0, 80, 62, 171, 20, 132, 0,
		99, 78, 0, 129, 64, 217, 148, 201, 197, 207, 5, 240, 219, 221, 102, 36, 248, 193, 244, 60, 16, 208, 209, 55, 124, 0, 22, 79, 158, 0, 129, 64, 89, 78, 0, 177, 84, 82, 164, 84, 42, 158, 89, 220, 27, 13, 226, 242, 163, 174, 140, 190, 64, 60, 146, 122, 87, 42, 228, 16, 75, 37, 69, 185, 1, 104, 154, 30, 57, 151, 205, 142, 219, 218, 193, 11, 66, 70, 159, 132, 74, 47, 100, 49, 37, 26, 149, 123, 20, 64, 52, 18, 11, 112, 28, 247, 76, 226, 209, 68, 20, 231, 123, 172, 99, 250, 139, 149, 218, 212, 59, 199, 37, 16, 137, 198, 251, 115, 2, 12, 242, 124, 87, 183, 61, 243, 199, 43, 196,
		133, 32, 60, 181, 218, 102, 183, 25, 28, 159, 249, 36, 148, 82, 52, 230, 20, 76, 74, 141, 237, 142, 135, 224, 5, 254, 78, 78, 0, 2, 184, 208, 98, 110, 27, 149, 208, 29, 245, 96, 143, 229, 91, 88, 3, 67, 57, 120, 129, 199, 53, 87, 75, 70, 113, 0, 88, 50, 121, 14, 36, 79, 237, 184, 165, 245, 38, 8, 224, 66, 78, 128, 96, 192, 123, 186, 169, 233, 140, 63, 204, 133, 17, 228, 134, 219, 72, 39, 213, 162, 182, 116, 57, 74, 85, 67, 71, 186, 61, 228, 64, 40, 17, 206, 40, 174, 145, 228, 97, 85, 89, 21, 0, 192, 31, 15, 99, 128, 141, 226, 244, 233, 75, 254, 96, 192, 123, 122, 100, 236, 168, 11, 137, 223, 239, 31, 20, 209, 82, 165, 87, 227, 93, 116, 147, 191, 5, 131, 238, 21, 208, 20, 13, 146, 32, 81, 172, 152, 4, 154, 28, 90, 149, 35, 212, 131, 91, 126, 43, 124, 145, 244, 2, 148, 82, 52, 62, 127, 181, 22, 227, 229, 26, 132, 185, 56, 118, 95, 249, 21, 231, 46, 53, 227, 122, 147, 229, 155, 155, 157, 29, 185, 1, 0, 64, 165, 156, 98, 238, 115, 186, 54, 21, 84, 106, 84, 110, 214, 131, 121, 218, 185, 32, 71, 156, 247, 106, 177, 10, 214, 192, 29, 244, 244, 15, 183, 90, 137, 74, 139, 207, 170, 222, 69, 169, 186, 8, 73, 126, 16, 223, 181, 158, 196, 253, 62, 23, 238, 54, 154, 31, 134, 220, 201, 13, 62, 95, 207, 168, 234, 206, 8, 224, 243, 245, 112, 90, 133, 174, 61, 234, 99, 55, 138, 103, 137, 72, 87, 204, 141, 217, 5, 229, 160, 136, 225, 112, 154, 164, 81, 85, 84, 137, 98, 121, 30, 202, 243, 39, 96, 229, 180, 74, 172, 159, 85, 3, 141, 84, 14, 110, 48, 137, 31, 218, 154, 208, 238, 182, 195, 123, 212, 154, 140, 247, 244, 175, 238, 236, 248, 103, 84, 1, 142, 9, 0, 0, 78, 231, 67, 155, 156, 200, 183, 199, 252, 241, 149, 68, 41, 200, 132, 192, 97, 102, 254, 140, 148, 223, 30, 116, 160, 47, 214, 7, 49, 65, 96, 130, 66, 131, 233, 154, 113, 169, 143, 207, 225, 219, 151, 113, 209, 110, 133, 247, 248, 157, 100, 168, 211, 189, 185, 205, 124, 245, 247, 177, 116, 178, 94, 74, 93, 206, 222, 235, 10, 34, 191, 45, 122, 63, 182, 172, 190, 166, 86, 54, 190, 104, 92, 202, 167, 145, 168, 161, 149, 105, 33, 2, 11, 133, 88, 10, 154, 28, 78, 197, 123, 98, 56, 185, 255, 184, 127, 160, 203, 243, 94, 155, 249, 234, 47, 217, 52, 178, 2, 0, 128, 243, 81, 111, 151, 74, 34, 255, 233, 216, 177, 179, 133, 247, 238, 63, 168, 144, 203, 229, 228, 184, 113, 133, 160, 40, 234, 113, 2, 22, 52, 73, 129, 227, 18, 104, 49, 223, 192, 190, 253, 71, 146, 251, 126, 60, 252, 115, 212, 231, 55, 181, 154, 155, 91, 115, 229, 127, 174, 63, 167, 12, 195, 76, 225, 73, 241, 251, 178, 60, 233, 91, 211, 167, 78, 157, 171, 82, 169, 10, 121, 129, 67, 40, 24, 246, 62, 176, 245, 118, 198, 226, 241, 179, 36, 207,
		53, 90, 44, 150, 49, 175, 225, 255, 59, 251, 23, 242, 156, 177, 101, 178, 79, 133, 255, 0, 0, 0, 0, 73, 69, 78, 68, 174, 66, 96, 130}
)
