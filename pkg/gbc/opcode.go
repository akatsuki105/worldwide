package gbc

import (
	"fmt"
	"gbc/pkg/util"
)

func (cpu *CPU) a16Fetch() uint16 {
	value := cpu.d16Fetch()
	return value
}

func (cpu *CPU) a16FetchJP() uint16 {
	lower := uint16(cpu.FetchMemory8(cpu.Reg.PC + 1)) // M = 1: nn read: memory access for low byte
	cpu.timer(1)
	upper := uint16(cpu.FetchMemory8(cpu.Reg.PC + 2)) // M = 2: nn read: memory access for high byte
	cpu.timer(1)
	value := (upper << 8) | lower
	return value
}

func (cpu *CPU) d8Fetch() byte {
	value := cpu.FetchMemory8(cpu.Reg.PC + 1)
	return value
}

func (cpu *CPU) d16Fetch() uint16 {
	lower, upper := uint16(cpu.FetchMemory8(cpu.Reg.PC+1)), uint16(cpu.FetchMemory8(cpu.Reg.PC+2))
	return (upper << 8) | lower
}

// LD R8,R8
func ldR8R8(cpu *CPU, op1, op2 int) {
	cpu.Reg.R[op1] = cpu.Reg.R[op2]
	cpu.Reg.PC++
}

// ------ LD A, *

// LD A,(BC)
func op0x0a(cpu *CPU, operand1, operand2 int) {
	cpu.Reg.R[A] = cpu.FetchMemory8(cpu.Reg.BC())
	cpu.Reg.PC++
}

// LD A,(DE)
func op0x1a(cpu *CPU, operand1, operand2 int) {
	cpu.Reg.R[A] = cpu.FetchMemory8(cpu.Reg.DE())
	cpu.Reg.PC++
}

// LD A,(HL+)
func op0x2a(cpu *CPU, operand1, operand2 int) {
	cpu.Reg.R[A] = cpu.FetchMemory8(cpu.Reg.HL())
	cpu.Reg.setHL(cpu.Reg.HL() + 1)
	cpu.Reg.PC++
}

// LD A,(HL-)
func op0x3a(cpu *CPU, operand1, operand2 int) {
	cpu.Reg.R[A] = cpu.FetchMemory8(cpu.Reg.HL())
	cpu.Reg.setHL(cpu.Reg.HL() - 1)
	cpu.Reg.PC++
}

// LD A,u8
func op0x3e(cpu *CPU, operand1, operand2 int) {
	cpu.Reg.R[A] = cpu.FetchMemory8(cpu.Reg.PC + 1)
	cpu.Reg.PC += 2
}

// LD A, (HL)
func op0x7e(cpu *CPU, operand1, operand2 int) {
	cpu.Reg.R[A] = cpu.FetchMemory8(cpu.Reg.HL())
	cpu.Reg.PC++
}

// LD A, (u16)
func op0xfa(cpu *CPU, operand1, operand2 int) {
	cpu.Reg.R[A] = cpu.FetchMemory8(cpu.a16FetchJP())
	cpu.Reg.PC += 3
	cpu.timer(2)
}

// LD A,(FF00+C)
func op0xf2(cpu *CPU, operand1, operand2 int) {
	addr := 0xff00 + uint16(cpu.Reg.R[C])
	cpu.Reg.R[A] = cpu.fetchIO(addr)
	cpu.Reg.PC++ // mistake?(https://www.pastraiser.com/cpu/gameboy/gameboy_opcodes.html)
}

// ------ LD B, *

// LD B,u8
func op0x06(cpu *CPU, operand1, operand2 int) {
	cpu.Reg.R[B] = cpu.d8Fetch()
	cpu.Reg.PC += 2
}

// LD B,(HL)
func op0x46(cpu *CPU, operand1, operand2 int) {
	cpu.Reg.R[B] = cpu.FetchMemory8(cpu.Reg.HL())
	cpu.Reg.PC++
}

// ------ LD C, *

// LD C,u8
func op0x0e(cpu *CPU, operand1, operand2 int) {
	value := cpu.d8Fetch()
	cpu.Reg.R[C] = value
	cpu.Reg.PC += 2
}

// LD C,(HL)
func op0x4e(cpu *CPU, operand1, operand2 int) {
	value := cpu.FetchMemory8(cpu.Reg.HL())
	cpu.Reg.R[C] = value
	cpu.Reg.PC++
}

// ------ LD D, *

// LD D,u8
func op0x16(cpu *CPU, operand1, operand2 int) {
	cpu.Reg.R[D] = cpu.d8Fetch()
	cpu.Reg.PC += 2
}

// LD D,(HL)
func op0x56(cpu *CPU, operand1, operand2 int) {
	cpu.Reg.R[D] = cpu.FetchMemory8(cpu.Reg.HL())
	cpu.Reg.PC++
}

// ------ LD E, *

// LD E,u8
func op0x1e(cpu *CPU, operand1, operand2 int) {
	cpu.Reg.R[E] = cpu.d8Fetch()
	cpu.Reg.PC += 2
}

// LD E,(HL)
func op0x5e(cpu *CPU, operand1, operand2 int) {
	cpu.Reg.R[E] = cpu.FetchMemory8(cpu.Reg.HL())
	cpu.Reg.PC++
}

// ------ LD H, *

// LD H,u8
func op0x26(cpu *CPU, operand1, operand2 int) {
	cpu.Reg.R[H] = cpu.d8Fetch()
	cpu.Reg.PC += 2
}

// LD H,(HL)
func op0x66(cpu *CPU, operand1, operand2 int) {
	cpu.Reg.R[H] = cpu.FetchMemory8(cpu.Reg.HL())
	cpu.Reg.PC++
}

// ------ LD L, *

// LD L,u8
func op0x2e(cpu *CPU, operand1, operand2 int) {
	cpu.Reg.R[L] = cpu.d8Fetch()
	cpu.Reg.PC += 2
}

// LD L,(HL)
func op0x6e(cpu *CPU, operand1, operand2 int) {
	cpu.Reg.R[L] = cpu.FetchMemory8(cpu.Reg.HL())
	cpu.Reg.PC++
}

// ------ LD (HL), *

// LD (HL),u8
func op0x36(cpu *CPU, operand1, operand2 int) {
	value := cpu.d8Fetch()
	cpu.timer(1)
	cpu.SetMemory8(cpu.Reg.HL(), value)
	cpu.Reg.PC += 2
	cpu.timer(2)
}

// LD (HL),R8
func ldHLR8(cpu *CPU, unused, op int) {
	cpu.SetMemory8(cpu.Reg.HL(), cpu.Reg.R[op])
	cpu.Reg.PC++
}

// ------ others ld

// LD (u16),SP
func op0x08(cpu *CPU, operand1, operand2 int) {
	// Store SP into addresses n16 (LSB) and n16 + 1 (MSB).
	addr := cpu.a16Fetch()
	upper, lower := byte(cpu.Reg.SP>>8), byte(cpu.Reg.SP) // MSB
	cpu.SetMemory8(addr, lower)
	cpu.SetMemory8(addr+1, upper)
	cpu.Reg.PC += 3
	cpu.timer(5)
}

// LD (u16),A
func op0xea(cpu *CPU, operand1, operand2 int) {
	cpu.SetMemory8(cpu.a16FetchJP(), cpu.Reg.R[A])
	cpu.Reg.PC += 3
	cpu.timer(2)
}

// LD BC,u16
func op0x01(cpu *CPU, operand1, operand2 int) {
	cpu.Reg.setBC(cpu.d16Fetch())
	cpu.Reg.PC += 3
}

// LD DE,u16
func op0x11(cpu *CPU, operand1, operand2 int) {
	cpu.Reg.setDE(cpu.d16Fetch())
	cpu.Reg.PC += 3
}

// LD HL,u16
func op0x21(cpu *CPU, operand1, operand2 int) {
	cpu.Reg.setHL(cpu.d16Fetch())
	cpu.Reg.PC += 3
}

// LD SP,u16
func op0x31(cpu *CPU, operand1, operand2 int) {
	cpu.Reg.SP = cpu.d16Fetch()
	cpu.Reg.PC += 3
}

// LD HL,SP+i8
func op0xf8(cpu *CPU, operand1, operand2 int) {
	delta := int8(cpu.FetchMemory8(cpu.Reg.PC + 1))
	value := int32(cpu.Reg.SP) + int32(delta)
	carryBits := uint32(cpu.Reg.SP) ^ uint32(delta) ^ uint32(value)
	cpu.Reg.setHL(uint16(value))
	cpu.setF(flagZ, false)
	cpu.setF(flagN, false)
	cpu.setF(flagC, util.Bit(carryBits, 8))
	cpu.setF(flagH, util.Bit(carryBits, 4))
	cpu.Reg.PC += 2
}

// LD SP,HL
func op0xf9(cpu *CPU, operand1, operand2 int) {
	cpu.Reg.SP = cpu.Reg.HL()
	cpu.Reg.PC++
}

// LD (FF00+C),A
func op0xe2(cpu *CPU, operand1, operand2 int) {
	addr := 0xff00 + uint16(cpu.Reg.R[C])
	cpu.SetMemory8(addr, cpu.Reg.R[A])
	cpu.Reg.PC++ // mistake?(https://www.pastraiser.com/cpu/gameboy/gameboy_opcodes.html)
}

// LD (BC),A
func op0x02(cpu *CPU, operand1, operand2 int) {
	cpu.SetMemory8(cpu.Reg.BC(), cpu.Reg.R[A])
	cpu.Reg.PC++
}

// LD (DE),A
func op0x12(cpu *CPU, operand1, operand2 int) {
	cpu.SetMemory8(cpu.Reg.DE(), cpu.Reg.R[A])
	cpu.Reg.PC++
}

// LD (HL+),A
func op0x22(cpu *CPU, operand1, operand2 int) {
	cpu.SetMemory8(cpu.Reg.HL(), cpu.Reg.R[A])
	cpu.Reg.setHL(cpu.Reg.HL() + 1)
	cpu.Reg.PC++
}

// LD (HL-),A
func op0x32(cpu *CPU, operand1, operand2 int) {
	// (HL)=A, HL=HL-1
	cpu.SetMemory8(cpu.Reg.HL(), cpu.Reg.R[A])
	cpu.Reg.setHL(cpu.Reg.HL() - 1)
	cpu.Reg.PC++
}

// LD Load
func LD(cpu *CPU, operand1, operand2 int) {
	errMsg := fmt.Sprintf("Error: LD %d %d", operand1, operand2)
	panic(errMsg)
}

// LDH Load High Byte
func LDH(cpu *CPU, operand1, operand2 int) {
	if operand1 == OP_A && operand2 == OP_a8_PAREN { // LD A,($FF00+a8)
		addr := 0xff00 + uint16(cpu.FetchMemory8(cpu.Reg.PC+1))
		cpu.timer(1)
		value := cpu.fetchIO(addr)

		cpu.Reg.R[A] = value
		cpu.Reg.PC += 2
		cpu.timer(2)
	} else if operand1 == OP_a8_PAREN && operand2 == OP_A { // LD ($FF00+a8),A
		addr := 0xff00 + uint16(cpu.FetchMemory8(cpu.Reg.PC+1))
		cpu.timer(1)
		cpu.setIO(addr, cpu.Reg.R[A])
		cpu.Reg.PC += 2
		cpu.timer(2)
	} else {
		panic(fmt.Errorf("error: LDH %d %d", operand1, operand2))
	}
}

// NOP No operation
func (cpu *CPU) NOP(operand1, operand2 int) {
	cpu.Reg.PC++
}

// INC Increment

func incR8(cpu *CPU, op, _ int) {
	value := cpu.Reg.R[op] + 1
	carryBits := cpu.Reg.R[op] ^ 1 ^ value
	cpu.Reg.R[op] = value

	cpu.setF(flagZ, value == 0)
	cpu.setF(flagN, false)
	cpu.setF(flagH, util.Bit(carryBits, 4))
	cpu.Reg.PC++
}

func (cpu *CPU) INC(operand1, operand2 int) {
	var value, carryBits byte

	switch operand1 {
	case OP_HL_PAREN:
		value = cpu.FetchMemory8(cpu.Reg.HL()) + 1
		cpu.timer(1)
		carryBits = cpu.FetchMemory8(cpu.Reg.HL()) ^ 1 ^ value
		cpu.SetMemory8(cpu.Reg.HL(), value)
		cpu.timer(2)
	case OP_BC:
		cpu.Reg.setBC(cpu.Reg.BC() + 1)
	case OP_DE:
		cpu.Reg.setDE(cpu.Reg.DE() + 1)
	case OP_HL:
		cpu.Reg.setHL(cpu.Reg.HL() + 1)
	case OP_SP:
		cpu.Reg.SP++
	default:
		panic(fmt.Errorf("error: INC %d %d", operand1, operand2))
	}

	if operand1 != OP_BC && operand1 != OP_DE && operand1 != OP_HL && operand1 != OP_SP {
		cpu.setF(flagZ, value == 0)
		cpu.setF(flagN, false)
		cpu.setF(flagH, util.Bit(carryBits, 4))
	}
	cpu.Reg.PC++
}

// DEC Decrement

func decR8(cpu *CPU, op, _ int) {
	value := cpu.Reg.R[op] - 1
	carryBits := cpu.Reg.R[op] ^ 1 ^ value
	cpu.Reg.R[op] = value

	cpu.setF(flagZ, value == 0)
	cpu.setF(flagN, true)
	cpu.setF(flagH, util.Bit(carryBits, 4))
	cpu.Reg.PC++
}

func (cpu *CPU) DEC(operand1, operand2 int) {
	var value byte
	var carryBits byte

	switch operand1 {
	case OP_HL_PAREN:
		value = cpu.FetchMemory8(cpu.Reg.HL()) - 1
		cpu.timer(1)
		carryBits = cpu.FetchMemory8(cpu.Reg.HL()) ^ 1 ^ value
		cpu.SetMemory8(cpu.Reg.HL(), value)
		cpu.timer(2)
	case OP_BC:
		cpu.Reg.setBC(cpu.Reg.BC() - 1)
	case OP_DE:
		cpu.Reg.setDE(cpu.Reg.DE() - 1)
	case OP_HL:
		cpu.Reg.setHL(cpu.Reg.HL() - 1)
	case OP_SP:
		cpu.Reg.SP--
	default:
		panic(fmt.Errorf("error: DEC %d %d", operand1, operand2))
	}

	if operand1 != OP_BC && operand1 != OP_DE && operand1 != OP_HL && operand1 != OP_SP {
		cpu.setF(flagZ, value == 0)
		cpu.setF(flagN, true)
		cpu.setF(flagH, util.Bit(carryBits, 4))
	}
	cpu.Reg.PC++
}

// --------- JR ---------

// JR i8
func op0x18(cpu *CPU, operand1, operand2 int) {
	delta := int8(cpu.FetchMemory8(cpu.Reg.PC + 1))
	destination := uint16(int32(cpu.Reg.PC+2) + int32(delta)) // PC+2 because of time after fetch(pc is incremented)
	cpu.Reg.PC = destination
	cpu.timer(3)
}

// JR NZ,i8
func op0x20(cpu *CPU, operand1, operand2 int) {
	if !cpu.f(flagZ) {
		delta := int8(cpu.FetchMemory8(cpu.Reg.PC + 1))
		destination := uint16(int32(cpu.Reg.PC+2) + int32(delta)) // PC+2 because of time after fetch(pc is incremented)
		cpu.Reg.PC = destination
		cpu.timer(3)
	} else {
		cpu.Reg.PC += 2
		cpu.timer(2)
	}
}

// JR Z,i8
func op0x28(cpu *CPU, operand1, operand2 int) {
	if cpu.f(flagZ) {
		delta := int8(cpu.FetchMemory8(cpu.Reg.PC + 1))
		destination := uint16(int32(cpu.Reg.PC+2) + int32(delta)) // PC+2 because of time after fetch(pc is incremented)
		cpu.Reg.PC = destination
		cpu.timer(3)
	} else {
		cpu.Reg.PC += 2
		cpu.timer(2)
	}
}

// JR NC,i8
func op0x30(cpu *CPU, operand1, operand2 int) {
	if !cpu.f(flagC) {
		delta := int8(cpu.FetchMemory8(cpu.Reg.PC + 1))
		destination := uint16(int32(cpu.Reg.PC+2) + int32(delta)) // PC+2 because of time after fetch(pc is incremented)
		cpu.Reg.PC = destination
		cpu.timer(3)
	} else {
		cpu.Reg.PC += 2
		cpu.timer(2)
	}
}

// JR C,i8
func op0x38(cpu *CPU, operand1, operand2 int) {
	if cpu.f(flagC) {
		delta := int8(cpu.FetchMemory8(cpu.Reg.PC + 1))
		destination := uint16(int32(cpu.Reg.PC+2) + int32(delta)) // PC+2 because of time after fetch(pc is incremented)
		cpu.Reg.PC = destination
		cpu.timer(3)
	} else {
		cpu.Reg.PC += 2
		cpu.timer(2)
	}
}

// JR Jump relatively
func JR(cpu *CPU, operand1, operand2 int) {
	result := true
	switch operand1 {
	case OP_r8:
		delta := int8(cpu.FetchMemory8(cpu.Reg.PC + 1))
		destination := uint16(int32(cpu.Reg.PC+2) + int32(delta)) // PC+2 because of time after fetch(pc is incremented)
		cpu.Reg.PC = destination
	case OP_Z:
		if cpu.f(flagZ) {
			delta := int8(cpu.FetchMemory8(cpu.Reg.PC + 1))
			destination := uint16(int32(cpu.Reg.PC+2) + int32(delta)) // PC+2 because of time after fetch(pc is incremented)
			cpu.Reg.PC = destination
		} else {
			cpu.Reg.PC += 2
			result = false
		}
	case OP_C:
		if cpu.f(flagC) {
			delta := int8(cpu.FetchMemory8(cpu.Reg.PC + 1))
			destination := uint16(int32(cpu.Reg.PC+2) + int32(delta)) // PC+2 because of time after fetch(pc is incremented)
			cpu.Reg.PC = destination
		} else {
			cpu.Reg.PC += 2
			result = false
		}
	case OP_NZ:
		if !cpu.f(flagZ) {
			delta := int8(cpu.FetchMemory8(cpu.Reg.PC + 1))
			destination := uint16(int32(cpu.Reg.PC+2) + int32(delta)) // PC+2 because of time after fetch(pc is incremented)
			cpu.Reg.PC = destination
		} else {
			cpu.Reg.PC += 2
			result = false
		}
	case OP_NC:
		if !cpu.f(flagC) {
			delta := int8(cpu.FetchMemory8(cpu.Reg.PC + 1))
			destination := uint16(int32(cpu.Reg.PC+2) + int32(delta)) // PC+2 because of time after fetch(pc is incremented)
			cpu.Reg.PC = destination
		} else {
			cpu.Reg.PC += 2
			result = false
		}
	default:
		panic(fmt.Errorf("error: JR %d %d", operand1, operand2))
	}

	if result {
		cpu.timer(3)
	} else {
		cpu.timer(2)
	}
}

var pending bool

func halt(cpu *CPU, _, _ int) {
	cpu.Reg.PC++
	cpu.halt = true

	// ref: https://rednex.github.io/rgbds/gbz80.7.html#HALT
	if !cpu.Reg.IME {
		IE, IF := cpu.RAM[IEIO], cpu.RAM[IFIO]
		pending = IE&IF != 0
	}
}

func (cpu *CPU) pend() {
	// Some pending
	cpu.halt = false
	PC := cpu.Reg.PC
	cpu.exec()
	cpu.Reg.PC = PC

	// IME turns on due to EI delay.
	cpu.halt = cpu.Reg.IME
}

// STOP stop CPU
func (cpu *CPU) STOP(operand1, operand2 int) {
	if operand1 == OP_0 && operand2 == OP_NONE {
		cpu.Reg.PC += 2
		KEY1 := cpu.FetchMemory8(KEY1IO)
		if util.Bit(KEY1, 0) {
			if util.Bit(KEY1, 7) {
				KEY1 = 0x00
				cpu.boost = 1
			} else {
				KEY1 = 0x80
				cpu.boost = 2
			}
			cpu.SetMemory8(KEY1IO, KEY1)
		}
	} else {
		panic(fmt.Errorf("error: STOP %d %d", operand1, operand2))
	}
}

// XOR xor
func (cpu *CPU) XOR(operand1, operand2 int) {
	var value byte
	switch operand1 {
	case OP_B:
		value = cpu.Reg.R[A] ^ cpu.Reg.R[B]
	case OP_C:
		value = cpu.Reg.R[A] ^ cpu.Reg.R[C]
	case OP_D:
		value = cpu.Reg.R[A] ^ cpu.Reg.R[D]
	case OP_E:
		value = cpu.Reg.R[A] ^ cpu.Reg.R[E]
	case OP_H:
		value = cpu.Reg.R[A] ^ cpu.Reg.R[H]
	case OP_L:
		value = cpu.Reg.R[A] ^ cpu.Reg.R[L]
	case OP_HL_PAREN:
		value = cpu.Reg.R[A] ^ cpu.FetchMemory8(cpu.Reg.HL())
	case OP_A:
		value = 0
	case OP_d8:
		value = cpu.Reg.R[A] ^ cpu.FetchMemory8(cpu.Reg.PC+1)
		cpu.Reg.PC++
	default:
		panic(fmt.Errorf("error: XOR %d %d", operand1, operand2))
	}

	cpu.Reg.R[A] = value
	cpu.setF(flagZ, value == 0)
	cpu.setF(flagN, false)
	cpu.setF(flagH, false)
	cpu.setF(flagC, false)
	cpu.Reg.PC++
}

// JP Jump
func JP(cpu *CPU, operand1, operand2 int) {
	cycle := 1

	switch operand1 {
	case OP_a16:
		destination := cpu.a16FetchJP()
		cycle++
		cpu.Reg.PC = destination
	case OP_HL_PAREN:
		cpu.Reg.PC = cpu.Reg.HL()
	case OP_Z:
		destination := cpu.a16FetchJP()
		if cpu.f(flagZ) {
			cycle++
			cpu.Reg.PC = destination
		} else {
			cpu.Reg.PC += 3
		}
	case OP_C:
		destination := cpu.a16FetchJP()
		if cpu.f(flagC) {
			cycle++
			cpu.Reg.PC = destination
		} else {
			cpu.Reg.PC += 3
		}
	case OP_NZ:
		destination := cpu.a16FetchJP()
		if !cpu.f(flagZ) {
			cycle++
			cpu.Reg.PC = destination
		} else {
			cpu.Reg.PC += 3
		}
	case OP_NC:
		destination := cpu.a16FetchJP()
		if !cpu.f(flagC) {
			cycle++
			cpu.Reg.PC = destination
		} else {
			cpu.Reg.PC += 3
		}
	default:
		panic(fmt.Errorf("error: JP %d %d", operand1, operand2))
	}

	cpu.timer(cycle)
}

// RET Return
func (cpu *CPU) RET(op1, op2 int) (result bool) {
	result = true

	switch op1 {
	case OP_NONE: // PC=(SP), SP=SP+2
		cpu.popPC()
	case OP_Z:
		if cpu.f(flagZ) {
			cpu.popPC()
		} else {
			cpu.Reg.PC++
			result = false
		}
	case OP_C:
		if cpu.f(flagC) {
			cpu.popPC()
		} else {
			cpu.Reg.PC++
			result = false
		}
	case OP_NZ:
		if !cpu.f(flagZ) {
			cpu.popPC()
		} else {
			cpu.Reg.PC++
			result = false
		}
	case OP_NC:
		if !cpu.f(flagC) {
			cpu.popPC()
		} else {
			cpu.Reg.PC++
			result = false
		}
	}

	return result
}

// RETI Return Interrupt
func (cpu *CPU) RETI(operand1, operand2 int) {
	cpu.popPC()
	cpu.Reg.IME = true
}

// CALL Call subroutine
func CALL(cpu *CPU, operand1, operand2 int) {

	switch operand1 {
	case OP_a16:
		destination := cpu.a16FetchJP()
		cpu.Reg.PC += 3
		cpu.timer(1)
		cpu.pushPCCALL()
		cpu.timer(1)
		cpu.Reg.PC = destination
	case OP_Z:
		if cpu.f(flagZ) {
			destination := cpu.a16FetchJP()
			cpu.Reg.PC += 3
			cpu.timer(1)
			cpu.pushPCCALL()
			cpu.timer(1)
			cpu.Reg.PC = destination
		} else {
			cpu.Reg.PC += 3
			cpu.timer(3)
		}
	case OP_C:
		if cpu.f(flagC) {
			destination := cpu.a16FetchJP()
			cpu.Reg.PC += 3
			cpu.timer(1)
			cpu.pushPCCALL()
			cpu.timer(1)
			cpu.Reg.PC = destination
		} else {
			cpu.Reg.PC += 3
			cpu.timer(3)
		}
	case OP_NZ:
		if !cpu.f(flagZ) {
			destination := cpu.a16FetchJP()
			cpu.Reg.PC += 3
			cpu.timer(1)
			cpu.pushPCCALL()
			cpu.timer(1)
			cpu.Reg.PC = destination
		} else {
			cpu.Reg.PC += 3
			cpu.timer(3)
		}
	case OP_NC:
		if !cpu.f(flagC) {
			destination := cpu.a16FetchJP()
			cpu.Reg.PC += 3
			cpu.timer(1)
			cpu.pushPCCALL()
			cpu.timer(1)
			cpu.Reg.PC = destination
		} else {
			cpu.Reg.PC += 3
			cpu.timer(3)
		}
	default:
		panic(fmt.Errorf("error: CALL %d %d", operand1, operand2))
	}
}

// DI Disable Interrupt
func di(cpu *CPU, _, _ int) {
	cpu.Reg.IME = false
	cpu.Reg.PC++
	if cpu.IMESwitch.Working && cpu.IMESwitch.Value {
		cpu.IMESwitch.Working = false // https://gbdev.gg8.se/wiki/articles/Interrupts 『The effect of EI is delayed by one instruction. This means that EI followed immediately by DI does not allow interrupts between the EI and the DI.』
	}
}

// EI Enable Interrupt
func ei(cpu *CPU, _, _ int) {
	// ref: https://github.com/Gekkio/mooneye-gb/blob/master/tests/acceptance/halt_ime0_ei.s#L23
	next := cpu.FetchMemory8(cpu.Reg.PC + 1) // next opcode
	HALT := byte(0x76)
	if next == HALT {
		cpu.Reg.IME = true
		cpu.Reg.PC++
		return
	}

	if !cpu.IMESwitch.Working {
		cpu.IMESwitch = IMESwitch{
			Count:   2,
			Value:   true,
			Working: true,
		}
	}
	cpu.Reg.PC++
}

// CP Compare
func cp(cpu *CPU, _, op int) {
	var value, carryBits byte
	value = cpu.Reg.R[A] - cpu.Reg.R[op]
	carryBits = cpu.Reg.R[A] ^ cpu.Reg.R[op] ^ value
	cpu.setCSub(cpu.Reg.R[A], cpu.Reg.R[op])

	cpu.setF(flagZ, value == 0)
	cpu.setF(flagN, true)
	cpu.setF(flagH, util.Bit(carryBits, 4))
	cpu.Reg.PC++
}

func (cpu *CPU) CP(operand1, operand2 int) {
	var value, carryBits byte
	switch operand1 {
	case OP_d8:
		value = cpu.Reg.R[A] - cpu.d8Fetch()
		carryBits = cpu.Reg.R[A] ^ cpu.d8Fetch() ^ value
		cpu.setCSub(cpu.Reg.R[A], cpu.d8Fetch())
		cpu.Reg.PC++
	case OP_HL_PAREN:
		value = cpu.Reg.R[A] - cpu.FetchMemory8(cpu.Reg.HL())
		carryBits = cpu.Reg.R[A] ^ cpu.FetchMemory8(cpu.Reg.HL()) ^ value
		cpu.setCSub(cpu.Reg.R[A], cpu.FetchMemory8(cpu.Reg.HL()))
	}
	cpu.setF(flagZ, value == 0)
	cpu.setF(flagN, true)
	cpu.setF(flagH, util.Bit(carryBits, 4))
	cpu.Reg.PC++
}

// AND And instruction

func andR8(cpu *CPU, _, op int) {
	value := cpu.Reg.R[A] & cpu.Reg.R[op]
	cpu.Reg.R[A] = value
	cpu.setF(flagZ, value == 0)
	cpu.setF(flagN, false)
	cpu.setF(flagH, true)
	cpu.setF(flagC, false)
	cpu.Reg.PC++
}

func (cpu *CPU) AND(operand1, operand2 int) {
	var value byte
	switch operand1 {
	case OP_HL_PAREN:
		value = cpu.Reg.R[A] & cpu.FetchMemory8(cpu.Reg.HL())
	case OP_d8:
		value = cpu.Reg.R[A] & cpu.d8Fetch()
		cpu.Reg.PC++
	}

	cpu.Reg.R[A] = value
	cpu.setF(flagZ, value == 0)
	cpu.setF(flagN, false)
	cpu.setF(flagH, true)
	cpu.setF(flagC, false)
	cpu.Reg.PC++
}

// OR or
func orR8(cpu *CPU, _, op int) {
	value := cpu.Reg.R[A] | cpu.Reg.R[op]
	cpu.Reg.R[A] = value
	cpu.setF(flagZ, value == 0)

	cpu.setF(flagN, false)
	cpu.setF(flagH, false)
	cpu.setF(flagC, false)
	cpu.Reg.PC++
}

func (cpu *CPU) OR(operand1, operand2 int) {
	switch operand1 {
	case OP_d8:
		value := cpu.Reg.R[A] | cpu.FetchMemory8(cpu.Reg.PC+1)
		cpu.Reg.R[A] = value
		cpu.setF(flagZ, value == 0)
		cpu.Reg.PC++
	case OP_HL_PAREN:
		value := cpu.Reg.R[A] | cpu.FetchMemory8(cpu.Reg.HL())
		cpu.Reg.R[A] = value
		cpu.setF(flagZ, value == 0)
	}

	cpu.setF(flagN, false)
	cpu.setF(flagH, false)
	cpu.setF(flagC, false)
	cpu.Reg.PC++
}

// ADD Addition
func addR8(cpu *CPU, _, op int) {
	value := uint16(cpu.Reg.R[A]) + uint16(cpu.Reg.R[op])
	carryBits := uint16(cpu.Reg.R[A]) ^ uint16(cpu.Reg.R[op]) ^ value
	cpu.Reg.R[A] = byte(value)
	cpu.setF(flagZ, byte(value) == 0)
	cpu.setF(flagN, false)
	cpu.setF(flagH, util.Bit(carryBits, 4))
	cpu.setF(flagC, util.Bit(carryBits, 8))
	cpu.Reg.PC++
}

func (cpu *CPU) ADD(operand1, operand2 int) {
	switch operand1 {
	case OP_A:
		switch operand2 {
		case OP_d8:
			value := uint16(cpu.Reg.R[A]) + uint16(cpu.d8Fetch())
			carryBits := uint16(cpu.Reg.R[A]) ^ uint16(cpu.d8Fetch()) ^ value
			cpu.Reg.R[A] = byte(value)
			cpu.setF(flagZ, byte(value) == 0)
			cpu.setF(flagN, false)
			cpu.setF(flagH, util.Bit(carryBits, 4))
			cpu.setF(flagC, util.Bit(carryBits, 8))
			cpu.Reg.PC += 2
		case OP_HL_PAREN:
			value := uint16(cpu.Reg.R[A]) + uint16(cpu.FetchMemory8(cpu.Reg.HL()))
			carryBits := uint16(cpu.Reg.R[A]) ^ uint16(cpu.FetchMemory8(cpu.Reg.HL())) ^ value
			cpu.Reg.R[A] = byte(value)
			cpu.setF(flagZ, byte(value) == 0)
			cpu.setF(flagN, false)
			cpu.setF(flagH, util.Bit(carryBits, 4))
			cpu.setF(flagC, util.Bit(carryBits, 8))
			cpu.Reg.PC++
		}
	case OP_HL:
		switch operand2 {
		case OP_BC:
			value := uint32(cpu.Reg.HL()) + uint32(cpu.Reg.BC())
			carryBits := uint32(cpu.Reg.HL()) ^ uint32(cpu.Reg.BC()) ^ value
			cpu.Reg.setHL(uint16(value))
			cpu.setF(flagN, false)
			cpu.setF(flagH, util.Bit(carryBits, 12))
			cpu.setF(flagC, util.Bit(carryBits, 16))
			cpu.Reg.PC++
		case OP_DE:
			value := uint32(cpu.Reg.HL()) + uint32(cpu.Reg.DE())
			carryBits := uint32(cpu.Reg.HL()) ^ uint32(cpu.Reg.DE()) ^ value
			cpu.Reg.setHL(uint16(value))
			cpu.setF(flagN, false)
			cpu.setF(flagH, util.Bit(carryBits, 12))
			cpu.setF(flagC, util.Bit(carryBits, 16))
			cpu.Reg.PC++
		case OP_HL:
			value := uint32(cpu.Reg.HL()) + uint32(cpu.Reg.HL())
			carryBits := value
			cpu.Reg.setHL(uint16(value))
			cpu.setF(flagN, false)
			cpu.setF(flagH, util.Bit(carryBits, 12))
			cpu.setF(flagC, util.Bit(carryBits, 16))
			cpu.Reg.PC++
		case OP_SP:
			value := uint32(cpu.Reg.HL()) + uint32(cpu.Reg.SP)
			carryBits := uint32(cpu.Reg.HL()) ^ uint32(cpu.Reg.SP) ^ value
			cpu.Reg.setHL(uint16(value))
			cpu.setF(flagN, false)
			cpu.setF(flagH, util.Bit(carryBits, 12))
			cpu.setF(flagC, util.Bit(carryBits, 16))
			cpu.Reg.PC++
		}
	case OP_SP:
		switch operand2 {
		case OP_r8:
			delta := int8(cpu.FetchMemory8(cpu.Reg.PC + 1))
			value := int32(cpu.Reg.SP) + int32(delta)
			carryBits := uint32(cpu.Reg.SP) ^ uint32(delta) ^ uint32(value)
			cpu.Reg.SP = uint16(value)
			cpu.setF(flagZ, false)
			cpu.setF(flagN, false)
			cpu.setF(flagH, util.Bit(carryBits, 4))
			cpu.setF(flagC, util.Bit(carryBits, 8))
			cpu.Reg.PC += 2
		}
	}
}

// CPL Complement A Register
func (cpu *CPU) CPL(operand1, operand2 int) {
	cpu.Reg.R[A] = ^cpu.Reg.R[A]
	cpu.setF(flagN, true)
	cpu.setF(flagH, true)
	cpu.Reg.PC++
}

// PREFIXCB is extend instruction
func (cpu *CPU) PREFIXCB(op1, op2 int) {
	if op1 == OP_NONE && op2 == OP_NONE {
		cpu.Reg.PC++
		cpu.timer(1)
		op := prefixCBs[cpu.FetchMemory8(cpu.Reg.PC)]
		instruction, op1, op2, cycle, handler := op.Ins, op.Operand1, op.Operand2, op.Cycle1, op.Handler

		if handler != nil {
			handler(cpu, op1, op2)
		} else {
			switch instruction {
			case INS_RRC:
				cpu.RRC(op1, op2)
			case INS_SRA:
				cpu.SRA(op1, op2)
			case INS_SRL:
				cpu.SRL(op1, op2)
			case INS_BIT:
				cpu.BIT(op1, op2)
			default:
				panic(fmt.Errorf("eip: 0x%04x opcode: %v", cpu.Reg.PC, op))
			}
		}

		if cycle > 1 {
			cpu.timer(cycle - 1)
		}
	} else {
		panic(fmt.Errorf("error: PREFIXCB %d %d", op1, op2))
	}
}

// RLC Rotate n left carry => bit0
func rlcR8(cpu *CPU, op, _ int) {
	value := cpu.Reg.R[op]
	bit7 := value >> 7
	value = (value << 1)
	value = util.SetLSB(value, bit7 != 0)
	cpu.Reg.R[op] = value

	cpu.setF(flagZ, value == 0)
	cpu.setF(flagN, false)
	cpu.setF(flagH, false)
	cpu.setF(flagC, bit7 != 0)
	cpu.Reg.PC++
}

func rlcHL(cpu *CPU, _, _ int) {
	value := cpu.FetchMemory8(cpu.Reg.HL())
	cpu.timer(1)
	bit7 := value >> 7
	value = (value << 1)
	value = util.SetLSB(value, bit7 != 0)
	cpu.SetMemory8(cpu.Reg.HL(), value)
	cpu.timer(2)

	cpu.setF(flagZ, value == 0)
	cpu.setF(flagN, false)
	cpu.setF(flagH, false)
	cpu.setF(flagC, bit7 != 0)
	cpu.Reg.PC++
}

// RLCA Rotate register A left.
func (cpu *CPU) RLCA(operand1, operand2 int) {
	var value byte
	var bit7 byte
	value = cpu.Reg.R[A]
	bit7 = value >> 7
	value = (value << 1)
	value = util.SetLSB(value, bit7 != 0)
	cpu.Reg.R[A] = value

	cpu.setF(flagZ, false)
	cpu.setF(flagN, false)
	cpu.setF(flagH, false)
	cpu.setF(flagC, bit7 != 0)
	cpu.Reg.PC++
}

// RRC Rotate n right carry => bit7
func rrcR8(cpu *CPU, op, _ int) {
	value := cpu.Reg.R[op]
	bit0 := value % 2
	value = (value >> 1)
	value = util.SetMSB(value, bit0 != 0)
	cpu.Reg.R[op] = value

	cpu.setF(flagZ, value == 0)
	cpu.setF(flagN, false)
	cpu.setF(flagH, false)
	cpu.setF(flagC, bit0 != 0)
	cpu.Reg.PC++
}

func (cpu *CPU) RRC(operand1, operand2 int) {
	value := cpu.FetchMemory8(cpu.Reg.HL())
	cpu.timer(1)
	bit0 := value % 2
	value = (value >> 1)
	value = util.SetMSB(value, bit0 != 0)
	cpu.SetMemory8(cpu.Reg.HL(), value)
	cpu.timer(2)

	cpu.setF(flagZ, value == 0)
	cpu.setF(flagN, false)
	cpu.setF(flagH, false)
	cpu.setF(flagC, bit0 != 0)
	cpu.Reg.PC++
}

// RRCA Rotate register A right.
func (cpu *CPU) RRCA(operand1, operand2 int) {
	var value byte
	var lsb bool

	value, lsb = cpu.Reg.R[A], util.Bit(cpu.Reg.R[A], 0)
	value = (value >> 1)
	value = util.SetMSB(value, lsb)
	cpu.Reg.R[A] = value

	cpu.setF(flagZ, false)
	cpu.setF(flagN, false)
	cpu.setF(flagH, false)
	cpu.setF(flagC, lsb)
	cpu.Reg.PC++
}

// RL Rotate n rigth through carry bit7 => bit0
func rl(cpu *CPU, _, op int) {
	carry, value := cpu.f(flagC), cpu.Reg.R[op]
	bit7 := value >> 7
	value = (value << 1)
	value = util.SetLSB(value, carry)
	cpu.Reg.R[op] = value

	cpu.setF(flagZ, value == 0)
	cpu.setF(flagN, false)
	cpu.setF(flagH, false)
	cpu.setF(flagC, bit7 != 0)
	cpu.Reg.PC++
}

func rlHL(cpu *CPU, _, _ int) {
	var value, bit7 byte
	carry := cpu.f(flagC)
	value = cpu.FetchMemory8(cpu.Reg.HL())
	cpu.timer(1)
	bit7 = value >> 7
	value = (value << 1)
	value = util.SetLSB(value, carry)
	cpu.SetMemory8(cpu.Reg.HL(), value)
	cpu.timer(2)

	cpu.setF(flagZ, value == 0)
	cpu.setF(flagN, false)
	cpu.setF(flagH, false)
	cpu.setF(flagC, bit7 != 0)
	cpu.Reg.PC++
}

// RLA Rotate register A left through carry.
func (cpu *CPU) RLA(operand1, operand2 int) {
	var value, bit7 byte
	carry := cpu.f(flagC)

	value = cpu.Reg.R[A]
	bit7 = value >> 7
	value = (value << 1)
	value = util.SetLSB(value, carry)
	cpu.Reg.R[A] = value

	cpu.setF(flagZ, false)
	cpu.setF(flagN, false)
	cpu.setF(flagH, false)
	cpu.setF(flagC, bit7 != 0)
	cpu.Reg.PC++
}

// rr Rotate n right through carry bit0 => bit7
func rr(cpu *CPU, op, _ int) {
	value, lsb, carry := cpu.Reg.R[op], util.Bit(cpu.Reg.R[op], 0), cpu.f(flagC)
	value >>= 1
	value = util.SetMSB(value, carry)
	cpu.Reg.R[op] = value

	cpu.setF(flagZ, value == 0)
	cpu.setF(flagN, false)
	cpu.setF(flagH, false)
	cpu.setF(flagC, lsb)
	cpu.Reg.PC++
}

func rrHL(cpu *CPU, _, _ int) {
	carry := cpu.f(flagC)
	value := cpu.FetchMemory8(cpu.Reg.HL())
	cpu.timer(1)
	lsb := util.Bit(value, 0)
	value >>= 1
	value = util.SetMSB(value, carry)
	cpu.SetMemory8(cpu.Reg.HL(), value)
	cpu.timer(2)

	cpu.setF(flagZ, value == 0)
	cpu.setF(flagN, false)
	cpu.setF(flagH, false)
	cpu.setF(flagC, lsb)
	cpu.Reg.PC++
}

// Shift Left
func sla(cpu *CPU, op, _ int) {
	value := cpu.Reg.R[op]
	bit7 := value >> 7
	value = (value << 1)
	cpu.Reg.R[op] = value

	cpu.setF(flagZ, value == 0)
	cpu.setF(flagN, false)
	cpu.setF(flagH, false)
	cpu.setF(flagC, bit7 != 0)
	cpu.Reg.PC++
}

func slaHL(cpu *CPU, _, _ int) {
	value := cpu.FetchMemory8(cpu.Reg.HL())
	cpu.timer(1)
	bit7 := value >> 7
	value = (value << 1)
	cpu.SetMemory8(cpu.Reg.HL(), value)
	cpu.timer(2)

	cpu.setF(flagZ, value == 0)
	cpu.setF(flagN, false)
	cpu.setF(flagH, false)
	cpu.setF(flagC, bit7 != 0)
	cpu.Reg.PC++
}

// SRA Shift Right MSBit dosen't change
func (cpu *CPU) SRA(operand1, operand2 int) {
	var value byte
	var lsb, msb bool
	if operand1 == OP_B && operand2 == OP_NONE {
		value, lsb, msb = cpu.Reg.R[B], util.Bit(cpu.Reg.R[B], 0), util.Bit(cpu.Reg.R[B], 7)
		value = (value >> 1)
		value = util.SetMSB(value, msb)
		cpu.Reg.R[B] = value
	} else if operand1 == OP_C && operand2 == OP_NONE {
		value, lsb, msb = cpu.Reg.R[C], util.Bit(cpu.Reg.R[C], 0), util.Bit(cpu.Reg.R[C], 7)
		value = (value >> 1)
		value = util.SetMSB(value, msb)
		cpu.Reg.R[C] = value
	} else if operand1 == OP_D && operand2 == OP_NONE {
		value, lsb, msb = cpu.Reg.R[D], util.Bit(cpu.Reg.R[D], 0), util.Bit(cpu.Reg.R[D], 7)
		value = (value >> 1)
		value = util.SetMSB(value, msb)
		cpu.Reg.R[D] = value
	} else if operand1 == OP_E && operand2 == OP_NONE {
		value, lsb, msb = cpu.Reg.R[E], util.Bit(cpu.Reg.R[E], 0), util.Bit(cpu.Reg.R[E], 7)
		value = (value >> 1)
		value = util.SetMSB(value, msb)
		cpu.Reg.R[E] = value
	} else if operand1 == OP_H && operand2 == OP_NONE {
		value, lsb, msb = cpu.Reg.R[H], util.Bit(cpu.Reg.R[H], 0), util.Bit(cpu.Reg.R[H], 7)
		value = (value >> 1)
		value = util.SetMSB(value, msb)
		cpu.Reg.R[H] = value
	} else if operand1 == OP_L && operand2 == OP_NONE {
		value, lsb, msb = cpu.Reg.R[L], util.Bit(cpu.Reg.R[L], 0), util.Bit(cpu.Reg.R[L], 7)
		value = (value >> 1)
		value = util.SetMSB(value, msb)
		cpu.Reg.R[L] = value
	} else if operand1 == OP_HL_PAREN && operand2 == OP_NONE {
		value = cpu.FetchMemory8(cpu.Reg.HL())
		cpu.timer(1)
		lsb, msb = util.Bit(value, 0), util.Bit(value, 7)
		value = (value >> 1)
		value = util.SetMSB(value, msb)
		cpu.SetMemory8(cpu.Reg.HL(), value)
		cpu.timer(2)
	} else if operand1 == OP_A && operand2 == OP_NONE {
		value, lsb, msb = cpu.Reg.R[A], util.Bit(cpu.Reg.R[A], 0), util.Bit(cpu.Reg.R[A], 7)
		value = (value >> 1)
		value = util.SetMSB(value, msb)
		cpu.Reg.R[A] = value
	} else {
		panic(fmt.Errorf("error: SRA %d %d", operand1, operand2))
	}

	cpu.setF(flagZ, value == 0)
	cpu.setF(flagN, false)
	cpu.setF(flagH, false)
	cpu.setF(flagC, lsb)
	cpu.Reg.PC++
}

// SWAP Swap n[5:8] and n[0:4]
func swap(cpu *CPU, _, op int) {
	b := cpu.Reg.R[op]
	lower := b & 0b1111
	upper := b >> 4
	cpu.Reg.R[op] = (lower << 4) | upper

	cpu.setF(flagZ, cpu.Reg.R[op] == 0)
	cpu.setF(flagN, false)
	cpu.setF(flagH, false)
	cpu.setF(flagC, false)
	cpu.Reg.PC++
}

func swapHL(cpu *CPU, _, _ int) {
	data := cpu.FetchMemory8(cpu.Reg.HL())
	cpu.timer(1)
	data03 := data & 0x0f
	data47 := data >> 4
	value := (data03 << 4) | data47
	cpu.SetMemory8(cpu.Reg.HL(), value)
	cpu.timer(2)

	cpu.setF(flagZ, value == 0)
	cpu.setF(flagN, false)
	cpu.setF(flagH, false)
	cpu.setF(flagC, false)
	cpu.Reg.PC++
}

// SRL Shift Right MSBit = 0
func (cpu *CPU) SRL(operand1, operand2 int) {
	var value byte
	var bit0 byte

	switch operand1 {
	case OP_B:
		value = cpu.Reg.R[B]
		bit0 = value % 2
		value = (value >> 1)
		cpu.Reg.R[B] = value
	case OP_C:
		value = cpu.Reg.R[C]
		bit0 = value % 2
		value = (value >> 1)
		cpu.Reg.R[C] = value
	case OP_D:
		value = cpu.Reg.R[D]
		bit0 = value % 2
		value = (value >> 1)
		cpu.Reg.R[D] = value
	case OP_E:
		value = cpu.Reg.R[E]
		bit0 = value % 2
		value = (value >> 1)
		cpu.Reg.R[E] = value
	case OP_H:
		value = cpu.Reg.R[H]
		bit0 = value % 2
		value = (value >> 1)
		cpu.Reg.R[H] = value
	case OP_L:
		value = cpu.Reg.R[L]
		bit0 = value % 2
		value = (value >> 1)
		cpu.Reg.R[L] = value
	case OP_HL_PAREN:
		value = cpu.FetchMemory8(cpu.Reg.HL())
		cpu.timer(1)
		bit0 = value % 2
		value = (value >> 1)
		cpu.SetMemory8(cpu.Reg.HL(), value)
		cpu.timer(2)
	case OP_A:
		value = cpu.Reg.R[A]
		bit0 = value % 2
		value = (value >> 1)
		cpu.Reg.R[A] = value
	default:
		panic(fmt.Errorf("error: SRL %d %d", operand1, operand2))
	}

	cpu.setF(flagZ, value == 0)
	cpu.setF(flagN, false)
	cpu.setF(flagH, false)
	cpu.setF(flagC, bit0 == 1)
	cpu.Reg.PC++
}

// BIT Test bit n
func (cpu *CPU) BIT(operand1, operand2 int) {
	var value bool
	targetBit := operand1 - OP_0
	switch operand2 {
	case OP_B:
		value = util.Bit(cpu.Reg.R[B], targetBit)
	case OP_C:
		value = util.Bit(cpu.Reg.R[C], targetBit)
	case OP_D:
		value = util.Bit(cpu.Reg.R[D], targetBit)
	case OP_E:
		value = util.Bit(cpu.Reg.R[E], targetBit)
	case OP_H:
		value = util.Bit(cpu.Reg.R[H], targetBit)
	case OP_L:
		value = util.Bit(cpu.Reg.R[L], targetBit)
	case OP_HL_PAREN:
		value = util.Bit(cpu.FetchMemory8(cpu.Reg.HL()), targetBit)
	case OP_A:
		value = util.Bit(cpu.Reg.R[A], targetBit)
	}

	cpu.setF(flagZ, !value)
	cpu.setF(flagN, false)
	cpu.setF(flagH, true)
	cpu.Reg.PC++
}

func res(cpu *CPU, bit, r8 int) {
	mask := ^(byte(1) << bit)
	cpu.Reg.R[r8] &= mask
	cpu.Reg.PC++
}

func resHL(cpu *CPU, bit, _ int) {
	mask := ^(byte(1) << bit)
	value := cpu.FetchMemory8(cpu.Reg.HL()) & mask
	cpu.timer(1)
	cpu.SetMemory8(cpu.Reg.HL(), value)
	cpu.timer(2)
	cpu.Reg.PC++
}

func set(cpu *CPU, bit, r8 int) {
	mask := byte(1) << bit
	cpu.Reg.R[r8] |= mask
	cpu.Reg.PC++
}

func setHL(cpu *CPU, bit, _ int) {
	mask := byte(1) << bit
	value := cpu.FetchMemory8(cpu.Reg.HL()) | mask
	cpu.timer(1)
	cpu.SetMemory8(cpu.Reg.HL(), value)
	cpu.timer(2)
	cpu.Reg.PC++
}

// PUSH value
func (cpu *CPU) PUSH(operand1, operand2 int) {
	cpu.timer(1)
	switch operand1 {
	case OP_BC:
		cpu.pushBC()
	case OP_DE:
		cpu.pushDE()
	case OP_HL:
		cpu.pushHL()
	case OP_AF:
		cpu.pushAF()
	default:
		panic(fmt.Errorf("error: PUSH %d %d", operand1, operand2))
	}
	cpu.Reg.PC++
	cpu.timer(2)
}

// POP value
func (cpu *CPU) POP(operand1, operand2 int) {
	switch operand1 {
	case OP_BC:
		cpu.popBC()
	case OP_DE:
		cpu.popDE()
	case OP_HL:
		cpu.popHL()
	case OP_AF:
		cpu.popAF()
	default:
		panic(fmt.Errorf("error: POP %d %d", operand1, operand2))
	}
	cpu.Reg.PC++
	cpu.timer(2)
}

// SUB subtract
func subR8(cpu *CPU, _, op int) {
	value := cpu.Reg.R[A] - cpu.Reg.R[op]
	carryBits := cpu.Reg.R[A] ^ cpu.Reg.R[op] ^ value
	cpu.setCSub(cpu.Reg.R[A], cpu.Reg.R[op])
	cpu.Reg.R[A] = value
	cpu.setF(flagZ, value == 0)
	cpu.setF(flagN, true)
	cpu.setF(flagH, util.Bit(carryBits, 4))
	cpu.Reg.PC++
}

func (cpu *CPU) SUB(op1, _ int) {
	switch op1 {
	case OP_d8:
		value := cpu.Reg.R[A] - cpu.d8Fetch()
		carryBits := cpu.Reg.R[A] ^ cpu.d8Fetch() ^ value
		cpu.setCSub(cpu.Reg.R[A], cpu.d8Fetch())
		cpu.Reg.R[A] = value
		cpu.setF(flagZ, value == 0)
		cpu.setF(flagN, true)
		cpu.setF(flagH, util.Bit(carryBits, 4))
		cpu.Reg.PC += 2
	case OP_HL_PAREN:
		value := cpu.Reg.R[A] - cpu.FetchMemory8(cpu.Reg.HL())
		carryBits := cpu.Reg.R[A] ^ cpu.FetchMemory8(cpu.Reg.HL()) ^ value
		cpu.setCSub(cpu.Reg.R[A], cpu.FetchMemory8(cpu.Reg.HL()))
		cpu.Reg.R[A] = value
		cpu.setF(flagZ, value == 0)
		cpu.setF(flagN, true)
		cpu.setF(flagH, util.Bit(carryBits, 4))
		cpu.Reg.PC++
	}
}

// RRA Rotate register A right through carry.
func (cpu *CPU) RRA(operand1, operand2 int) {
	carry := cpu.f(flagC)
	regA := cpu.Reg.R[A]
	cpu.setF(flagC, util.Bit(regA, 0))
	if carry {
		regA = (1 << 7) | (regA >> 1)
	} else {
		regA = (0 << 7) | (regA >> 1)
	}
	cpu.Reg.R[A] = regA
	cpu.setF(flagZ, false)
	cpu.setF(flagN, false)
	cpu.setF(flagH, false)
	cpu.Reg.PC++
}

// ADC Add the value n8 plus the carry flag to A
func adcAR8(cpu *CPU, _, op int) {
	var carry, value, value4 byte
	var value16 uint16
	if cpu.f(flagC) {
		carry = 1
	}

	value = cpu.Reg.R[op] + carry + cpu.Reg.R[A]
	value4 = (cpu.Reg.R[op] & 0b1111) + carry + (cpu.Reg.R[A] & 0b1111)
	value16 = uint16(cpu.Reg.R[op]) + uint16(carry) + uint16(cpu.Reg.R[A])
	cpu.Reg.R[A] = value

	cpu.setF(flagZ, value == 0)
	cpu.setF(flagN, false)
	cpu.setF(flagH, util.Bit(value4, 4))
	cpu.setF(flagC, util.Bit(value16, 8))
	cpu.Reg.PC++
}

func (cpu *CPU) ADC(_, op2 int) {
	var carry, value, value4 byte
	var value16 uint16
	if cpu.f(flagC) {
		carry = 1
	}

	switch op2 {
	case OP_HL_PAREN:
		data := cpu.FetchMemory8(cpu.Reg.HL())
		value = data + carry + cpu.Reg.R[A]
		value4 = (data & 0x0f) + carry + (cpu.Reg.R[A] & 0b1111)
		value16 = uint16(data) + uint16(cpu.Reg.R[A]) + uint16(carry)
	case OP_d8:
		data := cpu.d8Fetch()
		value = data + carry + cpu.Reg.R[A]
		value4 = (data & 0x0f) + carry + (cpu.Reg.R[A] & 0b1111)
		value16 = uint16(data) + uint16(cpu.Reg.R[A]) + uint16(carry)
		cpu.Reg.PC++
	}
	cpu.Reg.R[A] = value
	cpu.setF(flagZ, value == 0)
	cpu.setF(flagN, false)
	cpu.setF(flagH, util.Bit(value4, 4))
	cpu.setF(flagC, util.Bit(value16, 8))
	cpu.Reg.PC++
}

// SBC Subtract the value n8 and the carry flag from A

func sbcAR8(cpu *CPU, _, op int) {
	var carry, value, value4 byte
	var value16 uint16
	if cpu.f(flagC) {
		carry = 1
	}

	value = cpu.Reg.R[A] - (cpu.Reg.R[op] + carry)
	value4 = (cpu.Reg.R[A] & 0b1111) - ((cpu.Reg.R[op] & 0b1111) + carry)
	value16 = uint16(cpu.Reg.R[A]) - (uint16(cpu.Reg.R[op]) + uint16(carry))
	cpu.Reg.R[A] = value

	cpu.setF(flagZ, value == 0)
	cpu.setF(flagN, true)
	cpu.setF(flagH, util.Bit(value4, 4))
	cpu.setF(flagC, util.Bit(value16, 8))
	cpu.Reg.PC++
}

func (cpu *CPU) SBC(_, op2 int) {
	var carry, value, value4 byte
	var value16 uint16
	if cpu.f(flagC) {
		carry = 1
	}

	switch op2 {
	case OP_HL_PAREN:
		data := cpu.FetchMemory8(cpu.Reg.HL())
		value = cpu.Reg.R[A] - (data + carry)
		value4 = (cpu.Reg.R[A] & 0b1111) - ((data & 0x0f) + carry)
		value16 = uint16(cpu.Reg.R[A]) - (uint16(data) + uint16(carry))
	case OP_d8:
		data := cpu.d8Fetch()
		value = cpu.Reg.R[A] - (data + carry)
		value4 = (cpu.Reg.R[A] & 0b1111) - ((data & 0x0f) + carry)
		value16 = uint16(cpu.Reg.R[A]) - (uint16(data) + uint16(carry))
		cpu.Reg.PC++
	}
	cpu.Reg.R[A] = value

	cpu.setF(flagZ, value == 0)
	cpu.setF(flagN, true)
	cpu.setF(flagH, util.Bit(value4, 4))
	cpu.setF(flagC, util.Bit(value16, 8))
	cpu.Reg.PC++
}

// DAA Decimal adjust
func (cpu *CPU) DAA(operand1, operand2 int) {
	a := uint8(cpu.Reg.R[A])
	// ref: https://forums.nesdev.com/viewtopic.php?f=20&t=15944
	if !cpu.f(flagN) {
		if cpu.f(flagC) || a > 0x99 {
			a += 0x60
			cpu.setF(flagC, true)
		}
		if cpu.f(flagH) || (a&0x0f) > 0x09 {
			a += 0x06
		}
	} else {
		if cpu.f(flagC) {
			a -= 0x60
		}
		if cpu.f(flagH) {
			a -= 0x06
		}
	}

	cpu.Reg.R[A] = a
	cpu.setF(flagZ, a == 0)
	cpu.setF(flagH, false)
	cpu.Reg.PC++
}

// RST Push present address and jump to vector address
func (cpu *CPU) RST(operand1, operand2 int) {
	destination := uint16(operand1)
	cpu.Reg.PC++
	cpu.pushPC()
	cpu.Reg.PC = destination
}

func scf(cpu *CPU, _, _ int) {
	cpu.setF(flagN, false)
	cpu.setF(flagH, false)
	cpu.setF(flagC, true)
	cpu.Reg.PC++
}

// CCF Complement Carry Flag
func ccf(cpu *CPU, _, _ int) {
	cpu.setF(flagN, false)
	cpu.setF(flagH, false)
	cpu.setF(flagC, !cpu.f(flagC))
	cpu.Reg.PC++
}
