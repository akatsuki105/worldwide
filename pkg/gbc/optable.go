package gbc

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

type Inst struct {
	id                 int
	Operand1, Operand2 int
	Cycle1, Cycle2     int // cond is true(1)/false(2)
	Handler            func(*GBC, int, int)
}

var nilOpcode = Inst{id: INS_NONE}

var gbz80insts [256]Inst = [256]Inst{
	/* 0x0x */ {INS_NOP, 0, 0, 1, 1, nop}, {INS_LD, BC, 0, 3, 3, ld16i}, {INS_LD, BC, A, 2, 2, ldm16r}, {INS_INC, BC, 0, 2, 2, inc16}, {INS_INC, B, 0, 1, 1, inc8}, {INS_DEC, B, 0, 1, 1, dec8}, {INS_LD, B, 0, 2, 2, ld8i}, {INS_RLCA, 0, 0, 1, 1, rlca}, {INS_LD, OP_a16_PAREN, OP_SP, 5, 5, op0x08}, {INS_ADD, HL, BC, 2, 2, addHL}, {INS_LD, A, BC, 2, 2, ld8m}, {INS_DEC, BC, 0, 2, 2, dec16}, {INS_INC, C, 0, 1, 1, inc8}, {INS_DEC, C, 0, 1, 1, dec8}, {INS_LD, C, 0, 2, 2, ld8i}, {INS_RRCA, 0, 0, 1, 1, rrca},
	/* 0x1x */ {INS_STOP, 0, 0, 1, 1, stop}, {INS_LD, DE, 0, 3, 3, ld16i}, {INS_LD, DE, A, 2, 2, ldm16r}, {INS_INC, DE, 0, 2, 2, inc16}, {INS_INC, D, 0, 1, 1, inc8}, {INS_DEC, D, 0, 1, 1, dec8}, {INS_LD, D, 0, 2, 2, ld8i}, {INS_RLA, 0, 0, 1, 1, rla}, {INS_JR, 0, 0, 0, 0, jr}, {INS_ADD, HL, DE, 2, 2, addHL}, {INS_LD, A, DE, 2, 2, ld8m}, {INS_DEC, DE, 0, 2, 2, dec16}, {INS_INC, E, 0, 1, 1, inc8}, {INS_DEC, E, 0, 1, 1, dec8}, {INS_LD, E, 0, 2, 2, ld8i}, {INS_RRA, 0, 0, 1, 1, rra},
	/* 0x2x */ {INS_JR, flagZ, 0, 0, 0, jrncc}, {INS_LD, HL, 0, 3, 3, ld16i}, {INS_LD, HLI, A, 2, 2, ldm16r}, {INS_INC, HL, 0, 2, 2, inc16}, {INS_INC, H, 0, 1, 1, inc8}, {INS_DEC, H, 0, 1, 1, dec8}, {INS_LD, H, 0, 2, 2, ld8i}, {INS_DAA, 0, 0, 1, 1, daa}, {INS_JR, flagZ, 0, 0, 0, jrcc}, {INS_ADD, HL, HL, 2, 2, addHL}, {INS_LD, A, HLI, 2, 2, ld8m}, {INS_DEC, HL, 0, 2, 2, dec16}, {INS_INC, L, 0, 1, 1, inc8}, {INS_DEC, L, 0, 1, 1, dec8}, {INS_LD, L, 0, 2, 2, ld8i}, {INS_CPL, 0, 0, 1, 1, cpl},
	/* 0x3x */ {INS_JR, flagC, 0, 0, 0, jrncc}, {INS_LD, SP, 0, 3, 3, ld16i}, {INS_LD, HLD, A, 2, 2, ldm16r}, {INS_INC, SP, 0, 2, 2, inc16}, {INS_INC, 0, 0, 2, 2, incHL}, {INS_DEC, 0, 0, 2, 2, decHL}, {INS_LD, OP_HL_PAREN, OP_d8, 2, 2, op0x36}, {INS_SCF, 0, 0, 1, 1, scf}, {INS_JR, flagC, 0, 0, 0, jrcc}, {INS_ADD, HL, SP, 2, 2, addHL}, {INS_LD, A, HLD, 2, 2, ld8m}, {INS_DEC, SP, 0, 2, 2, dec16}, {INS_INC, A, 0, 1, 1, inc8}, {INS_DEC, A, 0, 1, 1, dec8}, {INS_LD, A, 0, 2, 2, ld8i}, {INS_CCF, 0, 0, 1, 1, ccf},
	/* 0x4x */ {INS_LD, B, B, 1, 1, ld8r}, {INS_LD, B, C, 1, 1, ld8r}, {INS_LD, B, D, 1, 1, ld8r}, {INS_LD, B, E, 1, 1, ld8r}, {INS_LD, B, H, 1, 1, ld8r}, {INS_LD, B, L, 1, 1, ld8r}, {INS_LD, B, HL, 2, 2, ld8m}, {INS_LD, B, A, 1, 1, ld8r}, {INS_LD, C, B, 1, 1, ld8r}, {INS_LD, C, C, 1, 1, ld8r}, {INS_LD, C, D, 1, 1, ld8r}, {INS_LD, C, E, 1, 1, ld8r}, {INS_LD, C, H, 1, 1, ld8r}, {INS_LD, C, L, 1, 1, ld8r}, {INS_LD, C, HL, 2, 2, ld8m}, {INS_LD, C, A, 1, 1, ld8r},
	/* 0x5x */ {INS_LD, D, B, 1, 1, ld8r}, {INS_LD, D, C, 1, 1, ld8r}, {INS_LD, D, D, 1, 1, ld8r}, {INS_LD, D, E, 1, 1, ld8r}, {INS_LD, D, H, 1, 1, ld8r}, {INS_LD, D, L, 1, 1, ld8r}, {INS_LD, D, HL, 2, 2, ld8m}, {INS_LD, D, A, 1, 1, ld8r}, {INS_LD, E, B, 1, 1, ld8r}, {INS_LD, E, C, 1, 1, ld8r}, {INS_LD, E, D, 1, 1, ld8r}, {INS_LD, E, E, 1, 1, ld8r}, {INS_LD, E, H, 1, 1, ld8r}, {INS_LD, E, L, 1, 1, ld8r}, {INS_LD, E, HL, 2, 2, ld8m}, {INS_LD, E, A, 1, 1, ld8r},
	/* 0x6x */ {INS_LD, H, B, 1, 1, ld8r}, {INS_LD, H, C, 1, 1, ld8r}, {INS_LD, H, D, 1, 1, ld8r}, {INS_LD, H, E, 1, 1, ld8r}, {INS_LD, H, H, 1, 1, ld8r}, {INS_LD, H, L, 1, 1, ld8r}, {INS_LD, H, HL, 2, 2, ld8m}, {INS_LD, H, A, 1, 1, ld8r}, {INS_LD, L, B, 1, 1, ld8r}, {INS_LD, L, C, 1, 1, ld8r}, {INS_LD, L, D, 1, 1, ld8r}, {INS_LD, L, E, 1, 1, ld8r}, {INS_LD, L, H, 1, 1, ld8r}, {INS_LD, L, L, 1, 1, ld8r}, {INS_LD, L, HL, 2, 2, ld8m}, {INS_LD, L, A, 1, 1, ld8r},
	/* 0x7x */ {INS_LD, 0, B, 2, 2, ldHLR8}, {INS_LD, 0, C, 2, 2, ldHLR8}, {INS_LD, 0, D, 2, 2, ldHLR8}, {INS_LD, 0, E, 2, 2, ldHLR8}, {INS_LD, 0, H, 2, 2, ldHLR8}, {INS_LD, 0, L, 2, 2, ldHLR8}, {INS_HALT, 0, 0, 1, 1, halt}, {INS_LD, 0, A, 2, 2, ldHLR8}, {INS_LD, A, B, 1, 1, ld8r}, {INS_LD, A, C, 1, 1, ld8r}, {INS_LD, A, D, 1, 1, ld8r}, {INS_LD, A, E, 1, 1, ld8r}, {INS_LD, A, H, 1, 1, ld8r}, {INS_LD, A, L, 1, 1, ld8r}, {INS_LD, A, HL, 2, 2, ld8m}, {INS_LD, A, A, 1, 1, ld8r},
	/* 0x8x */ {INS_ADD, 0, B, 1, 1, add8}, {INS_ADD, 0, C, 1, 1, add8}, {INS_ADD, 0, D, 1, 1, add8}, {INS_ADD, 0, E, 1, 1, add8}, {INS_ADD, 0, H, 1, 1, add8}, {INS_ADD, 0, L, 1, 1, add8}, {INS_ADD, OP_A, OP_HL_PAREN, 2, 2, addaHL}, {INS_ADD, 0, A, 1, 1, add8}, {INS_ADC, 0, B, 1, 1, adc8}, {INS_ADC, 0, C, 1, 1, adc8}, {INS_ADC, 0, D, 1, 1, adc8}, {INS_ADC, 0, E, 1, 1, adc8}, {INS_ADC, 0, H, 1, 1, adc8}, {INS_ADC, 0, L, 1, 1, adc8}, {INS_ADC, 0, 0, 2, 2, adcaHL}, {INS_ADC, 0, A, 1, 1, adc8},
	/* 0x9x */ {INS_SUB, 0, B, 1, 1, sub8}, {INS_SUB, 0, C, 1, 1, sub8}, {INS_SUB, 0, D, 1, 1, sub8}, {INS_SUB, 0, E, 1, 1, sub8}, {INS_SUB, 0, H, 1, 1, sub8}, {INS_SUB, 0, L, 1, 1, sub8}, {INS_SUB, 0, 0, 2, 2, subaHL}, {INS_SUB, 0, A, 1, 1, sub8}, {INS_SBC, 0, B, 1, 1, sbc8}, {INS_SBC, 0, C, 1, 1, sbc8}, {INS_SBC, 0, D, 1, 1, sbc8}, {INS_SBC, 0, E, 1, 1, sbc8}, {INS_SBC, 0, H, 1, 1, sbc8}, {INS_SBC, 0, L, 1, 1, sbc8}, {INS_SBC, 0, 0, 2, 2, sbcaHL}, {INS_SBC, 0, A, 1, 1, sbc8},
	/* 0xax */ {INS_AND, A, B, 1, 1, and8}, {INS_AND, A, C, 1, 1, and8}, {INS_AND, A, D, 1, 1, and8}, {INS_AND, A, E, 1, 1, and8}, {INS_AND, A, H, 1, 1, and8}, {INS_AND, A, L, 1, 1, and8}, {INS_AND, OP_HL_PAREN, OP_NONE, 2, 2, andaHL}, {INS_AND, A, A, 1, 1, and8}, {INS_XOR, 0, B, 1, 1, xor8}, {INS_XOR, 0, C, 1, 1, xor8}, {INS_XOR, 0, D, 1, 1, xor8}, {INS_XOR, 0, E, 1, 1, xor8}, {INS_XOR, 0, H, 1, 1, xor8}, {INS_XOR, 0, L, 1, 1, xor8}, {INS_XOR, OP_HL_PAREN, OP_NONE, 2, 2, xoraHL}, {INS_XOR, 0, A, 1, 1, xor8},
	/* 0xbx */ {INS_OR, A, B, 1, 1, or8}, {INS_OR, A, C, 1, 1, or8}, {INS_OR, A, D, 1, 1, or8}, {INS_OR, A, E, 1, 1, or8}, {INS_OR, A, H, 1, 1, or8}, {INS_OR, A, L, 1, 1, or8}, {INS_OR, OP_HL_PAREN, OP_NONE, 2, 2, oraHL}, {INS_OR, A, A, 1, 1, or8}, {INS_CP, 0, B, 1, 1, cp}, {INS_CP, 0, C, 1, 1, cp}, {INS_CP, 0, D, 1, 1, cp}, {INS_CP, 0, E, 1, 1, cp}, {INS_CP, 0, H, 1, 1, cp}, {INS_CP, 0, L, 1, 1, cp}, {INS_CP, OP_HL_PAREN, OP_NONE, 2, 2, cpaHL}, {INS_CP, 0, A, 1, 1, cp},
	/* 0xcx */ {INS_RET, flagZ, 0, 2, 2, retncc}, {INS_POP, C, B, 2, 2, pop}, {INS_JP, flagZ, 0, 1, 1, jpncc}, {INS_JP, 0, 0, 2, 2, jp}, {INS_CALL, flagZ, 0, 0, 0, callncc}, {INS_PUSH, B, C, 2, 2, push}, {INS_ADD, OP_A, OP_d8, 2, 2, addu8}, {INS_RST, 0x00, 0, 4, 4, rst}, {INS_RET, flagZ, 0, 2, 2, retcc}, {INS_RET, 0, 0, 4, 4, ret}, {INS_JP, flagZ, 0, 1, 1, jpcc}, {INS_PREFIX, 0, 0, 0, 0, prefixCB}, {INS_CALL, flagZ, 0, 0, 0, callcc}, {INS_CALL, 0, 0, 0, 0, call}, {INS_ADC, 0, 0, 2, 2, adcu8}, {INS_RST, 0x08, 0, 4, 4, rst},
	/* 0xdx */ {INS_RET, flagC, 0, 2, 2, retncc}, {INS_POP, E, D, 2, 2, pop}, {INS_JP, flagC, 0, 1, 1, jpncc}, nilOpcode, {INS_CALL, flagC, 0, 0, 0, callncc}, {INS_PUSH, D, E, 2, 2, push}, {INS_SUB, 0, 0, 2, 2, subu8}, {INS_RST, 0x10, 0, 4, 4, rst}, {INS_RET, flagC, 0, 2, 2, retcc}, {INS_RETI, 0, 0, 4, 4, reti}, {INS_JP, flagC, 0, 1, 1, jpcc}, nilOpcode, {INS_CALL, flagC, 0, 0, 0, callcc}, nilOpcode, {INS_SBC, 0, 0, 2, 2, sbcu8}, {INS_RST, 0x18, 0, 4, 4, rst},
	/* 0xex */ {INS_LDH, OP_a8_PAREN, OP_A, 2, 2, op0xe0}, {INS_POP, L, H, 2, 2, pop}, {INS_LD, OP_C_PAREN, OP_A, 2, 2, op0xe2}, nilOpcode, nilOpcode, {INS_PUSH, H, L, 2, 2, push}, {INS_AND, OP_d8, OP_NONE, 2, 2, andu8}, {INS_RST, 0x20, 0, 4, 4, rst}, {INS_ADD, 0, 0, 4, 4, addSPi8}, {INS_JP, 0, 0, 1, 1, jpHL}, {INS_LD, OP_a16_PAREN, OP_A, 2, 2, op0xea}, nilOpcode, nilOpcode, nilOpcode, {INS_XOR, OP_d8, OP_NONE, 2, 2, xoru8}, {INS_RST, 0x28, 0, 4, 4, rst},
	/* 0xfx */ {INS_LDH, OP_A, OP_a8_PAREN, 2, 2, op0xf0}, {INS_POP, 0, 0, 2, 2, popAF}, {INS_LD, OP_A, OP_C_PAREN, 2, 2, op0xf2}, {INS_DI, 0, 0, 1, 1, di}, nilOpcode, {INS_PUSH, A, F, 2, 2, pushAF}, {INS_OR, OP_d8, OP_NONE, 2, 2, oru8}, {INS_RST, 0x30, 0, 4, 4, rst}, {INS_LD, OP_HL, OP_SP_PLUS_r8, 3, 3, op0xf8}, {INS_LD, OP_SP, OP_HL, 2, 2, op0xf9}, {INS_LD, OP_A, OP_a16_PAREN, 2, 2, ldau16}, {INS_EI, 0, 0, 1, 1, ei}, nilOpcode, nilOpcode, {INS_CP, OP_d8, OP_NONE, 2, 2, cpu8}, {INS_RST, 0x38, 0, 4, 4, rst},
}

var gbz80instsCb [256]Inst = [256]Inst{
	/* 0x0x */ {INS_RLC, B, 0, 2, 2, rlc}, {INS_RLC, C, 0, 2, 2, rlc}, {INS_RLC, D, 0, 2, 2, rlc}, {INS_RLC, E, 0, 2, 2, rlc}, {INS_RLC, H, 0, 2, 2, rlc}, {INS_RLC, L, 0, 2, 2, rlc}, {INS_RLC, 0, 0, 3, 3, rlcHL}, {INS_RLC, A, 0, 2, 2, rlc}, {INS_RRC, B, 0, 2, 2, rrc}, {INS_RRC, C, 0, 2, 2, rrc}, {INS_RRC, D, 0, 2, 2, rrc}, {INS_RRC, E, 0, 2, 2, rrc}, {INS_RRC, H, 0, 2, 2, rrc}, {INS_RRC, L, 0, 2, 2, rrc}, {INS_RRC, 0, 0, 3, 3, rrcHL}, {INS_RRC, A, 0, 2, 2, rrc},
	/* 0x1x */ {INS_RL, 0, B, 2, 2, rl}, {INS_RL, 0, C, 2, 2, rl}, {INS_RL, 0, D, 2, 2, rl}, {INS_RL, 0, E, 2, 2, rl}, {INS_RL, 0, H, 2, 2, rl}, {INS_RL, 0, L, 2, 2, rl}, {INS_RL, 0, 0, 3, 3, rlHL}, {INS_RL, 0, A, 2, 2, rl}, {INS_RR, B, 0, 2, 2, rr}, {INS_RR, C, 0, 2, 2, rr}, {INS_RR, D, 0, 2, 2, rr}, {INS_RR, E, 0, 2, 2, rr}, {INS_RR, H, 0, 2, 2, rr}, {INS_RR, L, 0, 2, 2, rr}, {INS_RR, 0, 0, 3, 3, rrHL}, {INS_RR, A, 0, 2, 2, rr},
	/* 0x2x */ {INS_SLA, B, 0, 2, 2, sla}, {INS_SLA, C, 0, 2, 2, sla}, {INS_SLA, D, 0, 2, 2, sla}, {INS_SLA, E, 0, 2, 2, sla}, {INS_SLA, H, 0, 2, 2, sla}, {INS_SLA, L, 0, 2, 2, sla}, {INS_SLA, 0, 0, 3, 3, slaHL}, {INS_SLA, A, 0, 2, 2, sla}, {INS_SRA, B, 0, 2, 2, sra}, {INS_SRA, C, 0, 2, 2, sra}, {INS_SRA, D, 0, 2, 2, sra}, {INS_SRA, E, 0, 2, 2, sra}, {INS_SRA, H, 0, 2, 2, sra}, {INS_SRA, L, 0, 2, 2, sra}, {INS_SRA, 0, 0, 3, 3, sraHL}, {INS_SRA, A, 0, 2, 2, sra},
	/* 0x3x */ {INS_SWAP, 0, B, 2, 2, swap}, {INS_SWAP, 0, C, 2, 2, swap}, {INS_SWAP, 0, D, 2, 2, swap}, {INS_SWAP, 0, E, 2, 2, swap}, {INS_SWAP, 0, H, 2, 2, swap}, {INS_SWAP, 0, L, 2, 2, swap}, {INS_SWAP, 0, 0, 3, 3, swapHL}, {INS_SWAP, 0, A, 2, 2, swap}, {INS_SRL, B, 0, 2, 2, srl}, {INS_SRL, C, 0, 2, 2, srl}, {INS_SRL, D, 0, 2, 2, srl}, {INS_SRL, E, 0, 2, 2, srl}, {INS_SRL, H, 0, 2, 2, srl}, {INS_SRL, L, 0, 2, 2, srl}, {INS_SRL, 0, 0, 3, 3, srlHL}, {INS_SRL, A, 0, 2, 2, srl},

	/* 0x4x */ {INS_BIT, 0, B, 2, 2, bit}, {INS_BIT, 0, C, 2, 2, bit}, {INS_BIT, 0, D, 2, 2, bit}, {INS_BIT, 0, E, 2, 2, bit}, {INS_BIT, 0, H, 2, 2, bit}, {INS_BIT, 0, L, 2, 2, bit}, {INS_BIT, 0, 0, 3, 3, bitHL}, {INS_BIT, 0, A, 2, 2, bit}, {INS_BIT, 1, B, 2, 2, bit}, {INS_BIT, 1, C, 2, 2, bit}, {INS_BIT, 1, D, 2, 2, bit}, {INS_BIT, 1, E, 2, 2, bit}, {INS_BIT, 1, H, 2, 2, bit}, {INS_BIT, 1, L, 2, 2, bit}, {INS_BIT, 1, 0, 3, 3, bitHL}, {INS_BIT, 1, A, 2, 2, bit},
	/* 0x5x */ {INS_BIT, 2, B, 2, 2, bit}, {INS_BIT, 2, C, 2, 2, bit}, {INS_BIT, 2, D, 2, 2, bit}, {INS_BIT, 2, E, 2, 2, bit}, {INS_BIT, 2, H, 2, 2, bit}, {INS_BIT, 2, L, 2, 2, bit}, {INS_BIT, 2, 0, 3, 3, bitHL}, {INS_BIT, 2, A, 2, 2, bit}, {INS_BIT, 3, B, 2, 2, bit}, {INS_BIT, 3, C, 2, 2, bit}, {INS_BIT, 3, D, 2, 2, bit}, {INS_BIT, 3, E, 2, 2, bit}, {INS_BIT, 3, H, 2, 2, bit}, {INS_BIT, 3, L, 2, 2, bit}, {INS_BIT, 3, 0, 3, 3, bitHL}, {INS_BIT, 3, A, 2, 2, bit},
	/* 0x6x */ {INS_BIT, 4, B, 2, 2, bit}, {INS_BIT, 4, C, 2, 2, bit}, {INS_BIT, 4, D, 2, 2, bit}, {INS_BIT, 4, E, 2, 2, bit}, {INS_BIT, 4, H, 2, 2, bit}, {INS_BIT, 4, L, 2, 2, bit}, {INS_BIT, 4, 0, 3, 3, bitHL}, {INS_BIT, 4, A, 2, 2, bit}, {INS_BIT, 5, B, 2, 2, bit}, {INS_BIT, 5, C, 2, 2, bit}, {INS_BIT, 5, D, 2, 2, bit}, {INS_BIT, 5, E, 2, 2, bit}, {INS_BIT, 5, H, 2, 2, bit}, {INS_BIT, 5, L, 2, 2, bit}, {INS_BIT, 5, 0, 3, 3, bitHL}, {INS_BIT, 5, A, 2, 2, bit},
	/* 0x7x */ {INS_BIT, 6, B, 2, 2, bit}, {INS_BIT, 6, C, 2, 2, bit}, {INS_BIT, 6, D, 2, 2, bit}, {INS_BIT, 6, E, 2, 2, bit}, {INS_BIT, 6, H, 2, 2, bit}, {INS_BIT, 6, L, 2, 2, bit}, {INS_BIT, 6, 0, 3, 3, bitHL}, {INS_BIT, 6, A, 2, 2, bit}, {INS_BIT, 7, B, 2, 2, bit}, {INS_BIT, 7, C, 2, 2, bit}, {INS_BIT, 7, D, 2, 2, bit}, {INS_BIT, 7, E, 2, 2, bit}, {INS_BIT, 7, H, 2, 2, bit}, {INS_BIT, 7, L, 2, 2, bit}, {INS_BIT, 7, 0, 3, 3, bitHL}, {INS_BIT, 7, A, 2, 2, bit},

	/* 0x8x */ {INS_RES, 0, B, 2, 2, res}, {INS_RES, 0, C, 2, 2, res}, {INS_RES, 0, D, 2, 2, res}, {INS_RES, 0, E, 2, 2, res}, {INS_RES, 0, H, 2, 2, res}, {INS_RES, 0, L, 2, 2, res}, {INS_RES, 0, 0, 3, 3, resHL}, {INS_RES, 0, A, 2, 2, res}, {INS_RES, 1, B, 2, 2, res}, {INS_RES, 1, C, 2, 2, res}, {INS_RES, 1, D, 2, 2, res}, {INS_RES, 1, E, 2, 2, res}, {INS_RES, 1, H, 2, 2, res}, {INS_RES, 1, L, 2, 2, res}, {INS_RES, 1, 0, 3, 3, resHL}, {INS_RES, 1, A, 2, 2, res},
	/* 0x9x */ {INS_RES, 2, B, 2, 2, res}, {INS_RES, 2, C, 2, 2, res}, {INS_RES, 2, D, 2, 2, res}, {INS_RES, 2, E, 2, 2, res}, {INS_RES, 2, H, 2, 2, res}, {INS_RES, 2, L, 2, 2, res}, {INS_RES, 2, 0, 3, 3, resHL}, {INS_RES, 2, A, 2, 2, res}, {INS_RES, 3, B, 2, 2, res}, {INS_RES, 3, C, 2, 2, res}, {INS_RES, 3, D, 2, 2, res}, {INS_RES, 3, E, 2, 2, res}, {INS_RES, 3, H, 2, 2, res}, {INS_RES, 3, L, 2, 2, res}, {INS_RES, 3, 0, 3, 3, resHL}, {INS_RES, 3, A, 2, 2, res},
	/* 0xax */ {INS_RES, 4, B, 2, 2, res}, {INS_RES, 4, C, 2, 2, res}, {INS_RES, 4, D, 2, 2, res}, {INS_RES, 4, E, 2, 2, res}, {INS_RES, 4, H, 2, 2, res}, {INS_RES, 4, L, 2, 2, res}, {INS_RES, 4, 0, 3, 3, resHL}, {INS_RES, 4, A, 2, 2, res}, {INS_RES, 5, B, 2, 2, res}, {INS_RES, 5, C, 2, 2, res}, {INS_RES, 5, D, 2, 2, res}, {INS_RES, 5, E, 2, 2, res}, {INS_RES, 5, H, 2, 2, res}, {INS_RES, 5, L, 2, 2, res}, {INS_RES, 5, 0, 3, 3, resHL}, {INS_RES, 5, A, 2, 2, res},
	/* 0xbx */ {INS_RES, 6, B, 2, 2, res}, {INS_RES, 6, C, 2, 2, res}, {INS_RES, 6, D, 2, 2, res}, {INS_RES, 6, E, 2, 2, res}, {INS_RES, 6, H, 2, 2, res}, {INS_RES, 6, L, 2, 2, res}, {INS_RES, 6, 0, 3, 3, resHL}, {INS_RES, 6, A, 2, 2, res}, {INS_RES, 7, B, 2, 2, res}, {INS_RES, 7, C, 2, 2, res}, {INS_RES, 7, D, 2, 2, res}, {INS_RES, 7, E, 2, 2, res}, {INS_RES, 7, H, 2, 2, res}, {INS_RES, 7, L, 2, 2, res}, {INS_RES, 7, 0, 3, 3, resHL}, {INS_RES, 7, A, 2, 2, res},

	/* 0xcx */ {INS_SET, 0, B, 2, 2, set}, {INS_SET, 0, C, 2, 2, set}, {INS_SET, 0, D, 2, 2, set}, {INS_SET, 0, E, 2, 2, set}, {INS_SET, 0, H, 2, 2, set}, {INS_SET, 0, L, 2, 2, set}, {INS_SET, 0, 0, 3, 3, setHL}, {INS_SET, 0, A, 2, 2, set}, {INS_SET, 1, B, 2, 2, set}, {INS_SET, 1, C, 2, 2, set}, {INS_SET, 1, D, 2, 2, set}, {INS_SET, 1, E, 2, 2, set}, {INS_SET, 1, H, 2, 2, set}, {INS_SET, 1, L, 2, 2, set}, {INS_SET, 1, 0, 3, 3, setHL}, {INS_SET, 1, A, 2, 2, set},
	/* 0xdx */ {INS_SET, 2, B, 2, 2, set}, {INS_SET, 2, C, 2, 2, set}, {INS_SET, 2, D, 2, 2, set}, {INS_SET, 2, E, 2, 2, set}, {INS_SET, 2, H, 2, 2, set}, {INS_SET, 2, L, 2, 2, set}, {INS_SET, 2, 0, 3, 3, setHL}, {INS_SET, 2, A, 2, 2, set}, {INS_SET, 3, B, 2, 2, set}, {INS_SET, 3, C, 2, 2, set}, {INS_SET, 3, D, 2, 2, set}, {INS_SET, 3, E, 2, 2, set}, {INS_SET, 3, H, 2, 2, set}, {INS_SET, 3, L, 2, 2, set}, {INS_SET, 3, 0, 3, 3, setHL}, {INS_SET, 3, A, 2, 2, set},
	/* 0xex */ {INS_SET, 4, B, 2, 2, set}, {INS_SET, 4, C, 2, 2, set}, {INS_SET, 4, D, 2, 2, set}, {INS_SET, 4, E, 2, 2, set}, {INS_SET, 4, H, 2, 2, set}, {INS_SET, 4, L, 2, 2, set}, {INS_SET, 4, 0, 3, 3, setHL}, {INS_SET, 4, A, 2, 2, set}, {INS_SET, 5, B, 2, 2, set}, {INS_SET, 5, C, 2, 2, set}, {INS_SET, 5, D, 2, 2, set}, {INS_SET, 5, E, 2, 2, set}, {INS_SET, 5, H, 2, 2, set}, {INS_SET, 5, L, 2, 2, set}, {INS_SET, 5, 0, 3, 3, setHL}, {INS_SET, 5, A, 2, 2, set},
	/* 0xfx */ {INS_SET, 6, B, 2, 2, set}, {INS_SET, 6, C, 2, 2, set}, {INS_SET, 6, D, 2, 2, set}, {INS_SET, 6, E, 2, 2, set}, {INS_SET, 6, H, 2, 2, set}, {INS_SET, 6, L, 2, 2, set}, {INS_SET, 6, 0, 3, 3, setHL}, {INS_SET, 6, A, 2, 2, set}, {INS_SET, 7, B, 2, 2, set}, {INS_SET, 7, C, 2, 2, set}, {INS_SET, 7, D, 2, 2, set}, {INS_SET, 7, E, 2, 2, set}, {INS_SET, 7, H, 2, 2, set}, {INS_SET, 7, L, 2, 2, set}, {INS_SET, 7, 0, 3, 3, setHL}, {INS_SET, 7, A, 2, 2, set},
}
