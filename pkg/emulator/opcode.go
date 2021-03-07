package emulator

import (
	"fmt"
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
	lower := uint16(cpu.FetchMemory8(cpu.Reg.PC + 1))
	upper := uint16(cpu.FetchMemory8(cpu.Reg.PC + 2))
	value := (upper << 8) | lower
	return value
}

// ------ LD A, *

// LD A,(BC)
func op0x0a(cpu *CPU, operand1, operand2 int) {
	cpu.setAReg(cpu.FetchMemory8(cpu.Reg.BC))
	cpu.Reg.PC++
}

// LD A,(DE)
func op0x1a(cpu *CPU, operand1, operand2 int) {
	cpu.setAReg(cpu.FetchMemory8(cpu.Reg.DE))
	cpu.Reg.PC++
}

// LD A,(HL+)
func op0x2a(cpu *CPU, operand1, operand2 int) {
	cpu.setAReg(cpu.FetchMemory8(cpu.Reg.HL))
	cpu.Reg.HL++
	cpu.Reg.PC++
}

// LD A,(HL-)
func op0x3a(cpu *CPU, operand1, operand2 int) {
	cpu.setAReg(cpu.FetchMemory8(cpu.Reg.HL))
	cpu.Reg.HL--
	cpu.Reg.PC++
}

// LD A,u8
func op0x3e(cpu *CPU, operand1, operand2 int) {
	cpu.setAReg(cpu.FetchMemory8(cpu.Reg.PC + 1))
	cpu.Reg.PC += 2
}

// LD A, B
func op0x78(cpu *CPU, operand1, operand2 int) {
	cpu.setAReg(cpu.getBReg())
	cpu.Reg.PC++
}

// LD A, C
func op0x79(cpu *CPU, operand1, operand2 int) {
	cpu.setAReg(cpu.getCReg())
	cpu.Reg.PC++
}

// LD A, D
func op0x7a(cpu *CPU, operand1, operand2 int) {
	cpu.setAReg(cpu.getDReg())
	cpu.Reg.PC++
}

// LD A, E
func op0x7b(cpu *CPU, operand1, operand2 int) {
	cpu.setAReg(cpu.getEReg())
	cpu.Reg.PC++
}

// LD A, H
func op0x7c(cpu *CPU, operand1, operand2 int) {
	cpu.setAReg(cpu.getHReg())
	cpu.Reg.PC++
}

// LD A, L
func op0x7d(cpu *CPU, operand1, operand2 int) {
	cpu.setAReg(cpu.getLReg())
	cpu.Reg.PC++
}

// LD A, (HL)
func op0x7e(cpu *CPU, operand1, operand2 int) {
	value := cpu.FetchMemory8(cpu.Reg.HL)
	cpu.setAReg(value)
	cpu.Reg.PC++
}

// LD A, A
func op0x7f(cpu *CPU, operand1, operand2 int) {
	cpu.setAReg(cpu.getAReg())
	cpu.Reg.PC++
}

// LD A, (u16)
func op0xfa(cpu *CPU, operand1, operand2 int) {
	addr := cpu.a16FetchJP()
	cpu.setAReg(cpu.FetchMemory8(addr))
	cpu.Reg.PC += 3
	cpu.timer(2)
}

// LD A,(FF00+C)
func op0xf2(cpu *CPU, operand1, operand2 int) {
	addr := 0xff00 + uint16(cpu.getCReg())
	cpu.setAReg(cpu.fetchIO(addr))
	cpu.Reg.PC++ // 誤植(https://www.pastraiser.com/cpu/gameboy/gameboy_opcodes.html)
}

// ------ LD B, *

// LD B,u8
func op0x06(cpu *CPU, operand1, operand2 int) {
	value := cpu.d8Fetch()
	cpu.setBReg(value)
	cpu.Reg.PC += 2
}

// LD B,B
func op0x40(cpu *CPU, operand1, operand2 int) {
	cpu.setBReg(cpu.getBReg())
	cpu.Reg.PC++
}

// LD B,C
func op0x41(cpu *CPU, operand1, operand2 int) {
	cpu.setBReg(cpu.getCReg())
	cpu.Reg.PC++
}

// LD B,D
func op0x42(cpu *CPU, operand1, operand2 int) {
	cpu.setBReg(cpu.getDReg())
	cpu.Reg.PC++
}

// LD B,E
func op0x43(cpu *CPU, operand1, operand2 int) {
	cpu.setBReg(cpu.getEReg())
	cpu.Reg.PC++
}

// LD B,H
func op0x44(cpu *CPU, operand1, operand2 int) {
	cpu.setBReg(cpu.getHReg())
	cpu.Reg.PC++
}

// LD B,L
func op0x45(cpu *CPU, operand1, operand2 int) {
	cpu.setBReg(cpu.getLReg())
	cpu.Reg.PC++
}

// LD B,(HL)
func op0x46(cpu *CPU, operand1, operand2 int) {
	value := cpu.FetchMemory8(cpu.Reg.HL)
	cpu.setBReg(value)
	cpu.Reg.PC++
}

// LD B,A
func op0x47(cpu *CPU, operand1, operand2 int) {
	cpu.setBReg(cpu.getAReg())
	cpu.Reg.PC++
}

// ------ LD C, *

// LD C,u8
func op0x0e(cpu *CPU, operand1, operand2 int) {
	value := cpu.d8Fetch()
	cpu.setCReg(value)
	cpu.Reg.PC += 2
}

// LD C,B
func op0x48(cpu *CPU, operand1, operand2 int) {
	cpu.setCReg(cpu.getBReg())
	cpu.Reg.PC++
}

// LD C,C
func op0x49(cpu *CPU, operand1, operand2 int) {
	cpu.setCReg(cpu.getCReg())
	cpu.Reg.PC++
}

// LD C,D
func op0x4a(cpu *CPU, operand1, operand2 int) {
	cpu.setCReg(cpu.getDReg())
	cpu.Reg.PC++
}

// LD C,E
func op0x4b(cpu *CPU, operand1, operand2 int) {
	cpu.setCReg(cpu.getEReg())
	cpu.Reg.PC++
}

// LD C,H
func op0x4c(cpu *CPU, operand1, operand2 int) {
	cpu.setCReg(cpu.getHReg())
	cpu.Reg.PC++
}

// LD C,L
func op0x4d(cpu *CPU, operand1, operand2 int) {
	cpu.setCReg(cpu.getLReg())
	cpu.Reg.PC++
}

// LD C,(HL)
func op0x4e(cpu *CPU, operand1, operand2 int) {
	value := cpu.FetchMemory8(cpu.Reg.HL)
	cpu.setCReg(value)
	cpu.Reg.PC++
}

// LD C,A
func op0x4f(cpu *CPU, operand1, operand2 int) {
	cpu.setCReg(cpu.getAReg())
	cpu.Reg.PC++
}

// ------ LD D, *

// LD D,u8
func op0x16(cpu *CPU, operand1, operand2 int) {
	value := cpu.d8Fetch()
	cpu.setDReg(value)
	cpu.Reg.PC += 2
}

// LD D,B
func op0x50(cpu *CPU, operand1, operand2 int) {
	cpu.setDReg(cpu.getBReg())
	cpu.Reg.PC++
}

// LD D,C
func op0x51(cpu *CPU, operand1, operand2 int) {
	cpu.setDReg(cpu.getCReg())
	cpu.Reg.PC++
}

// LD D,D
func op0x52(cpu *CPU, operand1, operand2 int) {
	cpu.setDReg(cpu.getDReg())
	cpu.Reg.PC++
}

// LD D,E
func op0x53(cpu *CPU, operand1, operand2 int) {
	cpu.setDReg(cpu.getEReg())
	cpu.Reg.PC++
}

// LD D,H
func op0x54(cpu *CPU, operand1, operand2 int) {
	cpu.setDReg(cpu.getHReg())
	cpu.Reg.PC++
}

// LD D,L
func op0x55(cpu *CPU, operand1, operand2 int) {
	cpu.setDReg(cpu.getLReg())
	cpu.Reg.PC++
}

// LD D,(HL)
func op0x56(cpu *CPU, operand1, operand2 int) {
	value := cpu.FetchMemory8(cpu.Reg.HL)
	cpu.setDReg(value)
	cpu.Reg.PC++
}

// LD D,A
func op0x57(cpu *CPU, operand1, operand2 int) {
	cpu.setDReg(cpu.getAReg())
	cpu.Reg.PC++
}

// ------ LD E, *

// LD E,u8
func op0x1e(cpu *CPU, operand1, operand2 int) {
	value := cpu.d8Fetch()
	cpu.setEReg(value)
	cpu.Reg.PC += 2
}

// LD E,B
func op0x58(cpu *CPU, operand1, operand2 int) {
	cpu.setEReg(cpu.getBReg())
	cpu.Reg.PC++
}

// LD E,C
func op0x59(cpu *CPU, operand1, operand2 int) {
	cpu.setEReg(cpu.getCReg())
	cpu.Reg.PC++
}

// LD E,D
func op0x5a(cpu *CPU, operand1, operand2 int) {
	cpu.setEReg(cpu.getDReg())
	cpu.Reg.PC++
}

// LD E,E
func op0x5b(cpu *CPU, operand1, operand2 int) {
	cpu.setEReg(cpu.getEReg())
	cpu.Reg.PC++
}

// LD E,H
func op0x5c(cpu *CPU, operand1, operand2 int) {
	cpu.setEReg(cpu.getHReg())
	cpu.Reg.PC++
}

// LD E,L
func op0x5d(cpu *CPU, operand1, operand2 int) {
	cpu.setEReg(cpu.getLReg())
	cpu.Reg.PC++
}

// LD E,(HL)
func op0x5e(cpu *CPU, operand1, operand2 int) {
	value := cpu.FetchMemory8(cpu.Reg.HL)
	cpu.setEReg(value)
	cpu.Reg.PC++
}

// LD E,A
func op0x5f(cpu *CPU, operand1, operand2 int) {
	cpu.setEReg(cpu.getAReg())
	cpu.Reg.PC++
}

// ------ LD H, *

// LD H,u8
func op0x26(cpu *CPU, operand1, operand2 int) {
	value := cpu.d8Fetch()
	cpu.setHReg(value)
	cpu.Reg.PC += 2
}

// LD H,B
func op0x60(cpu *CPU, operand1, operand2 int) {
	cpu.setHReg(cpu.getBReg())
	cpu.Reg.PC++
}

// LD H,C
func op0x61(cpu *CPU, operand1, operand2 int) {
	cpu.setHReg(cpu.getCReg())
	cpu.Reg.PC++
}

// LD H,D
func op0x62(cpu *CPU, operand1, operand2 int) {
	cpu.setHReg(cpu.getDReg())
	cpu.Reg.PC++
}

// LD H,E
func op0x63(cpu *CPU, operand1, operand2 int) {
	cpu.setHReg(cpu.getEReg())
	cpu.Reg.PC++
}

// LD H,H
func op0x64(cpu *CPU, operand1, operand2 int) {
	cpu.setHReg(cpu.getHReg())
	cpu.Reg.PC++
}

// LD H,L
func op0x65(cpu *CPU, operand1, operand2 int) {
	cpu.setHReg(cpu.getLReg())
	cpu.Reg.PC++
}

// LD H,(HL)
func op0x66(cpu *CPU, operand1, operand2 int) {
	value := cpu.FetchMemory8(cpu.Reg.HL)
	cpu.setHReg(value)
	cpu.Reg.PC++
}

// LD H,A
func op0x67(cpu *CPU, operand1, operand2 int) {
	cpu.setHReg(cpu.getAReg())
	cpu.Reg.PC++
}

// ------ LD L, *

// LD L,u8
func op0x2e(cpu *CPU, operand1, operand2 int) {
	value := cpu.d8Fetch()
	cpu.setLReg(value)
	cpu.Reg.PC += 2
}

// LD L,B
func op0x68(cpu *CPU, operand1, operand2 int) {
	cpu.setLReg(cpu.getBReg())
	cpu.Reg.PC++
}

// LD L,C
func op0x69(cpu *CPU, operand1, operand2 int) {
	cpu.setLReg(cpu.getCReg())
	cpu.Reg.PC++
}

// LD L,D
func op0x6a(cpu *CPU, operand1, operand2 int) {
	cpu.setLReg(cpu.getDReg())
	cpu.Reg.PC++
}

// LD L,E
func op0x6b(cpu *CPU, operand1, operand2 int) {
	cpu.setLReg(cpu.getEReg())
	cpu.Reg.PC++
}

// LD L,H
func op0x6c(cpu *CPU, operand1, operand2 int) {
	cpu.setLReg(cpu.getHReg())
	cpu.Reg.PC++
}

// LD L,L
func op0x6d(cpu *CPU, operand1, operand2 int) {
	cpu.setLReg(cpu.getLReg())
	cpu.Reg.PC++
}

// LD L,(HL)
func op0x6e(cpu *CPU, operand1, operand2 int) {
	value := cpu.FetchMemory8(cpu.Reg.HL)
	cpu.setLReg(value)
	cpu.Reg.PC++
}

// LD L,A
func op0x6f(cpu *CPU, operand1, operand2 int) {
	cpu.setLReg(cpu.getAReg())
	cpu.Reg.PC++
}

// ------ LD (HL), *

// LD (HL),u8
func op0x36(cpu *CPU, operand1, operand2 int) {
	value := cpu.d8Fetch()
	cpu.timer(1)
	cpu.SetMemory8(cpu.Reg.HL, value)
	cpu.Reg.PC += 2
	cpu.timer(2)
}

// LD (HL),B
func op0x70(cpu *CPU, operand1, operand2 int) {
	cpu.SetMemory8(cpu.Reg.HL, cpu.getBReg())
	cpu.Reg.PC++
}

// LD (HL),C
func op0x71(cpu *CPU, operand1, operand2 int) {
	cpu.SetMemory8(cpu.Reg.HL, cpu.getCReg())
	cpu.Reg.PC++
}

// LD (HL),D
func op0x72(cpu *CPU, operand1, operand2 int) {
	cpu.SetMemory8(cpu.Reg.HL, cpu.getDReg())
	cpu.Reg.PC++
}

// LD (HL),E
func op0x73(cpu *CPU, operand1, operand2 int) {
	cpu.SetMemory8(cpu.Reg.HL, cpu.getEReg())
	cpu.Reg.PC++
}

// LD (HL),H
func op0x74(cpu *CPU, operand1, operand2 int) {
	cpu.SetMemory8(cpu.Reg.HL, cpu.getHReg())
	cpu.Reg.PC++
}

// LD (HL),L
func op0x75(cpu *CPU, operand1, operand2 int) {
	cpu.SetMemory8(cpu.Reg.HL, cpu.getLReg())
	cpu.Reg.PC++
}

// LD (HL),A
func op0x77(cpu *CPU, operand1, operand2 int) {
	cpu.SetMemory8(cpu.Reg.HL, cpu.getAReg())
	cpu.Reg.PC++
}

// ------ その他のLD

// LD (u16),SP
func op0x08(cpu *CPU, operand1, operand2 int) {
	// Store SP into addresses n16 (LSB) and n16 + 1 (MSB).
	addr := cpu.a16Fetch()
	upper := byte(cpu.Reg.SP >> 8)     // MSB
	lower := byte(cpu.Reg.SP & 0x00ff) // LSB
	cpu.SetMemory8(addr, lower)
	cpu.SetMemory8(addr+1, upper)
	cpu.Reg.PC += 3
	cpu.timer(5)
}

// LD (u16),A
func op0xea(cpu *CPU, operand1, operand2 int) {
	addr := cpu.a16FetchJP()
	cpu.SetMemory8(addr, cpu.getAReg())
	cpu.Reg.PC += 3
	cpu.timer(2)
}

// LD BC,u16
func op0x01(cpu *CPU, operand1, operand2 int) {
	value := cpu.d16Fetch()
	cpu.Reg.BC = value
	cpu.Reg.PC += 3
}

// LD DE,u16
func op0x11(cpu *CPU, operand1, operand2 int) {
	value := cpu.d16Fetch()
	cpu.Reg.DE = value
	cpu.Reg.PC += 3
}

// LD HL,u16
func op0x21(cpu *CPU, operand1, operand2 int) {
	cpu.Reg.HL = cpu.d16Fetch()
	cpu.Reg.PC += 3
}

// LD SP,u16
func op0x31(cpu *CPU, operand1, operand2 int) {
	value := cpu.d16Fetch()
	cpu.Reg.SP = value
	cpu.Reg.PC += 3
}

// LD HL,SP+i8
func op0xf8(cpu *CPU, operand1, operand2 int) {
	delta := int8(cpu.FetchMemory8(cpu.Reg.PC + 1))
	value := int32(cpu.Reg.SP) + int32(delta)
	carryBits := uint32(cpu.Reg.SP) ^ uint32(delta) ^ uint32(value)
	cpu.Reg.HL = uint16(value)
	cpu.clearZFlag()
	cpu.flagN(false)
	cpu.flagC8(uint16(carryBits))
	cpu.flagH8(byte(carryBits))
	cpu.Reg.PC += 2
}

// LD SP,HL
func op0xf9(cpu *CPU, operand1, operand2 int) {
	cpu.Reg.SP = cpu.Reg.HL
	cpu.Reg.PC++
}

// LD (FF00+C),A
func op0xe2(cpu *CPU, operand1, operand2 int) {
	addr := 0xff00 + uint16(cpu.getCReg())
	cpu.SetMemory8(addr, cpu.getAReg())
	cpu.Reg.PC++ // 誤植(https://www.pastraiser.com/cpu/gameboy/gameboy_opcodes.html)
}

// LD (BC),A
func op0x02(cpu *CPU, operand1, operand2 int) {
	cpu.SetMemory8(cpu.Reg.BC, cpu.getAReg())
	cpu.Reg.PC++
}

// LD (DE),A
func op0x12(cpu *CPU, operand1, operand2 int) {
	cpu.SetMemory8(cpu.Reg.DE, cpu.getAReg())
	cpu.Reg.PC++
}

// LD (HL+),A
func op0x22(cpu *CPU, operand1, operand2 int) {
	cpu.SetMemory8(cpu.Reg.HL, cpu.getAReg())
	cpu.Reg.HL++
	cpu.Reg.PC++
}

// LD (HL-),A
func op0x32(cpu *CPU, operand1, operand2 int) {
	// (HL)=A, HL=HL-1
	cpu.SetMemory8(cpu.Reg.HL, cpu.getAReg())
	cpu.Reg.HL--
	cpu.Reg.PC++
}

// LD Load
func LD(cpu *CPU, operand1, operand2 int) {
	errMsg := fmt.Sprintf("Error: LD %d %d", operand1, operand2)
	panic(errMsg)
}

// LDH Load High Byte
func LDH(cpu *CPU, operand1, operand2 int) {
	if operand1 == OPERAND_A && operand2 == OPERAND_a8_PAREN {
		// LD A,($FF00+a8)
		addr := 0xff00 + uint16(cpu.FetchMemory8(cpu.Reg.PC+1))
		cpu.timer(1)
		value := cpu.fetchIO(addr)

		cpu.setAReg(value)
		cpu.Reg.PC += 2
		cpu.timer(2)
	} else if operand1 == OPERAND_a8_PAREN && operand2 == OPERAND_A {
		// LD ($FF00+a8),A
		addr := 0xff00 + uint16(cpu.FetchMemory8(cpu.Reg.PC+1))
		cpu.timer(1)
		cpu.setIO(addr, cpu.getAReg())
		cpu.Reg.PC += 2
		cpu.timer(2)
	} else {
		errMsg := fmt.Sprintf("Error: LDH %s %s", operand1, operand2)
		panic(errMsg)
	}
}

// NOP No operation
func (cpu *CPU) NOP(operand1, operand2 int) {
	cpu.Reg.PC++
}

// INC Increment
func (cpu *CPU) INC(operand1, operand2 int) {
	var value, carryBits byte

	switch operand1 {
	case OPERAND_A:
		value = cpu.getAReg() + 1
		carryBits = cpu.getAReg() ^ 1 ^ value
		cpu.setAReg(value)
	case OPERAND_B:
		value = cpu.getBReg() + 1
		carryBits = cpu.getBReg() ^ 1 ^ value
		cpu.setBReg(value)
	case OPERAND_C:
		value = cpu.getCReg() + 1
		carryBits = cpu.getCReg() ^ 1 ^ value
		cpu.setCReg(value)
	case OPERAND_D:
		value = cpu.getDReg() + 1
		carryBits = cpu.getDReg() ^ 1 ^ value
		cpu.setDReg(value)
	case OPERAND_E:
		value = cpu.getEReg() + 1
		carryBits = cpu.getEReg() ^ 1 ^ value
		cpu.setEReg(value)
	case OPERAND_H:
		value = cpu.getHReg() + 1
		carryBits = cpu.getHReg() ^ 1 ^ value
		cpu.setHReg(value)
	case OPERAND_L:
		value = cpu.getLReg() + 1
		carryBits = cpu.getLReg() ^ 1 ^ value
		cpu.setLReg(value)
	case OPERAND_HL_PAREN:
		value = cpu.FetchMemory8(cpu.Reg.HL) + 1
		cpu.timer(1)
		carryBits = cpu.FetchMemory8(cpu.Reg.HL) ^ 1 ^ value
		cpu.SetMemory8(cpu.Reg.HL, value)
		cpu.timer(2)
	case OPERAND_BC:
		cpu.Reg.BC++
	case OPERAND_DE:
		cpu.Reg.DE++
	case OPERAND_HL:
		cpu.Reg.HL++
	case OPERAND_SP:
		cpu.Reg.SP++
	default:
		errMsg := fmt.Sprintf("Error: INC %s %s", operand1, operand2)
		panic(errMsg)
	}

	if operand1 != OPERAND_BC && operand1 != OPERAND_DE && operand1 != OPERAND_HL && operand1 != OPERAND_SP {
		cpu.flagZ(value)
		cpu.flagN(false)
		cpu.flagH8(carryBits)
	}
	cpu.Reg.PC++
}

// DEC Decrement
func (cpu *CPU) DEC(operand1, operand2 int) {
	var value byte
	var carryBits byte

	switch operand1 {
	case OPERAND_A:
		value = cpu.getAReg() - 1
		carryBits = cpu.getAReg() ^ 1 ^ value
		cpu.setAReg(value)
	case OPERAND_B:
		value = cpu.getBReg() - 1
		carryBits = cpu.getBReg() ^ 1 ^ value
		cpu.setBReg(value)
	case OPERAND_C:
		value = cpu.getCReg() - 1
		carryBits = cpu.getCReg() ^ 1 ^ value
		cpu.setCReg(value)
	case OPERAND_D:
		value = cpu.getDReg() - 1
		carryBits = cpu.getDReg() ^ 1 ^ value
		cpu.setDReg(value)
	case OPERAND_E:
		value = cpu.getEReg() - 1
		carryBits = cpu.getEReg() ^ 1 ^ value
		cpu.setEReg(value)
	case OPERAND_H:
		value = cpu.getHReg() - 1
		carryBits = cpu.getHReg() ^ 1 ^ value
		cpu.setHReg(value)
	case OPERAND_L:
		value = cpu.getLReg() - 1
		carryBits = cpu.getLReg() ^ 1 ^ value
		cpu.setLReg(value)
	case OPERAND_HL_PAREN:
		value = cpu.FetchMemory8(cpu.Reg.HL) - 1
		cpu.timer(1)
		carryBits = cpu.FetchMemory8(cpu.Reg.HL) ^ 1 ^ value
		cpu.SetMemory8(cpu.Reg.HL, value)
		cpu.timer(2)
	case OPERAND_BC:
		cpu.Reg.BC--
	case OPERAND_DE:
		cpu.Reg.DE--
	case OPERAND_HL:
		cpu.Reg.HL--
	case OPERAND_SP:
		cpu.Reg.SP--
	default:
		errMsg := fmt.Sprintf("Error: DEC %s %s", operand1, operand2)
		panic(errMsg)
	}

	if operand1 != OPERAND_BC && operand1 != OPERAND_DE && operand1 != OPERAND_HL && operand1 != OPERAND_SP {
		cpu.flagZ(value)
		cpu.flagN(true)
		cpu.flagH8(carryBits)
	}
	cpu.Reg.PC++
}

// --------- JR ---------

// JR i8
func op0x18(cpu *CPU, operand1, operand2 int) {
	delta := int8(cpu.FetchMemory8(cpu.Reg.PC + 1))
	destination := uint16(int32(cpu.Reg.PC+2) + int32(delta)) // この時点でのPCは命令フェッチ後のPCなので+2してあげる必要あり
	cpu.Reg.PC = destination
	cpu.timer(3)
}

// JR NZ,i8
func op0x20(cpu *CPU, operand1, operand2 int) {
	if !cpu.getZFlag() {
		delta := int8(cpu.FetchMemory8(cpu.Reg.PC + 1))
		destination := uint16(int32(cpu.Reg.PC+2) + int32(delta)) // この時点でのPCは命令フェッチ後のPCなので+2してあげる必要あり
		cpu.Reg.PC = destination
		cpu.timer(3)
	} else {
		cpu.Reg.PC += 2
		cpu.timer(2)
	}
}

// JR Z,i8
func op0x28(cpu *CPU, operand1, operand2 int) {
	if cpu.getZFlag() {
		delta := int8(cpu.FetchMemory8(cpu.Reg.PC + 1))
		destination := uint16(int32(cpu.Reg.PC+2) + int32(delta)) // この時点でのPCは命令フェッチ後のPCなので+2してあげる必要あり
		cpu.Reg.PC = destination
		cpu.timer(3)
	} else {
		cpu.Reg.PC += 2
		cpu.timer(2)
	}
}

// JR NC,i8
func op0x30(cpu *CPU, operand1, operand2 int) {
	if !cpu.getCFlag() {
		delta := int8(cpu.FetchMemory8(cpu.Reg.PC + 1))
		destination := uint16(int32(cpu.Reg.PC+2) + int32(delta)) // この時点でのPCは命令フェッチ後のPCなので+2してあげる必要あり
		cpu.Reg.PC = destination
		cpu.timer(3)
	} else {
		cpu.Reg.PC += 2
		cpu.timer(2)
	}
}

// JR C,i8
func op0x38(cpu *CPU, operand1, operand2 int) {
	if cpu.getCFlag() {
		delta := int8(cpu.FetchMemory8(cpu.Reg.PC + 1))
		destination := uint16(int32(cpu.Reg.PC+2) + int32(delta)) // この時点でのPCは命令フェッチ後のPCなので+2してあげる必要あり
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
	case OPERAND_r8:
		delta := int8(cpu.FetchMemory8(cpu.Reg.PC + 1))
		destination := uint16(int32(cpu.Reg.PC+2) + int32(delta)) // この時点でのPCは命令フェッチ後のPCなので+2してあげる必要あり
		cpu.Reg.PC = destination
	case OPERAND_Z:
		if cpu.getZFlag() {
			delta := int8(cpu.FetchMemory8(cpu.Reg.PC + 1))
			destination := uint16(int32(cpu.Reg.PC+2) + int32(delta)) // この時点でのPCは命令フェッチ後のPCなので+2してあげる必要あり
			cpu.Reg.PC = destination
		} else {
			cpu.Reg.PC += 2
			result = false
		}
	case OPERAND_C:
		if cpu.getCFlag() {
			delta := int8(cpu.FetchMemory8(cpu.Reg.PC + 1))
			destination := uint16(int32(cpu.Reg.PC+2) + int32(delta)) // この時点でのPCは命令フェッチ後のPCなので+2してあげる必要あり
			cpu.Reg.PC = destination
		} else {
			cpu.Reg.PC += 2
			result = false
		}
	case OPERAND_NZ:
		if !cpu.getZFlag() {
			delta := int8(cpu.FetchMemory8(cpu.Reg.PC + 1))
			destination := uint16(int32(cpu.Reg.PC+2) + int32(delta)) // この時点でのPCは命令フェッチ後のPCなので+2してあげる必要あり
			cpu.Reg.PC = destination
		} else {
			cpu.Reg.PC += 2
			result = false
		}
	case OPERAND_NC:
		if !cpu.getCFlag() {
			delta := int8(cpu.FetchMemory8(cpu.Reg.PC + 1))
			destination := uint16(int32(cpu.Reg.PC+2) + int32(delta)) // この時点でのPCは命令フェッチ後のPCなので+2してあげる必要あり
			cpu.Reg.PC = destination
		} else {
			cpu.Reg.PC += 2
			result = false
		}
	default:
		errMsg := fmt.Sprintf("Error: JR %s %s", operand1, operand2)
		panic(errMsg)
	}

	if result {
		cpu.timer(3)
	} else {
		cpu.timer(2)
	}
}

// HALT Halt
func (cpu *CPU) HALT(operand1, operand2 int) {
	cpu.Reg.PC++
	cpu.halt = true

	// ref: https://rednex.github.io/rgbds/gbz80.7.html#HALT
	if !cpu.Reg.IME {
		IE, IF := cpu.RAM[IEIO], cpu.RAM[IFIO]
		pending := IE&IF != 0
		if pending {
			// Some pending
			cpu.halt = false
			PC := cpu.Reg.PC
			cpu.exec()
			cpu.Reg.PC = PC

			// IME turns on due to EI delay.
			cpu.halt = cpu.Reg.IME
		}
	}
}

// STOP stop CPU
func (cpu *CPU) STOP(operand1, operand2 int) {
	if operand1 == OPERAND_0 && operand2 == OPERAND_NONE {
		cpu.Reg.PC += 2
		// 速度切り替え
		KEY1 := cpu.FetchMemory8(KEY1IO)
		if KEY1&0x01 == 1 {
			if KEY1>>7 == 1 {
				KEY1 = 0x00
				cpu.boost = 1
			} else {
				KEY1 = 0x80
				cpu.boost = 2
			}
			cpu.SetMemory8(KEY1IO, KEY1)
		}
	} else {
		errMsg := fmt.Sprintf("Error: STOP %s %s", operand1, operand2)
		panic(errMsg)
	}
}

// XOR xor
func (cpu *CPU) XOR(operand1, operand2 int) {
	var value byte
	switch operand1 {
	case OPERAND_B:
		value = cpu.getAReg() ^ cpu.getBReg()
	case OPERAND_C:
		value = cpu.getAReg() ^ cpu.getCReg()
	case OPERAND_D:
		value = cpu.getAReg() ^ cpu.getDReg()
	case OPERAND_E:
		value = cpu.getAReg() ^ cpu.getEReg()
	case OPERAND_H:
		value = cpu.getAReg() ^ cpu.getHReg()
	case OPERAND_L:
		value = cpu.getAReg() ^ cpu.getLReg()
	case OPERAND_HL_PAREN:
		value = cpu.getAReg() ^ cpu.FetchMemory8(cpu.Reg.HL)
	case OPERAND_A:
		value = cpu.getAReg() ^ cpu.getAReg()
	case OPERAND_d8:
		value = cpu.getAReg() ^ cpu.FetchMemory8(cpu.Reg.PC+1)
		cpu.Reg.PC++
	default:
		errMsg := fmt.Sprintf("Error: XOR %s %s", operand1, operand2)
		panic(errMsg)
	}

	cpu.setAReg(value)
	cpu.flagZ(value)
	cpu.flagN(false)
	cpu.clearHFlag()
	cpu.clearCFlag()
	cpu.Reg.PC++
}

// JP Jump
func JP(cpu *CPU, operand1, operand2 int) {
	cycle := 1

	switch operand1 {
	case OPERAND_a16:
		destination := cpu.a16FetchJP()
		cycle++
		cpu.Reg.PC = destination
	case OPERAND_HL_PAREN:
		destination := cpu.Reg.HL
		cpu.Reg.PC = destination
	case OPERAND_Z:
		destination := cpu.a16FetchJP()
		if cpu.getZFlag() {
			cycle++
			cpu.Reg.PC = destination
		} else {
			cpu.Reg.PC += 3
		}
	case OPERAND_C:
		destination := cpu.a16FetchJP()
		if cpu.getCFlag() {
			cycle++
			cpu.Reg.PC = destination
		} else {
			cpu.Reg.PC += 3
		}
	case OPERAND_NZ:
		destination := cpu.a16FetchJP()
		if !cpu.getZFlag() {
			cycle++
			cpu.Reg.PC = destination
		} else {
			cpu.Reg.PC += 3
		}
	case OPERAND_NC:
		destination := cpu.a16FetchJP()
		if !cpu.getCFlag() {
			cycle++
			cpu.Reg.PC = destination
		} else {
			cpu.Reg.PC += 3
		}
	default:
		errMsg := fmt.Sprintf("Error: JP %s %s", operand1, operand2)
		panic(errMsg)
	}

	cpu.timer(cycle)
}

// RET Return
func (cpu *CPU) RET(operand1, operand2 int) (result bool) {
	result = true

	switch operand1 {
	case OPERAND_NONE:
		// PC=(SP), SP=SP+2
		cpu.popPC()
	case OPERAND_Z:
		if cpu.getZFlag() {
			cpu.popPC()
		} else {
			cpu.Reg.PC++
			result = false
		}
	case OPERAND_C:
		if cpu.getCFlag() {
			cpu.popPC()
		} else {
			cpu.Reg.PC++
			result = false
		}
	case OPERAND_NZ:
		if !cpu.getZFlag() {
			cpu.popPC()
		} else {
			cpu.Reg.PC++
			result = false
		}
	case OPERAND_NC:
		if !cpu.getCFlag() {
			cpu.popPC()
		} else {
			cpu.Reg.PC++
			result = false
		}
	default:
		errMsg := fmt.Sprintf("Error: RET %s %s", operand1, operand2)
		panic(errMsg)
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
	case OPERAND_a16:
		destination := cpu.a16FetchJP()
		cpu.Reg.PC += 3
		cpu.timer(1)
		cpu.pushPCCALL()
		cpu.timer(1)
		cpu.Reg.PC = destination
	case OPERAND_Z:
		if cpu.getZFlag() {
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
	case OPERAND_C:
		if cpu.getCFlag() {
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
	case OPERAND_NZ:
		if !cpu.getZFlag() {
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
	case OPERAND_NC:
		if !cpu.getCFlag() {
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
		errMsg := fmt.Sprintf("Error: CALL %s %s", operand1, operand2)
		panic(errMsg)
	}
}

// DI Disable Interrupt
func (cpu *CPU) DI(operand1, operand2 int) {
	cpu.Reg.IME = false
	cpu.Reg.PC++
	if cpu.IMESwitch.Working && cpu.IMESwitch.Value {
		cpu.IMESwitch.Working = false // https://gbdev.gg8.se/wiki/articles/Interrupts 『The effect of EI is delayed by one instruction. This means that EI followed immediately by DI does not allow interrupts between the EI and the DI.』
	}
}

// EI Enable Interrupt
func (cpu *CPU) EI(operand1, operand2 int) {
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
func (cpu *CPU) CP(operand1, operand2 int) {
	var value, carryBits byte

	switch operand1 {
	case OPERAND_A:
		value = cpu.getAReg() - cpu.getAReg()
		carryBits = cpu.getAReg() ^ cpu.getAReg() ^ value
		cpu.flagC8Sub(cpu.getAReg(), cpu.getAReg())
	case OPERAND_B:
		value = cpu.getAReg() - cpu.getBReg()
		carryBits = cpu.getAReg() ^ cpu.getBReg() ^ value
		cpu.flagC8Sub(cpu.getAReg(), cpu.getBReg())
	case OPERAND_C:
		value = cpu.getAReg() - cpu.getCReg()
		carryBits = cpu.getAReg() ^ cpu.getCReg() ^ value
		cpu.flagC8Sub(cpu.getAReg(), cpu.getCReg())
	case OPERAND_D:
		value = cpu.getAReg() - cpu.getDReg()
		carryBits = cpu.getAReg() ^ cpu.getDReg() ^ value
		cpu.flagC8Sub(cpu.getAReg(), cpu.getDReg())
	case OPERAND_E:
		value = cpu.getAReg() - cpu.getEReg()
		carryBits = cpu.getAReg() ^ cpu.getEReg() ^ value
		cpu.flagC8Sub(cpu.getAReg(), cpu.getEReg())
	case OPERAND_H:
		value = cpu.getAReg() - cpu.getHReg()
		carryBits = cpu.getAReg() ^ cpu.getHReg() ^ value
		cpu.flagC8Sub(cpu.getAReg(), cpu.getHReg())
	case OPERAND_L:
		value = cpu.getAReg() - cpu.getLReg()
		carryBits = cpu.getAReg() ^ cpu.getLReg() ^ value
		cpu.flagC8Sub(cpu.getAReg(), cpu.getLReg())
	case OPERAND_d8:
		value = cpu.getAReg() - cpu.d8Fetch()
		carryBits = cpu.getAReg() ^ cpu.d8Fetch() ^ value
		cpu.flagC8Sub(cpu.getAReg(), cpu.d8Fetch())
		cpu.Reg.PC++
	case OPERAND_HL_PAREN:
		value = cpu.getAReg() - cpu.FetchMemory8(cpu.Reg.HL)
		carryBits = cpu.getAReg() ^ cpu.FetchMemory8(cpu.Reg.HL) ^ value
		cpu.flagC8Sub(cpu.getAReg(), cpu.FetchMemory8(cpu.Reg.HL))
	default:
		errMsg := fmt.Sprintf("Error: CP %s %s", operand1, operand2)
		panic(errMsg)
	}
	cpu.flagZ(value)
	cpu.flagN(true)
	cpu.flagH8(carryBits)
	cpu.Reg.PC++
}

// AND And instruction
func (cpu *CPU) AND(operand1, operand2 int) {
	var value byte
	switch operand1 {
	case OPERAND_A:
		value = cpu.getAReg() & cpu.getAReg()
	case OPERAND_B:
		value = cpu.getAReg() & cpu.getBReg()
	case OPERAND_C:
		value = cpu.getAReg() & cpu.getCReg()
	case OPERAND_D:
		value = cpu.getAReg() & cpu.getDReg()
	case OPERAND_E:
		value = cpu.getAReg() & cpu.getEReg()
	case OPERAND_H:
		value = cpu.getAReg() & cpu.getHReg()
	case OPERAND_L:
		value = cpu.getAReg() & cpu.getLReg()
	case OPERAND_HL_PAREN:
		value = cpu.getAReg() & cpu.FetchMemory8(cpu.Reg.HL)
	case OPERAND_d8:
		value = cpu.getAReg() & cpu.d8Fetch()
		cpu.Reg.PC++
	default:
		errMsg := fmt.Sprintf("Error: AND %s %s", operand1, operand2)
		panic(errMsg)
	}

	cpu.setAReg(value)
	cpu.flagZ(value)
	cpu.flagN(false)
	cpu.setHFlag()
	cpu.clearCFlag()
	cpu.Reg.PC++
}

// OR or
func (cpu *CPU) OR(operand1, operand2 int) {
	switch operand1 {
	case OPERAND_A:
		value := cpu.getAReg() | cpu.getAReg()
		cpu.setAReg(value)
		cpu.flagZ(value)
	case OPERAND_B:
		value := cpu.getAReg() | cpu.getBReg()
		cpu.setAReg(value)
		cpu.flagZ(value)
	case OPERAND_C:
		value := cpu.getAReg() | cpu.getCReg()
		cpu.setAReg(value)
		cpu.flagZ(value)
	case OPERAND_D:
		value := cpu.getAReg() | cpu.getDReg()
		cpu.setAReg(value)
		cpu.flagZ(value)
	case OPERAND_E:
		value := cpu.getAReg() | cpu.getEReg()
		cpu.setAReg(value)
		cpu.flagZ(value)
	case OPERAND_H:
		value := cpu.getAReg() | cpu.getHReg()
		cpu.setAReg(value)
		cpu.flagZ(value)
	case OPERAND_L:
		value := cpu.getAReg() | cpu.getLReg()
		cpu.setAReg(value)
		cpu.flagZ(value)
	case OPERAND_d8:
		value := cpu.getAReg() | cpu.FetchMemory8(cpu.Reg.PC+1)
		cpu.setAReg(value)
		cpu.flagZ(value)
		cpu.Reg.PC++
	case OPERAND_HL_PAREN:
		value := cpu.getAReg() | cpu.FetchMemory8(cpu.Reg.HL)
		cpu.setAReg(value)
		cpu.flagZ(value)
	default:
		errMsg := fmt.Sprintf("Error: OR %s %s", operand1, operand2)
		panic(errMsg)
	}

	cpu.flagN(false)
	cpu.clearHFlag()
	cpu.clearCFlag()
	cpu.Reg.PC++
}

// ADD Addition
func (cpu *CPU) ADD(operand1, operand2 int) {
	switch operand1 {
	case OPERAND_A:
		switch operand2 {
		case OPERAND_A:
			value := uint16(cpu.getAReg()) + uint16(cpu.getAReg())
			carryBits := uint16(cpu.getAReg()) ^ uint16(cpu.getAReg()) ^ value
			cpu.setAReg(byte(value))
			cpu.flagZ(byte(value))
			cpu.flagN(false)
			cpu.flagH8(byte(carryBits))
			cpu.flagC8(carryBits)
			cpu.Reg.PC++
		case OPERAND_B:
			value := uint16(cpu.getAReg()) + uint16(cpu.getBReg())
			carryBits := uint16(cpu.getAReg()) ^ uint16(cpu.getBReg()) ^ value
			cpu.setAReg(byte(value))
			cpu.flagZ(byte(value))
			cpu.flagN(false)
			cpu.flagH8(byte(carryBits))
			cpu.flagC8(carryBits)
			cpu.Reg.PC++
		case OPERAND_C:
			value := uint16(cpu.getAReg()) + uint16(cpu.getCReg())
			carryBits := uint16(cpu.getAReg()) ^ uint16(cpu.getCReg()) ^ value
			cpu.setAReg(byte(value))
			cpu.flagZ(byte(value))
			cpu.flagN(false)
			cpu.flagH8(byte(carryBits))
			cpu.flagC8(carryBits)
			cpu.Reg.PC++
		case OPERAND_D:
			value := uint16(cpu.getAReg()) + uint16(cpu.getDReg())
			carryBits := uint16(cpu.getAReg()) ^ uint16(cpu.getDReg()) ^ value
			cpu.setAReg(byte(value))
			cpu.flagZ(byte(value))
			cpu.flagN(false)
			cpu.flagH8(byte(carryBits))
			cpu.flagC8(carryBits)
			cpu.Reg.PC++
		case OPERAND_E:
			value := uint16(cpu.getAReg()) + uint16(cpu.getEReg())
			carryBits := uint16(cpu.getAReg()) ^ uint16(cpu.getEReg()) ^ value
			cpu.setAReg(byte(value))
			cpu.flagZ(byte(value))
			cpu.flagN(false)
			cpu.flagH8(byte(carryBits))
			cpu.flagC8(carryBits)
			cpu.Reg.PC++
		case OPERAND_H:
			value := uint16(cpu.getAReg()) + uint16(cpu.getHReg())
			carryBits := uint16(cpu.getAReg()) ^ uint16(cpu.getHReg()) ^ value
			cpu.setAReg(byte(value))
			cpu.flagZ(byte(value))
			cpu.flagN(false)
			cpu.flagH8(byte(carryBits))
			cpu.flagC8(carryBits)
			cpu.Reg.PC++
		case OPERAND_L:
			value := uint16(cpu.getAReg()) + uint16(cpu.getLReg())
			carryBits := uint16(cpu.getAReg()) ^ uint16(cpu.getLReg()) ^ value
			cpu.setAReg(byte(value))
			cpu.flagZ(byte(value))
			cpu.flagN(false)
			cpu.flagH8(byte(carryBits))
			cpu.flagC8(carryBits)
			cpu.Reg.PC++
		case OPERAND_d8:
			value := uint16(cpu.getAReg()) + uint16(cpu.d8Fetch())
			carryBits := uint16(cpu.getAReg()) ^ uint16(cpu.d8Fetch()) ^ value
			cpu.setAReg(byte(value))
			cpu.flagZ(byte(value))
			cpu.flagN(false)
			cpu.flagH8(byte(carryBits))
			cpu.flagC8(carryBits)
			cpu.Reg.PC += 2
		case OPERAND_HL_PAREN:
			value := uint16(cpu.getAReg()) + uint16(cpu.FetchMemory8(cpu.Reg.HL))
			carryBits := uint16(cpu.getAReg()) ^ uint16(cpu.FetchMemory8(cpu.Reg.HL)) ^ value
			cpu.setAReg(byte(value))
			cpu.flagZ(byte(value))
			cpu.flagN(false)
			cpu.flagH8(byte(carryBits))
			cpu.flagC8(carryBits)
			cpu.Reg.PC++
		}
	case OPERAND_HL:
		switch operand2 {
		case OPERAND_BC:
			value := uint32(cpu.Reg.HL) + uint32(cpu.Reg.BC)
			carryBits := uint32(cpu.Reg.HL) ^ uint32(cpu.Reg.BC) ^ value
			cpu.Reg.HL = uint16(value)
			cpu.flagN(false)
			cpu.flagH16(uint16(carryBits))
			cpu.flagC16(carryBits)
			cpu.Reg.PC++
		case OPERAND_DE:
			value := uint32(cpu.Reg.HL) + uint32(cpu.Reg.DE)
			carryBits := uint32(cpu.Reg.HL) ^ uint32(cpu.Reg.DE) ^ value
			cpu.Reg.HL = uint16(value)
			cpu.flagN(false)
			cpu.flagH16(uint16(carryBits))
			cpu.flagC16(carryBits)
			cpu.Reg.PC++
		case OPERAND_HL:
			value := uint32(cpu.Reg.HL) + uint32(cpu.Reg.HL)
			carryBits := uint32(cpu.Reg.HL) ^ uint32(cpu.Reg.HL) ^ value
			cpu.Reg.HL = uint16(value)
			cpu.flagN(false)
			cpu.flagH16(uint16(carryBits))
			cpu.flagC16(carryBits)
			cpu.Reg.PC++
		case OPERAND_SP:
			value := uint32(cpu.Reg.HL) + uint32(cpu.Reg.SP)
			carryBits := uint32(cpu.Reg.HL) ^ uint32(cpu.Reg.SP) ^ value
			cpu.Reg.HL = uint16(value)
			cpu.flagN(false)
			cpu.flagH16(uint16(carryBits))
			cpu.flagC16(carryBits)
			cpu.Reg.PC++
		}
	case OPERAND_SP:
		switch operand2 {
		case OPERAND_r8:
			delta := int8(cpu.FetchMemory8(cpu.Reg.PC + 1))
			value := int32(cpu.Reg.SP) + int32(delta)
			carryBits := uint32(cpu.Reg.SP) ^ uint32(delta) ^ uint32(value)
			cpu.Reg.SP = uint16(value)
			cpu.clearZFlag()
			cpu.flagN(false)
			cpu.flagH8(byte(carryBits))
			cpu.flagC8(uint16(carryBits))
			cpu.Reg.PC += 2
		}
	default:
		errMsg := fmt.Sprintf("Error: ADD %s %s", operand1, operand2)
		panic(errMsg)
	}
}

// CPL Complement A Register(Aレジスタのbitをすべて反転)
func (cpu *CPU) CPL(operand1, operand2 int) {
	A := ^cpu.getAReg()
	cpu.setAReg(A)
	cpu.flagN(true)
	cpu.setHFlag()
	cpu.Reg.PC++
}

// PREFIXCB 拡張命令
func (cpu *CPU) PREFIXCB(operand1, operand2 int) {
	if operand1 == OPERAND_NONE && operand2 == OPERAND_NONE {
		cpu.Reg.PC++
		cpu.timer(1)
		opcode := prefixCBs[cpu.FetchMemory8(cpu.Reg.PC)]
		instruction, operand1, operand2, cycle := opcode.Ins, opcode.Operand1, opcode.Operand2, opcode.Cycle1

		switch instruction {
		case INS_RLC:
			cpu.RLC(operand1, operand2)
		case INS_RRC:
			cpu.RRC(operand1, operand2)
		case INS_RL:
			cpu.RL(operand1, operand2)
		case INS_RR:
			cpu.RR(operand1, operand2)
		case INS_SLA:
			cpu.SLA(operand1, operand2)
		case INS_SRA:
			cpu.SRA(operand1, operand2)
		case INS_SWAP:
			cpu.SWAP(operand1, operand2)
		case INS_SRL:
			cpu.SRL(operand1, operand2)
		case INS_BIT:
			cpu.BIT(operand1, operand2)
		case INS_RES:
			cpu.RES(operand1, operand2)
		case INS_SET:
			cpu.SET(operand1, operand2)
		default:
			errMsg := fmt.Sprintf("eip: 0x%04x opcode: 0x%02x", cpu.Reg.PC, opcode)
			panic(errMsg)
		}

		if cycle > 1 {
			cpu.timer(cycle - 1)
		}
	} else {
		errMsg := fmt.Sprintf("Error: PREFIXCB %s %s", operand1, operand2)
		panic(errMsg)
	}
}

// RLC Rotate n left carry => bit0
func (cpu *CPU) RLC(operand1, operand2 int) {
	var value byte
	var bit7 byte
	if operand1 == OPERAND_B && operand2 == OPERAND_NONE {
		value = cpu.getBReg()
		bit7 = value >> 7
		value = (value << 1)
		if bit7 != 0 {
			value |= 1
		} else {
			value &= 0xfe
		}
		cpu.setBReg(value)
	} else if operand1 == OPERAND_C && operand2 == OPERAND_NONE {
		value = cpu.getCReg()
		bit7 = value >> 7
		value = (value << 1)
		if bit7 != 0 {
			value |= 1
		} else {
			value &= 0xfe
		}
		cpu.setCReg(value)
	} else if operand1 == OPERAND_D && operand2 == OPERAND_NONE {
		value = cpu.getDReg()
		bit7 = value >> 7
		value = (value << 1)
		if bit7 != 0 {
			value |= 1
		} else {
			value &= 0xfe
		}
		cpu.setDReg(value)
	} else if operand1 == OPERAND_E && operand2 == OPERAND_NONE {
		value = cpu.getEReg()
		bit7 = value >> 7
		value = (value << 1)
		if bit7 != 0 {
			value |= 1
		} else {
			value &= 0xfe
		}
		cpu.setEReg(value)
	} else if operand1 == OPERAND_H && operand2 == OPERAND_NONE {
		value = cpu.getHReg()
		bit7 = value >> 7
		value = (value << 1)
		if bit7 != 0 {
			value |= 1
		} else {
			value &= 0xfe
		}
		cpu.setHReg(value)
	} else if operand1 == OPERAND_L && operand2 == OPERAND_NONE {
		value = cpu.getLReg()
		bit7 = value >> 7
		value = (value << 1)
		if bit7 != 0 {
			value |= 1
		} else {
			value &= 0xfe
		}
		cpu.setLReg(value)
	} else if operand1 == OPERAND_HL_PAREN && operand2 == OPERAND_NONE {
		value = cpu.FetchMemory8(cpu.Reg.HL)
		cpu.timer(1)
		bit7 = value >> 7
		value = (value << 1)
		if bit7 != 0 {
			value |= 1
		} else {
			value &= 0xfe
		}
		cpu.SetMemory8(cpu.Reg.HL, value)
		cpu.timer(2)
	} else if operand1 == OPERAND_A && operand2 == OPERAND_NONE {
		value = cpu.getAReg()
		bit7 = value >> 7
		value = (value << 1)
		if bit7 != 0 {
			value |= 1
		} else {
			value &= 0xfe
		}
		cpu.setAReg(value)
	} else {
		errMsg := fmt.Sprintf("Error: RLC %s %s", operand1, operand2)
		panic(errMsg)
	}

	cpu.flagZ(value)
	cpu.flagN(false)
	cpu.clearHFlag()
	if bit7 != 0 {
		cpu.setCFlag()
	} else {
		cpu.clearCFlag()
	}
	cpu.Reg.PC++
}

// RLCA Rotate register A left.
func (cpu *CPU) RLCA(operand1, operand2 int) {
	var value byte
	var bit7 byte
	value = cpu.getAReg()
	bit7 = value >> 7
	value = (value << 1)
	if bit7 != 0 {
		value |= 1
	} else {
		value &= 0xfe
	}
	cpu.setAReg(value)

	cpu.clearZFlag()
	cpu.flagN(false)
	cpu.clearHFlag()
	if bit7 != 0 {
		cpu.setCFlag()
	} else {
		cpu.clearCFlag()
	}
	cpu.Reg.PC++
}

// RRC Rotate n right carry => bit7
func (cpu *CPU) RRC(operand1, operand2 int) {
	var value byte
	var bit0 byte
	if operand1 == OPERAND_B && operand2 == OPERAND_NONE {
		value = cpu.getBReg()
		bit0 = value % 2
		value = (value >> 1)
		if bit0 != 0 {
			value |= 0x80
		} else {
			value &= 0x7f
		}
		cpu.setBReg(value)
	} else if operand1 == OPERAND_C && operand2 == OPERAND_NONE {
		value = cpu.getCReg()
		bit0 = value % 2
		value = (value >> 1)
		if bit0 != 0 {
			value |= 0x80
		} else {
			value &= 0x7f
		}
		cpu.setCReg(value)
	} else if operand1 == OPERAND_D && operand2 == OPERAND_NONE {
		value = cpu.getDReg()
		bit0 = value % 2
		value = (value >> 1)
		if bit0 != 0 {
			value |= 0x80
		} else {
			value &= 0x7f
		}
		cpu.setDReg(value)
	} else if operand1 == OPERAND_E && operand2 == OPERAND_NONE {
		value = cpu.getEReg()
		bit0 = value % 2
		value = (value >> 1)
		if bit0 != 0 {
			value |= 0x80
		} else {
			value &= 0x7f
		}
		cpu.setEReg(value)
	} else if operand1 == OPERAND_H && operand2 == OPERAND_NONE {
		value = cpu.getHReg()
		bit0 = value % 2
		value = (value >> 1)
		if bit0 != 0 {
			value |= 0x80
		} else {
			value &= 0x7f
		}
		cpu.setHReg(value)
	} else if operand1 == OPERAND_L && operand2 == OPERAND_NONE {
		value = cpu.getLReg()
		bit0 = value % 2
		value = (value >> 1)
		if bit0 != 0 {
			value |= 0x80
		} else {
			value &= 0x7f
		}
		cpu.setLReg(value)
	} else if operand1 == OPERAND_HL_PAREN && operand2 == OPERAND_NONE {
		value = cpu.FetchMemory8(cpu.Reg.HL)
		cpu.timer(1)
		bit0 = value % 2
		value = (value >> 1)
		if bit0 != 0 {
			value |= 0x80
		} else {
			value &= 0x7f
		}
		cpu.SetMemory8(cpu.Reg.HL, value)
		cpu.timer(2)
	} else if operand1 == OPERAND_A && operand2 == OPERAND_NONE {
		value = cpu.getAReg()
		bit0 = value % 2
		value = (value >> 1)
		if bit0 != 0 {
			value |= 0x80
		} else {
			value &= 0x7f
		}
		cpu.setAReg(value)
	} else {
		errMsg := fmt.Sprintf("Error: RRC %s %s", operand1, operand2)
		panic(errMsg)
	}
	cpu.flagZ(value)
	cpu.flagN(false)
	cpu.clearHFlag()
	if bit0 != 0 {
		cpu.setCFlag()
	} else {
		cpu.clearCFlag()
	}
	cpu.Reg.PC++
}

// RRCA Rotate register A right.
func (cpu *CPU) RRCA(operand1, operand2 int) {
	var value byte
	var bit0 byte

	value = cpu.getAReg()
	bit0 = value % 2
	value = (value >> 1)
	if bit0 != 0 {
		value |= 0x80
	} else {
		value &= 0x7f
	}
	cpu.setAReg(value)

	cpu.clearZFlag()
	cpu.flagN(false)
	cpu.clearHFlag()
	if bit0 != 0 {
		cpu.setCFlag()
	} else {
		cpu.clearCFlag()
	}
	cpu.Reg.PC++
}

// RL Rotate n rigth through carry bit7 => bit0
func (cpu *CPU) RL(operand1, operand2 int) {
	var value byte
	var bit7 byte
	carry := cpu.getCFlag()

	switch operand1 {
	case OPERAND_A:
		value = cpu.getAReg()
		bit7 = value >> 7
		value = (value << 1)
		if carry {
			value |= 1
		} else {
			value &= 0xfe
		}
		cpu.setAReg(value)
	case OPERAND_B:
		value = cpu.getBReg()
		bit7 = value >> 7
		value = (value << 1)
		if carry {
			value |= 1
		} else {
			value &= 0xfe
		}
		cpu.setBReg(value)
	case OPERAND_C:
		value = cpu.getCReg()
		bit7 = value >> 7
		value = (value << 1)
		if carry {
			value |= 1
		} else {
			value &= 0xfe
		}
		cpu.setCReg(value)
	case OPERAND_D:
		value = cpu.getDReg()
		bit7 = value >> 7
		value = (value << 1)
		if carry {
			value |= 1
		} else {
			value &= 0xfe
		}
		cpu.setDReg(value)
	case OPERAND_E:
		value = cpu.getEReg()
		bit7 = value >> 7
		value = (value << 1)
		if carry {
			value |= 1
		} else {
			value &= 0xfe
		}
		cpu.setEReg(value)
	case OPERAND_H:
		value = cpu.getHReg()
		bit7 = value >> 7
		value = (value << 1)
		if carry {
			value |= 1
		} else {
			value &= 0xfe
		}
		cpu.setHReg(value)
	case OPERAND_L:
		value = cpu.getLReg()
		bit7 = value >> 7
		value = (value << 1)
		if carry {
			value |= 1
		} else {
			value &= 0xfe
		}
		cpu.setLReg(value)
	case OPERAND_HL_PAREN:
		value = cpu.FetchMemory8(cpu.Reg.HL)
		cpu.timer(1)
		bit7 = value >> 7
		value = (value << 1)
		if carry {
			value |= 1
		} else {
			value &= 0xfe
		}
		cpu.SetMemory8(cpu.Reg.HL, value)
		cpu.timer(2)
	default:
		errMsg := fmt.Sprintf("Error: RL %s %s", operand1, operand2)
		panic(errMsg)
	}

	cpu.flagZ(value)
	cpu.flagN(false)
	cpu.clearHFlag()
	if bit7 != 0 {
		cpu.setCFlag()
	} else {
		cpu.clearCFlag()
	}
	cpu.Reg.PC++
}

// RLA Rotate register A left through carry.
func (cpu *CPU) RLA(operand1, operand2 int) {
	var value byte
	var bit7 byte
	carry := cpu.getCFlag()

	value = cpu.getAReg()
	bit7 = value >> 7
	value = (value << 1)
	if carry {
		value |= 1
	} else {
		value &= 0xfe
	}
	cpu.setAReg(value)

	cpu.clearZFlag()
	cpu.flagN(false)
	cpu.clearHFlag()
	if bit7 != 0 {
		cpu.setCFlag()
	} else {
		cpu.clearCFlag()
	}
	cpu.Reg.PC++
}

// RR Rotate n right through carry bit0 => bit7
func (cpu *CPU) RR(operand1, operand2 int) {
	var value byte
	var bit0 byte
	carry := cpu.getCFlag()

	switch operand1 {
	case OPERAND_A:
		value = cpu.getAReg()
		bit0 = value % 2
		value = (value >> 1)
		if carry {
			value |= 0x80
		} else {
			value &= 0x7f
		}
		cpu.setAReg(value)
	case OPERAND_B:
		value = cpu.getBReg()
		bit0 = value % 2
		value = (value >> 1)
		if carry {
			value |= 0x80
		} else {
			value &= 0x7f
		}
		cpu.setBReg(value)
	case OPERAND_C:
		value = cpu.getCReg()
		bit0 = value % 2
		value = (value >> 1)
		if carry {
			value |= 0x80
		} else {
			value &= 0x7f
		}
		cpu.setCReg(value)
	case OPERAND_D:
		value = cpu.getDReg()
		bit0 = value % 2
		value = (value >> 1)
		if carry {
			value |= 0x80
		} else {
			value &= 0x7f
		}
		cpu.setDReg(value)
	case OPERAND_E:
		value = cpu.getEReg()
		bit0 = value % 2
		value = (value >> 1)
		if carry {
			value |= 0x80
		} else {
			value &= 0x7f
		}
		cpu.setEReg(value)
	case OPERAND_H:
		value = cpu.getHReg()
		bit0 = value % 2
		value = (value >> 1)
		if carry {
			value |= 0x80
		} else {
			value &= 0x7f
		}
		cpu.setHReg(value)
	case OPERAND_L:
		value = cpu.getLReg()
		bit0 = value % 2
		value = (value >> 1)
		if carry {
			value |= 0x80
		} else {
			value &= 0x7f
		}
		cpu.setLReg(value)
	case OPERAND_HL_PAREN:
		value = cpu.FetchMemory8(cpu.Reg.HL)
		cpu.timer(1)
		bit0 = value % 2
		value = (value >> 1)
		if carry {
			value |= 0x80
		} else {
			value &= 0x7f
		}
		cpu.SetMemory8(cpu.Reg.HL, value)
		cpu.timer(2)
	default:
		errMsg := fmt.Sprintf("Error: RR %s %s", operand1, operand2)
		panic(errMsg)
	}

	cpu.flagZ(value)
	cpu.flagN(false)
	cpu.clearHFlag()
	if bit0 != 0 {
		cpu.setCFlag()
	} else {
		cpu.clearCFlag()
	}
	cpu.Reg.PC++
}

// SLA Shift Left
func (cpu *CPU) SLA(operand1, operand2 int) {
	var value byte
	var bit7 byte
	if operand1 == OPERAND_B && operand2 == OPERAND_NONE {
		value = cpu.getBReg()
		bit7 = value >> 7
		value = (value << 1)
		cpu.setBReg(value)
	} else if operand1 == OPERAND_C && operand2 == OPERAND_NONE {
		value = cpu.getCReg()
		bit7 = value >> 7
		value = (value << 1)
		cpu.setCReg(value)
	} else if operand1 == OPERAND_D && operand2 == OPERAND_NONE {
		value = cpu.getDReg()
		bit7 = value >> 7
		value = (value << 1)
		cpu.setDReg(value)
	} else if operand1 == OPERAND_E && operand2 == OPERAND_NONE {
		value = cpu.getEReg()
		bit7 = value >> 7
		value = (value << 1)
		cpu.setEReg(value)
	} else if operand1 == OPERAND_H && operand2 == OPERAND_NONE {
		value = cpu.getHReg()
		bit7 = value >> 7
		value = (value << 1)
		cpu.setHReg(value)
	} else if operand1 == OPERAND_L && operand2 == OPERAND_NONE {
		value = cpu.getLReg()
		bit7 = value >> 7
		value = (value << 1)
		cpu.setLReg(value)
	} else if operand1 == OPERAND_HL_PAREN && operand2 == OPERAND_NONE {
		value = cpu.FetchMemory8(cpu.Reg.HL)
		cpu.timer(1)
		bit7 = value >> 7
		value = (value << 1)
		cpu.SetMemory8(cpu.Reg.HL, value)
		cpu.timer(2)
	} else if operand1 == OPERAND_A && operand2 == OPERAND_NONE {
		value = cpu.getAReg()
		bit7 = value >> 7
		value = (value << 1)
		cpu.setAReg(value)
	} else {
		errMsg := fmt.Sprintf("Error: SLA %s %s", operand1, operand2)
		panic(errMsg)
	}

	cpu.flagZ(value)
	cpu.flagN(false)
	cpu.clearHFlag()
	if bit7 != 0 {
		cpu.setCFlag()
	} else {
		cpu.clearCFlag()
	}
	cpu.Reg.PC++
}

// SRA Shift Right MSBit dosen't change
func (cpu *CPU) SRA(operand1, operand2 int) {
	var value byte
	var bit0 byte
	if operand1 == OPERAND_B && operand2 == OPERAND_NONE {
		value = cpu.getBReg()
		bit0 = value % 2
		bit7 := value >> 7
		value = (value >> 1)
		if bit7 != 0 {
			value |= 0x80
		} else {
			value &= 0x7f
		}
		cpu.setBReg(value)
	} else if operand1 == OPERAND_C && operand2 == OPERAND_NONE {
		value = cpu.getCReg()
		bit0 = value % 2
		bit7 := value >> 7
		value = (value >> 1)
		if bit7 != 0 {
			value |= 0x80
		} else {
			value &= 0x7f
		}
		cpu.setCReg(value)
	} else if operand1 == OPERAND_D && operand2 == OPERAND_NONE {
		value = cpu.getDReg()
		bit0 = value % 2
		bit7 := value >> 7
		value = (value >> 1)
		if bit7 != 0 {
			value |= 0x80
		} else {
			value &= 0x7f
		}
		cpu.setDReg(value)
	} else if operand1 == OPERAND_E && operand2 == OPERAND_NONE {
		value = cpu.getEReg()
		bit0 = value % 2
		bit7 := value >> 7
		value = (value >> 1)
		if bit7 != 0 {
			value |= 0x80
		} else {
			value &= 0x7f
		}
		cpu.setEReg(value)
	} else if operand1 == OPERAND_H && operand2 == OPERAND_NONE {
		value = cpu.getHReg()
		bit0 = value % 2
		bit7 := value >> 7
		value = (value >> 1)
		if bit7 != 0 {
			value |= 0x80
		} else {
			value &= 0x7f
		}
		cpu.setHReg(value)
	} else if operand1 == OPERAND_L && operand2 == OPERAND_NONE {
		value = cpu.getLReg()
		bit0 = value % 2
		bit7 := value >> 7
		value = (value >> 1)
		if bit7 != 0 {
			value |= 0x80
		} else {
			value &= 0x7f
		}
		cpu.setLReg(value)
	} else if operand1 == OPERAND_HL_PAREN && operand2 == OPERAND_NONE {
		value = cpu.FetchMemory8(cpu.Reg.HL)
		cpu.timer(1)
		bit0 = value % 2
		bit7 := value >> 7
		value = (value >> 1)
		if bit7 != 0 {
			value |= 0x80
		} else {
			value &= 0x7f
		}
		cpu.SetMemory8(cpu.Reg.HL, value)
		cpu.timer(2)
	} else if operand1 == OPERAND_A && operand2 == OPERAND_NONE {
		value = cpu.getAReg()
		bit0 = value % 2
		bit7 := value >> 7
		value = (value >> 1)
		if bit7 != 0 {
			value |= 0x80
		} else {
			value &= 0x7f
		}
		cpu.setAReg(value)
	} else {
		errMsg := fmt.Sprintf("Error: SRA %s %s", operand1, operand2)
		panic(errMsg)
	}

	cpu.flagZ(value)
	cpu.flagN(false)
	cpu.clearHFlag()
	if bit0 != 0 {
		cpu.setCFlag()
	} else {
		cpu.clearCFlag()
	}
	cpu.Reg.PC++
}

// SWAP Swap n[5:8] and n[0:4]
func (cpu *CPU) SWAP(operand1, operand2 int) {
	var value byte

	switch operand1 {
	case OPERAND_B:
		B := cpu.getBReg()
		B03 := B & 0x0f
		B47 := B >> 4
		value = (B03 << 4) | B47
		cpu.setBReg(value)
	case OPERAND_C:
		C := cpu.getCReg()
		C03 := C & 0x0f
		C47 := C >> 4
		value = (C03 << 4) | C47
		cpu.setCReg(value)
	case OPERAND_D:
		D := cpu.getDReg()
		D03 := D & 0x0f
		D47 := D >> 4
		value = (D03 << 4) | D47
		cpu.setDReg(value)
	case OPERAND_E:
		E := cpu.getEReg()
		E03 := E & 0x0f
		E47 := E >> 4
		value = (E03 << 4) | E47
		cpu.setEReg(value)
	case OPERAND_H:
		H := cpu.getHReg()
		H03 := H & 0x0f
		H47 := H >> 4
		value = (H03 << 4) | H47
		cpu.setHReg(value)
	case OPERAND_L:
		L := cpu.getLReg()
		L03 := L & 0x0f
		L47 := L >> 4
		value = (L03 << 4) | L47
		cpu.setLReg(value)
	case OPERAND_HL_PAREN:
		data := cpu.FetchMemory8(cpu.Reg.HL)
		cpu.timer(1)
		data03 := data & 0x0f
		data47 := data >> 4
		value = (data03 << 4) | data47
		cpu.SetMemory8(cpu.Reg.HL, value)
		cpu.timer(2)
	case OPERAND_A:
		A := cpu.getAReg()
		A03 := A & 0x0f
		A47 := A >> 4
		value = (A03 << 4) | A47
		cpu.setAReg(value)
	default:
		errMsg := fmt.Sprintf("Error: SWAP %s %s", operand1, operand2)
		panic(errMsg)
	}

	cpu.flagZ(value)
	cpu.flagN(false)
	cpu.clearHFlag()
	cpu.clearCFlag()
	cpu.Reg.PC++
}

// SRL Shift Right MSBit = 0
func (cpu *CPU) SRL(operand1, operand2 int) {
	var value byte
	var bit0 byte

	switch operand1 {
	case OPERAND_B:
		value = cpu.getBReg()
		bit0 = value % 2
		value = (value >> 1)
		cpu.setBReg(value)
	case OPERAND_C:
		value = cpu.getCReg()
		bit0 = value % 2
		value = (value >> 1)
		cpu.setCReg(value)
	case OPERAND_D:
		value = cpu.getDReg()
		bit0 = value % 2
		value = (value >> 1)
		cpu.setDReg(value)
	case OPERAND_E:
		value = cpu.getEReg()
		bit0 = value % 2
		value = (value >> 1)
		cpu.setEReg(value)
	case OPERAND_H:
		value = cpu.getHReg()
		bit0 = value % 2
		value = (value >> 1)
		cpu.setHReg(value)
	case OPERAND_L:
		value = cpu.getLReg()
		bit0 = value % 2
		value = (value >> 1)
		cpu.setLReg(value)
	case OPERAND_HL_PAREN:
		value = cpu.FetchMemory8(cpu.Reg.HL)
		cpu.timer(1)
		bit0 = value % 2
		value = (value >> 1)
		cpu.SetMemory8(cpu.Reg.HL, value)
		cpu.timer(2)
	case OPERAND_A:
		value = cpu.getAReg()
		bit0 = value % 2
		value = (value >> 1)
		cpu.setAReg(value)
	default:
		errMsg := fmt.Sprintf("Error: SRL %s %s", operand1, operand2)
		panic(errMsg)
	}

	cpu.flagZ(value)
	cpu.flagN(false)
	cpu.clearHFlag()
	if bit0 == 1 {
		cpu.setCFlag()
	} else {
		cpu.clearCFlag()
	}
	cpu.Reg.PC++
}

// BIT Test bit n
func (cpu *CPU) BIT(operand1, operand2 int) {
	var value byte

	var targetBit uint // ターゲットのbit
	switch operand1 {
	case OPERAND_0:
		targetBit = 0
	case OPERAND_1:
		targetBit = 1
	case OPERAND_2:
		targetBit = 2
	case OPERAND_3:
		targetBit = 3
	case OPERAND_4:
		targetBit = 4
	case OPERAND_5:
		targetBit = 5
	case OPERAND_6:
		targetBit = 6
	case OPERAND_7:
		targetBit = 7
	}

	switch operand2 {
	case OPERAND_B:
		value = (cpu.getBReg() >> targetBit) % 2
	case OPERAND_C:
		value = (cpu.getCReg() >> targetBit) % 2
	case OPERAND_D:
		value = (cpu.getDReg() >> targetBit) % 2
	case OPERAND_E:
		value = (cpu.getEReg() >> targetBit) % 2
	case OPERAND_H:
		value = (cpu.getHReg() >> targetBit) % 2
	case OPERAND_L:
		value = (cpu.getLReg() >> targetBit) % 2
	case OPERAND_HL_PAREN:
		value = (cpu.FetchMemory8(cpu.Reg.HL) >> targetBit) % 2
	case OPERAND_A:
		value = (cpu.getAReg() >> targetBit) % 2
	}

	cpu.flagZ(value)
	cpu.flagN(false)
	cpu.setHFlag()
	cpu.Reg.PC++
}

// RES Clear bit n
func (cpu *CPU) RES(operand1, operand2 int) {

	var targetBit uint // ターゲットのbit
	switch operand1 {
	case OPERAND_0:
		targetBit = 0
	case OPERAND_1:
		targetBit = 1
	case OPERAND_2:
		targetBit = 2
	case OPERAND_3:
		targetBit = 3
	case OPERAND_4:
		targetBit = 4
	case OPERAND_5:
		targetBit = 5
	case OPERAND_6:
		targetBit = 6
	case OPERAND_7:
		targetBit = 7
	}

	switch operand2 {
	case OPERAND_B:
		mask := ^(byte(1) << targetBit)
		B := cpu.getBReg() & mask
		cpu.setBReg(B)
	case OPERAND_C:
		mask := ^(byte(1) << targetBit)
		C := cpu.getCReg() & mask
		cpu.setCReg(C)
	case OPERAND_D:
		mask := ^(byte(1) << targetBit)
		D := cpu.getDReg() & mask
		cpu.setDReg(D)
	case OPERAND_E:
		mask := ^(byte(1) << targetBit)
		E := cpu.getEReg() & mask
		cpu.setEReg(E)
	case OPERAND_H:
		mask := ^(byte(1) << targetBit)
		H := cpu.getHReg() & mask
		cpu.setHReg(H)
	case OPERAND_L:
		mask := ^(byte(1) << targetBit)
		L := cpu.getLReg() & mask
		cpu.setLReg(L)
	case OPERAND_HL_PAREN:
		mask := ^(byte(1) << targetBit)
		value := cpu.FetchMemory8(cpu.Reg.HL) & mask
		cpu.timer(1)
		cpu.SetMemory8(cpu.Reg.HL, value)
		cpu.timer(2)
	case OPERAND_A:
		mask := ^(byte(1) << targetBit)
		A := cpu.getAReg() & mask
		cpu.setAReg(A)
	}
	cpu.Reg.PC++
}

// SET Clear bit n
func (cpu *CPU) SET(operand1, operand2 int) {

	var targetBit uint // ターゲットのbit
	switch operand1 {
	case OPERAND_0:
		targetBit = 0
	case OPERAND_1:
		targetBit = 1
	case OPERAND_2:
		targetBit = 2
	case OPERAND_3:
		targetBit = 3
	case OPERAND_4:
		targetBit = 4
	case OPERAND_5:
		targetBit = 5
	case OPERAND_6:
		targetBit = 6
	case OPERAND_7:
		targetBit = 7
	}

	switch operand2 {
	case OPERAND_B:
		mask := byte(1) << targetBit
		B := cpu.getBReg() | mask
		cpu.setBReg(B)
	case OPERAND_C:
		mask := byte(1) << targetBit
		C := cpu.getCReg() | mask
		cpu.setCReg(C)
	case OPERAND_D:
		mask := byte(1) << targetBit
		D := cpu.getDReg() | mask
		cpu.setDReg(D)
	case OPERAND_E:
		mask := byte(1) << targetBit
		E := cpu.getEReg() | mask
		cpu.setEReg(E)
	case OPERAND_H:
		mask := byte(1) << targetBit
		H := cpu.getHReg() | mask
		cpu.setHReg(H)
	case OPERAND_L:
		mask := byte(1) << targetBit
		L := cpu.getLReg() | mask
		cpu.setLReg(L)
	case OPERAND_HL_PAREN:
		mask := byte(1) << targetBit
		value := cpu.FetchMemory8(cpu.Reg.HL) | mask
		cpu.timer(1)
		cpu.SetMemory8(cpu.Reg.HL, value)
		cpu.timer(2)
	case OPERAND_A:
		mask := byte(1) << targetBit
		A := cpu.getAReg() | mask
		cpu.setAReg(A)
	}
	cpu.Reg.PC++
}

// PUSH スタックにPUSH
func (cpu *CPU) PUSH(operand1, operand2 int) {
	cpu.timer(1)
	switch operand1 {
	case OPERAND_BC:
		cpu.pushBC()
	case OPERAND_DE:
		cpu.pushDE()
	case OPERAND_HL:
		cpu.pushHL()
	case OPERAND_AF:
		cpu.pushAF()
	default:
		errMsg := fmt.Sprintf("Error: PUSH %s %s", operand1, operand2)
		panic(errMsg)
	}
	cpu.Reg.PC++
	cpu.timer(2)
}

// POP スタックからPOP
func (cpu *CPU) POP(operand1, operand2 int) {
	switch operand1 {
	case OPERAND_BC:
		cpu.popBC()
	case OPERAND_DE:
		cpu.popDE()
	case OPERAND_HL:
		cpu.popHL()
	case OPERAND_AF:
		cpu.popAF()
	default:
		errMsg := fmt.Sprintf("Error: POP %s %s", operand1, operand2)
		panic(errMsg)
	}
	cpu.Reg.PC++
	cpu.timer(2)
}

// SUB 減算
func (cpu *CPU) SUB(operand1, operand2 int) {
	switch operand1 {
	case OPERAND_A:
		value := cpu.getAReg() - cpu.getAReg()
		carryBits := cpu.getAReg() ^ cpu.getAReg() ^ value
		cpu.flagC8Sub(cpu.getAReg(), cpu.getAReg())
		cpu.setAReg(value)
		cpu.flagZ(value)
		cpu.flagN(true)
		cpu.flagH8(carryBits)
		cpu.Reg.PC++
	case OPERAND_B:
		value := cpu.getAReg() - cpu.getBReg()
		carryBits := cpu.getAReg() ^ cpu.getBReg() ^ value
		cpu.flagC8Sub(cpu.getAReg(), cpu.getBReg())
		cpu.setAReg(value)
		cpu.flagZ(value)
		cpu.flagN(true)
		cpu.flagH8(carryBits)
		cpu.Reg.PC++
	case OPERAND_C:
		value := cpu.getAReg() - cpu.getCReg()
		carryBits := cpu.getAReg() ^ cpu.getCReg() ^ value
		cpu.flagC8Sub(cpu.getAReg(), cpu.getCReg())
		cpu.setAReg(value)
		cpu.flagZ(value)
		cpu.flagN(true)
		cpu.flagH8(carryBits)
		cpu.Reg.PC++
	case OPERAND_D:
		value := cpu.getAReg() - cpu.getDReg()
		carryBits := cpu.getAReg() ^ cpu.getDReg() ^ value
		cpu.flagC8Sub(cpu.getAReg(), cpu.getDReg())
		cpu.setAReg(value)
		cpu.flagZ(value)
		cpu.flagN(true)
		cpu.flagH8(carryBits)
		cpu.Reg.PC++
	case OPERAND_E:
		value := cpu.getAReg() - cpu.getEReg()
		carryBits := cpu.getAReg() ^ cpu.getEReg() ^ value
		cpu.flagC8Sub(cpu.getAReg(), cpu.getEReg())
		cpu.setAReg(value)
		cpu.flagZ(value)
		cpu.flagN(true)
		cpu.flagH8(carryBits)
		cpu.Reg.PC++
	case OPERAND_H:
		value := cpu.getAReg() - cpu.getHReg()
		carryBits := cpu.getAReg() ^ cpu.getHReg() ^ value
		cpu.flagC8Sub(cpu.getAReg(), cpu.getHReg())
		cpu.setAReg(value)
		cpu.flagZ(value)
		cpu.flagN(true)
		cpu.flagH8(carryBits)
		cpu.Reg.PC++
	case OPERAND_L:
		value := cpu.getAReg() - cpu.getLReg()
		carryBits := cpu.getAReg() ^ cpu.getLReg() ^ value
		cpu.flagC8Sub(cpu.getAReg(), cpu.getLReg())
		cpu.setAReg(value)
		cpu.flagZ(value)
		cpu.flagN(true)
		cpu.flagH8(carryBits)
		cpu.Reg.PC++
	case OPERAND_d8:
		value := cpu.getAReg() - cpu.d8Fetch()
		carryBits := cpu.getAReg() ^ cpu.d8Fetch() ^ value
		cpu.flagC8Sub(cpu.getAReg(), cpu.d8Fetch())
		cpu.setAReg(value)
		cpu.flagZ(value)
		cpu.flagN(true)
		cpu.flagH8(carryBits)
		cpu.Reg.PC += 2
	case OPERAND_HL_PAREN:
		value := cpu.getAReg() - cpu.FetchMemory8(cpu.Reg.HL)
		carryBits := cpu.getAReg() ^ cpu.FetchMemory8(cpu.Reg.HL) ^ value
		cpu.flagC8Sub(cpu.getAReg(), cpu.FetchMemory8(cpu.Reg.HL))
		cpu.setAReg(value)
		cpu.flagZ(value)
		cpu.flagN(true)
		cpu.flagH8(carryBits)
		cpu.Reg.PC++
	default:
		errMsg := fmt.Sprintf("Error: SUB %s %s", operand1, operand2)
		panic(errMsg)
	}
}

// RRA Rotate register A right through carry.
func (cpu *CPU) RRA(operand1, operand2 int) {
	carry := cpu.getCFlag()
	A := cpu.getAReg()
	A0 := A % 2
	if A0 == 1 {
		cpu.setCFlag()
	} else {
		cpu.clearCFlag()
	}
	if carry {
		A = (1 << 7) | (A >> 1)
	} else {
		A = (0 << 7) | (A >> 1)
	}
	cpu.setAReg(A)
	cpu.clearZFlag()
	cpu.flagN(false)
	cpu.clearHFlag()
	cpu.Reg.PC++
}

// ADC Add the value n8 plus the carry flag to A
func (cpu *CPU) ADC(operand1, operand2 int) {
	var carry, value, value4 byte
	var value16 uint16
	if cpu.getCFlag() {
		carry = 1
	} else {
		carry = 0
	}

	switch operand1 {
	case OPERAND_A:
		switch operand2 {
		case OPERAND_A:
			value = cpu.getAReg() + carry + cpu.getAReg()
			value4 = cpu.getARegLower4() + carry + cpu.getARegLower4()
			value16 = uint16(cpu.getAReg()) + uint16(cpu.getAReg()) + uint16(carry)
		case OPERAND_B:
			value = cpu.getBReg() + carry + cpu.getAReg()
			value4 = cpu.getBRegLower4() + carry + cpu.getARegLower4()
			value16 = uint16(cpu.getBReg()) + uint16(carry) + uint16(cpu.getAReg())
		case OPERAND_C:
			value = cpu.getCReg() + carry + cpu.getAReg()
			value4 = cpu.getCRegLower4() + carry + cpu.getARegLower4()
			value16 = uint16(cpu.getCReg()) + uint16(carry) + uint16(cpu.getAReg())
		case OPERAND_D:
			value = cpu.getDReg() + carry + cpu.getAReg()
			value4 = cpu.getDRegLower4() + carry + cpu.getARegLower4()
			value16 = uint16(cpu.getDReg()) + uint16(carry) + uint16(cpu.getAReg())
		case OPERAND_E:
			value = cpu.getEReg() + carry + cpu.getAReg()
			value4 = cpu.getERegLower4() + carry + cpu.getARegLower4()
			value16 = uint16(cpu.getEReg()) + uint16(carry) + uint16(cpu.getAReg())
		case OPERAND_H:
			value = cpu.getHReg() + carry + cpu.getAReg()
			value4 = cpu.getHRegLower4() + carry + cpu.getARegLower4()
			value16 = uint16(cpu.getHReg()) + uint16(carry) + uint16(cpu.getAReg())
		case OPERAND_L:
			value = cpu.getLReg() + carry + cpu.getAReg()
			value4 = cpu.getLRegLower4() + carry + cpu.getARegLower4()
			value16 = uint16(cpu.getLReg()) + uint16(carry) + uint16(cpu.getAReg())
		case OPERAND_HL_PAREN:
			data := cpu.FetchMemory8(cpu.Reg.HL)
			value = data + carry + cpu.getAReg()
			value4 = (data & 0x0f) + carry + cpu.getARegLower4()
			value16 = uint16(data) + uint16(cpu.getAReg()) + uint16(carry)
		case OPERAND_d8:
			data := cpu.d8Fetch()
			value = data + carry + cpu.getAReg()
			value4 = (data & 0x0f) + carry + cpu.getARegLower4()
			value16 = uint16(data) + uint16(cpu.getAReg()) + uint16(carry)
			cpu.Reg.PC++
		}
	default:
		errMsg := fmt.Sprintf("Error: ADC %s %s", operand1, operand2)
		panic(errMsg)
	}
	cpu.setAReg(value)
	cpu.flagZ(value)
	cpu.flagN(false)
	cpu.flagH8(value4)
	cpu.flagC8(value16)
	cpu.Reg.PC++
}

// SBC Subtract the value n8 and the carry flag from A
func (cpu *CPU) SBC(operand1, operand2 int) {
	var carry, value, value4 byte
	var value16 uint16
	if cpu.getCFlag() {
		carry = 1
	} else {
		carry = 0
	}

	switch operand1 {
	case OPERAND_A:
		switch operand2 {
		case OPERAND_A:
			value = cpu.getAReg() - (cpu.getAReg() + carry)
			value4 = cpu.getARegLower4() - (cpu.getARegLower4() + carry)
			value16 = uint16(cpu.getAReg()) - (uint16(cpu.getAReg()) + uint16(carry))
		case OPERAND_B:
			value = cpu.getAReg() - (cpu.getBReg() + carry)
			value4 = cpu.getARegLower4() - (cpu.getBRegLower4() + carry)
			value16 = uint16(cpu.getAReg()) - (uint16(cpu.getBReg()) + uint16(carry))
		case OPERAND_C:
			value = cpu.getAReg() - (cpu.getCReg() + carry)
			value4 = cpu.getARegLower4() - (cpu.getCRegLower4() + carry)
			value16 = uint16(cpu.getAReg()) - (uint16(cpu.getCReg()) + uint16(carry))
		case OPERAND_D:
			value = cpu.getAReg() - (cpu.getDReg() + carry)
			value4 = cpu.getARegLower4() - (cpu.getDRegLower4() + carry)
			value16 = uint16(cpu.getAReg()) - (uint16(cpu.getDReg()) + uint16(carry))
		case OPERAND_E:
			value = cpu.getAReg() - (cpu.getEReg() + carry)
			value4 = cpu.getARegLower4() - (cpu.getERegLower4() + carry)
			value16 = uint16(cpu.getAReg()) - (uint16(cpu.getEReg()) + uint16(carry))
		case OPERAND_H:
			value = cpu.getAReg() - (cpu.getHReg() + carry)
			value4 = cpu.getARegLower4() - (cpu.getHRegLower4() + carry)
			value16 = uint16(cpu.getAReg()) - (uint16(cpu.getHReg()) + uint16(carry))
		case OPERAND_L:
			value = cpu.getAReg() - (cpu.getLReg() + carry)
			value4 = cpu.getARegLower4() - (cpu.getLRegLower4() + carry)
			value16 = uint16(cpu.getAReg()) - (uint16(cpu.getLReg()) + uint16(carry))
		case OPERAND_HL_PAREN:
			data := cpu.FetchMemory8(cpu.Reg.HL)
			value = cpu.getAReg() - (data + carry)
			value4 = cpu.getARegLower4() - ((data & 0x0f) + carry)
			value16 = uint16(cpu.getAReg()) - (uint16(data) + uint16(carry))
		case OPERAND_d8:
			data := cpu.d8Fetch()
			value = cpu.getAReg() - (data + carry)
			value4 = cpu.getARegLower4() - ((data & 0x0f) + carry)
			value16 = uint16(cpu.getAReg()) - (uint16(data) + uint16(carry))
			cpu.Reg.PC++
		}
	default:
		errMsg := fmt.Sprintf("Error: SBC %s %s", operand1, operand2)
		panic(errMsg)
	}
	cpu.setAReg(value)
	cpu.flagZ(value)
	cpu.flagN(true)
	cpu.flagH8(value4)
	cpu.flagC8(value16)
	cpu.Reg.PC++
}

// DAA Decimal adjust
func (cpu *CPU) DAA(operand1, operand2 int) {
	A := uint8(cpu.getAReg())
	// ref: https://forums.nesdev.com/viewtopic.php?f=20&t=15944
	if !cpu.getNFlag() {
		if cpu.getCFlag() || A > 0x99 {
			A += 0x60
			cpu.setCFlag()
		}
		if cpu.getHFlag() || (A&0x0f) > 0x09 {
			A += 0x06
		}
	} else {
		if cpu.getCFlag() {
			A -= 0x60
		}
		if cpu.getHFlag() {
			A -= 0x06
		}
	}

	cpu.setAReg(A)
	cpu.flagZ(A)
	cpu.clearHFlag()
	cpu.Reg.PC++
}

// RST Push present address and jump to vector address
func (cpu *CPU) RST(operand1, operand2 int) {

	var vector uint16
	switch operand1 {
	case OPERAND_00H:
		vector = 0x00
	case OPERAND_08H:
		vector = 0x08
	case OPERAND_10H:
		vector = 0x10
	case OPERAND_18H:
		vector = 0x18
	case OPERAND_20H:
		vector = 0x20
	case OPERAND_28H:
		vector = 0x28
	case OPERAND_30H:
		vector = 0x30
	case OPERAND_38H:
		vector = 0x38
	}

	destination := uint16(vector)
	cpu.Reg.PC++
	cpu.pushPC()
	cpu.Reg.PC = destination
}

// SCF Set Carry Flag
func (cpu *CPU) SCF(operand1, operand2 int) {
	cpu.flagN(false)
	cpu.clearHFlag()
	cpu.setCFlag()
	cpu.Reg.PC++
}

// CCF Complement Carry Flag
func (cpu *CPU) CCF(operand1, operand2 int) {
	cpu.flagN(false)
	cpu.clearHFlag()
	if cpu.getCFlag() {
		cpu.clearCFlag()
	} else {
		cpu.setCFlag()
	}
	cpu.Reg.PC++
}
