package emulator

import (
	"fmt"
	"strconv"
)

func (cpu *CPU) a16Fetch() uint16 {
	value := cpu.d16Fetch()
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

// LD Load
func (cpu *CPU) LD(operand1, operand2 string) {
	switch operand1 {
	case "A":
		switch operand2 {
		case "A":
			cpu.setAReg(cpu.getAReg())
			cpu.Reg.PC++
		case "B":
			cpu.setAReg(cpu.getBReg())
			cpu.Reg.PC++
		case "C":
			cpu.setAReg(cpu.getCReg())
			cpu.Reg.PC++
		case "D":
			cpu.setAReg(cpu.getDReg())
			cpu.Reg.PC++
		case "E":
			cpu.setAReg(cpu.getEReg())
			cpu.Reg.PC++
		case "H":
			cpu.setAReg(cpu.getHReg())
			cpu.Reg.PC++
		case "L":
			cpu.setAReg(cpu.getLReg())
			cpu.Reg.PC++
		case "d8":
			cpu.setAReg(cpu.FetchMemory8(cpu.Reg.PC + 1))
			cpu.Reg.PC += 2
		case "(C)":
			addr := 0xff00 + uint16(cpu.getCReg())
			cpu.setAReg(cpu.FetchMemory8(addr))
			cpu.Reg.PC++ // 誤植(https://www.pastraiser.com/cpu/gameboy/gameboy_opcodes.html)
		case "(BC)":
			cpu.setAReg(cpu.FetchMemory8(cpu.Reg.BC))
			cpu.Reg.PC++
		case "(DE)":
			cpu.setAReg(cpu.FetchMemory8(cpu.Reg.DE))
			cpu.Reg.PC++
		case "(HL)":
			value := cpu.FetchMemory8(cpu.Reg.HL)
			cpu.setAReg(value)
			cpu.Reg.PC++
		case "(HL+)":
			cpu.setAReg(cpu.FetchMemory8(cpu.Reg.HL))
			cpu.Reg.HL++
			cpu.Reg.PC++
		case "(HL-)":
			cpu.setAReg(cpu.FetchMemory8(cpu.Reg.HL))
			cpu.Reg.HL--
			cpu.Reg.PC++
		case "(a16)":
			addr := cpu.a16Fetch()
			cpu.setAReg(cpu.FetchMemory8(addr))
			cpu.Reg.PC += 3
		}
	case "B":
		switch operand2 {
		case "A":
			cpu.setBReg(cpu.getAReg())
			cpu.Reg.PC++
		case "B":
			cpu.setBReg(cpu.getBReg())
			cpu.Reg.PC++
		case "C":
			cpu.setBReg(cpu.getCReg())
			cpu.Reg.PC++
		case "D":
			cpu.setBReg(cpu.getDReg())
			cpu.Reg.PC++
		case "E":
			cpu.setBReg(cpu.getEReg())
			cpu.Reg.PC++
		case "H":
			cpu.setBReg(cpu.getHReg())
			cpu.Reg.PC++
		case "L":
			cpu.setBReg(cpu.getLReg())
			cpu.Reg.PC++
		case "d8":
			value := cpu.d8Fetch()
			cpu.setBReg(value)
			cpu.Reg.PC += 2
		case "(HL)":
			value := cpu.FetchMemory8(cpu.Reg.HL)
			cpu.setBReg(value)
			cpu.Reg.PC++
		}
	case "C":
		switch operand2 {
		case "A":
			cpu.setCReg(cpu.getAReg())
			cpu.Reg.PC++
		case "B":
			cpu.setCReg(cpu.getBReg())
			cpu.Reg.PC++
		case "C":
			cpu.setCReg(cpu.getCReg())
			cpu.Reg.PC++
		case "D":
			cpu.setCReg(cpu.getDReg())
			cpu.Reg.PC++
		case "E":
			cpu.setCReg(cpu.getEReg())
			cpu.Reg.PC++
		case "H":
			cpu.setCReg(cpu.getHReg())
			cpu.Reg.PC++
		case "L":
			cpu.setCReg(cpu.getLReg())
			cpu.Reg.PC++
		case "d8":
			value := cpu.d8Fetch()
			cpu.setCReg(value)
			cpu.Reg.PC += 2
		case "(HL)":
			value := cpu.FetchMemory8(cpu.Reg.HL)
			cpu.setCReg(value)
			cpu.Reg.PC++
		}
	case "D":
		switch operand2 {
		case "A":
			cpu.setDReg(cpu.getAReg())
			cpu.Reg.PC++
		case "B":
			cpu.setDReg(cpu.getBReg())
			cpu.Reg.PC++
		case "C":
			cpu.setDReg(cpu.getCReg())
			cpu.Reg.PC++
		case "D":
			cpu.setDReg(cpu.getDReg())
			cpu.Reg.PC++
		case "E":
			cpu.setDReg(cpu.getEReg())
			cpu.Reg.PC++
		case "H":
			cpu.setDReg(cpu.getHReg())
			cpu.Reg.PC++
		case "L":
			cpu.setDReg(cpu.getLReg())
			cpu.Reg.PC++
		case "d8":
			value := cpu.d8Fetch()
			cpu.setDReg(value)
			cpu.Reg.PC += 2
		case "(HL)":
			value := cpu.FetchMemory8(cpu.Reg.HL)
			cpu.setDReg(value)
			cpu.Reg.PC++
		}
	case "E":
		switch operand2 {
		case "A":
			cpu.setEReg(cpu.getAReg())
			cpu.Reg.PC++
		case "B":
			cpu.setEReg(cpu.getBReg())
			cpu.Reg.PC++
		case "C":
			cpu.setEReg(cpu.getCReg())
			cpu.Reg.PC++
		case "D":
			cpu.setEReg(cpu.getDReg())
			cpu.Reg.PC++
		case "E":
			cpu.setEReg(cpu.getEReg())
			cpu.Reg.PC++
		case "H":
			cpu.setEReg(cpu.getHReg())
			cpu.Reg.PC++
		case "L":
			cpu.setEReg(cpu.getLReg())
			cpu.Reg.PC++
		case "d8":
			value := cpu.d8Fetch()
			cpu.setEReg(value)
			cpu.Reg.PC += 2
		case "(HL)":
			value := cpu.FetchMemory8(cpu.Reg.HL)
			cpu.setEReg(value)
			cpu.Reg.PC++
		}
	case "H":
		switch operand2 {
		case "A":
			cpu.setHReg(cpu.getAReg())
			cpu.Reg.PC++
		case "B":
			cpu.setHReg(cpu.getBReg())
			cpu.Reg.PC++
		case "C":
			cpu.setHReg(cpu.getCReg())
			cpu.Reg.PC++
		case "D":
			cpu.setHReg(cpu.getDReg())
			cpu.Reg.PC++
		case "E":
			cpu.setHReg(cpu.getEReg())
			cpu.Reg.PC++
		case "H":
			cpu.setHReg(cpu.getHReg())
			cpu.Reg.PC++
		case "L":
			cpu.setHReg(cpu.getLReg())
			cpu.Reg.PC++
		case "d8":
			value := cpu.d8Fetch()
			cpu.setHReg(value)
			cpu.Reg.PC += 2
		case "(HL)":
			value := cpu.FetchMemory8(cpu.Reg.HL)
			cpu.setHReg(value)
			cpu.Reg.PC++
		}
	case "L":
		switch operand2 {
		case "A":
			cpu.setLReg(cpu.getAReg())
			cpu.Reg.PC++
		case "B":
			cpu.setLReg(cpu.getBReg())
			cpu.Reg.PC++
		case "C":
			cpu.setLReg(cpu.getCReg())
			cpu.Reg.PC++
		case "D":
			cpu.setLReg(cpu.getDReg())
			cpu.Reg.PC++
		case "E":
			cpu.setLReg(cpu.getEReg())
			cpu.Reg.PC++
		case "H":
			cpu.setLReg(cpu.getHReg())
			cpu.Reg.PC++
		case "L":
			cpu.setLReg(cpu.getLReg())
			cpu.Reg.PC++
		case "d8":
			value := cpu.d8Fetch()
			cpu.setLReg(value)
			cpu.Reg.PC += 2
		case "(HL)":
			value := cpu.FetchMemory8(cpu.Reg.HL)
			cpu.setLReg(value)
			cpu.Reg.PC++
		}
	case "BC":
		switch operand2 {
		case "d16":
			value := cpu.d16Fetch()
			cpu.Reg.BC = value
			cpu.Reg.PC += 3
		}
	case "DE":
		switch operand2 {
		case "d16":
			value := cpu.d16Fetch()
			cpu.Reg.DE = value
			cpu.Reg.PC += 3
		}
	case "HL":
		switch operand2 {
		case "d16":
			cpu.Reg.HL = cpu.d16Fetch()
			cpu.Reg.PC += 3
		case "SP+r8":
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
	case "SP":
		switch operand2 {
		case "d16":
			value := cpu.d16Fetch()
			cpu.Reg.SP = value
			cpu.Reg.PC += 3
		case "HL":
			cpu.Reg.SP = cpu.Reg.HL
			cpu.Reg.PC++
		}
	case "(C)":
		switch operand2 {
		case "A":
			addr := 0xff00 + uint16(cpu.getCReg())
			cpu.SetMemory8(addr, cpu.getAReg())
			cpu.Reg.PC++ // 誤植(https://www.pastraiser.com/cpu/gameboy/gameboy_opcodes.html)
		}
	case "(BC)":
		switch operand2 {
		case "A":
			cpu.SetMemory8(cpu.Reg.BC, cpu.getAReg())
			cpu.Reg.PC++
		}
	case "(DE)":
		switch operand2 {
		case "A":
			cpu.SetMemory8(cpu.Reg.DE, cpu.getAReg())
			cpu.Reg.PC++
		}
	case "(HL)":
		switch operand2 {
		case "A":
			cpu.SetMemory8(cpu.Reg.HL, cpu.getAReg())
			cpu.Reg.PC++
		case "B":
			cpu.SetMemory8(cpu.Reg.HL, cpu.getBReg())
			cpu.Reg.PC++
		case "C":
			cpu.SetMemory8(cpu.Reg.HL, cpu.getCReg())
			cpu.Reg.PC++
		case "D":
			cpu.SetMemory8(cpu.Reg.HL, cpu.getDReg())
			cpu.Reg.PC++
		case "E":
			cpu.SetMemory8(cpu.Reg.HL, cpu.getEReg())
			cpu.Reg.PC++
		case "H":
			cpu.SetMemory8(cpu.Reg.HL, cpu.getHReg())
			cpu.Reg.PC++
		case "L":
			cpu.SetMemory8(cpu.Reg.HL, cpu.getLReg())
			cpu.Reg.PC++
		case "d8":
			value := cpu.d8Fetch()
			cpu.SetMemory8(cpu.Reg.HL, value)
			cpu.Reg.PC += 2
		}
	case "(HL+)":
		switch operand2 {
		case "A":
			cpu.SetMemory8(cpu.Reg.HL, cpu.getAReg())
			cpu.Reg.HL++
			cpu.Reg.PC++
		}
	case "(HL-)":
		switch operand2 {
		case "A":
			// (HL)=A, HL=HL-1
			cpu.SetMemory8(cpu.Reg.HL, cpu.getAReg())
			cpu.Reg.HL--
			cpu.Reg.PC++
		}
	case "(a16)":
		switch operand2 {
		case "A":
			addr := cpu.a16Fetch()
			cpu.SetMemory8(addr, cpu.getAReg())
			cpu.Reg.PC += 3
		case "SP":
			// Store SP into addresses n16 (LSB) and n16 + 1 (MSB).
			addr := cpu.a16Fetch()
			upper := byte(cpu.Reg.SP >> 8)     // MSB
			lower := byte(cpu.Reg.SP & 0x00ff) // LSB
			cpu.SetMemory8(addr, lower)
			cpu.SetMemory8(addr+1, upper)
			cpu.Reg.PC += 3
		}
	default:
		errMsg := fmt.Sprintf("Error: LD %s %s", operand1, operand2)
		panic(errMsg)
	}
}

// LDH Load High Byte
func (cpu *CPU) LDH(operand1, operand2 string) {
	if operand1 == "A" && operand2 == "(a8)" {
		// LD A,($FF00+a8)
		addr := 0xff00 + uint16(cpu.FetchMemory8(cpu.Reg.PC+1))
		value := cpu.FetchMemory8(addr)

		if addr == JOYPADIO {
			// Joypad読み込み
			value = cpu.formatJoypad()
		}

		cpu.setAReg(value)
		cpu.Reg.PC += 2
	} else if operand1 == "(a8)" && operand2 == "A" {
		// LD ($FF00+a8),A
		addr := 0xff00 + uint16(cpu.FetchMemory8(cpu.Reg.PC+1))

		cpu.SetMemory8(addr, cpu.getAReg())
		cpu.Reg.PC += 2
	} else {
		errMsg := fmt.Sprintf("Error: LDH %s %s", operand1, operand2)
		panic(errMsg)
	}
}

// NOP No operation
func (cpu *CPU) NOP(operand1, operand2 string) {
	if operand1 == "*" && operand2 == "*" {
		cpu.Reg.PC++
	} else {
		errMsg := fmt.Sprintf("Error: NOP %s %s", operand1, operand2)
		panic(errMsg)
	}
}

// INC Increment
func (cpu *CPU) INC(operand1, operand2 string) {
	var value, carryBits byte

	switch operand1 {
	case "A":
		value = cpu.getAReg() + 1
		carryBits = cpu.getAReg() ^ 1 ^ value
		cpu.setAReg(value)
	case "B":
		value = cpu.getBReg() + 1
		carryBits = cpu.getBReg() ^ 1 ^ value
		cpu.setBReg(value)
	case "C":
		value = cpu.getCReg() + 1
		carryBits = cpu.getCReg() ^ 1 ^ value
		cpu.setCReg(value)
	case "D":
		value = cpu.getDReg() + 1
		carryBits = cpu.getDReg() ^ 1 ^ value
		cpu.setDReg(value)
	case "E":
		value = cpu.getEReg() + 1
		carryBits = cpu.getEReg() ^ 1 ^ value
		cpu.setEReg(value)
	case "H":
		value = cpu.getHReg() + 1
		carryBits = cpu.getHReg() ^ 1 ^ value
		cpu.setHReg(value)
	case "L":
		value = cpu.getLReg() + 1
		carryBits = cpu.getLReg() ^ 1 ^ value
		cpu.setLReg(value)
	case "(HL)":
		value = cpu.FetchMemory8(cpu.Reg.HL) + 1
		carryBits = cpu.FetchMemory8(cpu.Reg.HL) ^ 1 ^ value
		cpu.SetMemory8(cpu.Reg.HL, value)
	case "BC":
		cpu.Reg.BC++
	case "DE":
		cpu.Reg.DE++
	case "HL":
		cpu.Reg.HL++
	case "SP":
		cpu.Reg.SP++
	default:
		errMsg := fmt.Sprintf("Error: INC %s %s", operand1, operand2)
		panic(errMsg)
	}

	if operand1 != "BC" && operand1 != "DE" && operand1 != "HL" && operand1 != "SP" {
		cpu.flagZ(value)
		cpu.flagN(false)
		cpu.flagH8(carryBits)
	}
	cpu.Reg.PC++
}

// DEC Decrement
func (cpu *CPU) DEC(operand1, operand2 string) {
	var value byte
	var carryBits byte

	switch operand1 {
	case "A":
		value = cpu.getAReg() - 1
		carryBits = cpu.getAReg() ^ 1 ^ value
		cpu.setAReg(value)
	case "B":
		value = cpu.getBReg() - 1
		carryBits = cpu.getBReg() ^ 1 ^ value
		cpu.setBReg(value)
	case "C":
		value = cpu.getCReg() - 1
		carryBits = cpu.getCReg() ^ 1 ^ value
		cpu.setCReg(value)
	case "D":
		value = cpu.getDReg() - 1
		carryBits = cpu.getDReg() ^ 1 ^ value
		cpu.setDReg(value)
	case "E":
		value = cpu.getEReg() - 1
		carryBits = cpu.getEReg() ^ 1 ^ value
		cpu.setEReg(value)
	case "H":
		value = cpu.getHReg() - 1
		carryBits = cpu.getHReg() ^ 1 ^ value
		cpu.setHReg(value)
	case "L":
		value = cpu.getLReg() - 1
		carryBits = cpu.getLReg() ^ 1 ^ value
		cpu.setLReg(value)
	case "(HL)":
		value = cpu.FetchMemory8(cpu.Reg.HL) - 1
		carryBits = cpu.FetchMemory8(cpu.Reg.HL) ^ 1 ^ value
		cpu.SetMemory8(cpu.Reg.HL, value)
	case "BC":
		cpu.Reg.BC--
	case "DE":
		cpu.Reg.DE--
	case "HL":
		cpu.Reg.HL--
	case "SP":
		cpu.Reg.SP--
	default:
		errMsg := fmt.Sprintf("Error: DEC %s %s", operand1, operand2)
		panic(errMsg)
	}

	if operand1 != "BC" && operand1 != "DE" && operand1 != "HL" && operand1 != "SP" {
		cpu.flagZ(value)
		cpu.flagN(true)
		cpu.flagH8(carryBits)
	}
	cpu.Reg.PC++
}

// JR Jump relatively
func (cpu *CPU) JR(operand1, operand2 string) {
	switch operand1 {
	case "r8":
		delta := int8(cpu.FetchMemory8(cpu.Reg.PC + 1))
		destination := uint16(int32(cpu.Reg.PC+2) + int32(delta)) // この時点でのPCは命令フェッチ後のPCなので+2してあげる必要あり
		cpu.Reg.PC = destination
	case "Z":
		if cpu.getZFlag() {
			delta := int8(cpu.FetchMemory8(cpu.Reg.PC + 1))
			destination := uint16(int32(cpu.Reg.PC+2) + int32(delta)) // この時点でのPCは命令フェッチ後のPCなので+2してあげる必要あり
			cpu.Reg.PC = destination
		} else {
			cpu.Reg.PC += 2
		}
	case "C":
		if cpu.getCFlag() {
			delta := int8(cpu.FetchMemory8(cpu.Reg.PC + 1))
			destination := uint16(int32(cpu.Reg.PC+2) + int32(delta)) // この時点でのPCは命令フェッチ後のPCなので+2してあげる必要あり
			cpu.Reg.PC = destination
		} else {
			cpu.Reg.PC += 2
		}
	case "NZ":
		if !cpu.getZFlag() {
			delta := int8(cpu.FetchMemory8(cpu.Reg.PC + 1))
			destination := uint16(int32(cpu.Reg.PC+2) + int32(delta)) // この時点でのPCは命令フェッチ後のPCなので+2してあげる必要あり
			cpu.Reg.PC = destination
		} else {
			cpu.Reg.PC += 2
		}
	case "NC":
		if !cpu.getCFlag() {
			delta := int8(cpu.FetchMemory8(cpu.Reg.PC + 1))
			destination := uint16(int32(cpu.Reg.PC+2) + int32(delta)) // この時点でのPCは命令フェッチ後のPCなので+2してあげる必要あり
			cpu.Reg.PC = destination
		} else {
			cpu.Reg.PC += 2
		}
	default:
		errMsg := fmt.Sprintf("Error: JR %s %s", operand1, operand2)
		panic(errMsg)
	}
}

// HALT Halt
func (cpu *CPU) HALT(operand1, operand2 string) {
	if operand1 == "*" && operand2 == "*" {
		if cpu.interruptTrigger {
			cpu.Reg.PC++
			cpu.interruptTrigger = false
		}
	} else {
		errMsg := fmt.Sprintf("Error: HALT %s %s", operand1, operand2)
		panic(errMsg)
	}
}

// STOP stop CPU
func (cpu *CPU) STOP(operand1, operand2 string) {
	if operand1 == "0" && operand2 == "*" {
		if cpu.interruptTrigger {
			cpu.Reg.PC += 2
			// 速度切り替え
			KEY1 := cpu.FetchMemory8(KEY1IO)
			if KEY1&0x01 == 1 {
				if KEY1>>7 == 1 {
					KEY1 = 0x00
				} else {
					KEY1 = 0x80
				}
				cpu.SetMemory8(KEY1IO, KEY1)
			}
			cpu.interruptTrigger = false
		}
	} else {
		errMsg := fmt.Sprintf("Error: STOP %s %s", operand1, operand2)
		panic(errMsg)
	}
}

// XOR xor
func (cpu *CPU) XOR(operand1, operand2 string) {
	var value byte
	switch operand1 {
	case "B":
		value = cpu.getAReg() ^ cpu.getBReg()
	case "C":
		value = cpu.getAReg() ^ cpu.getCReg()
	case "D":
		value = cpu.getAReg() ^ cpu.getDReg()
	case "E":
		value = cpu.getAReg() ^ cpu.getEReg()
	case "H":
		value = cpu.getAReg() ^ cpu.getHReg()
	case "L":
		value = cpu.getAReg() ^ cpu.getLReg()
	case "(HL)":
		value = cpu.getAReg() ^ cpu.FetchMemory8(cpu.Reg.HL)
	case "A":
		value = cpu.getAReg() ^ cpu.getAReg()
	case "d8":
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
func (cpu *CPU) JP(operand1, operand2 string) {
	switch operand1 {
	case "a16":
		destination := cpu.a16Fetch()
		cpu.Reg.PC = destination
	case "(HL)":
		destination := cpu.Reg.HL
		cpu.Reg.PC = destination
	case "Z":
		if cpu.getZFlag() {
			destination := cpu.a16Fetch()
			cpu.Reg.PC = destination
		} else {
			cpu.Reg.PC += 3
		}
	case "C":
		if cpu.getCFlag() {
			destination := cpu.a16Fetch()
			cpu.Reg.PC = destination
		} else {
			cpu.Reg.PC += 3
		}
	case "NZ":
		if !cpu.getZFlag() {
			destination := cpu.a16Fetch()
			cpu.Reg.PC = destination
		} else {
			cpu.Reg.PC += 3
		}
	case "NC":
		if !cpu.getCFlag() {
			destination := cpu.a16Fetch()
			cpu.Reg.PC = destination
		} else {
			cpu.Reg.PC += 3
		}
	default:
		errMsg := fmt.Sprintf("Error: JP %s %s", operand1, operand2)
		panic(errMsg)
	}
}

// RET Return
func (cpu *CPU) RET(operand1, operand2 string) {
	switch operand1 {
	case "*":
		// PC=(SP), SP=SP+2
		cpu.popPC()
	case "Z":
		if cpu.getZFlag() {
			cpu.popPC()
		} else {
			cpu.Reg.PC++
		}
	case "C":
		if cpu.getCFlag() {
			cpu.popPC()
		} else {
			cpu.Reg.PC++
		}
	case "NZ":
		if !cpu.getZFlag() {
			cpu.popPC()
		} else {
			cpu.Reg.PC++
		}
	case "NC":
		if !cpu.getCFlag() {
			cpu.popPC()
		} else {
			cpu.Reg.PC++
		}
	default:
		errMsg := fmt.Sprintf("Error: RET %s %s", operand1, operand2)
		panic(errMsg)
	}
}

// RETI Return Interrupt
func (cpu *CPU) RETI(operand1, operand2 string) {
	if operand1 == "*" && operand2 == "*" {
		cpu.popPC()
		cpu.Reg.IME = true
	} else {
		errMsg := fmt.Sprintf("Error: RETI %s %s", operand1, operand2)
		panic(errMsg)
	}
}

// CALL Call subroutine
func (cpu *CPU) CALL(operand1, operand2 string) {
	switch operand1 {
	case "a16":
		destination := cpu.a16Fetch()
		cpu.Reg.PC += 3
		cpu.pushPC()
		cpu.Reg.PC = destination
	case "Z":
		if cpu.getZFlag() {
			destination := cpu.a16Fetch()
			cpu.Reg.PC += 3
			cpu.pushPC()
			cpu.Reg.PC = destination
		} else {
			cpu.Reg.PC += 3
		}
	case "C":
		if cpu.getCFlag() {
			destination := cpu.a16Fetch()
			cpu.Reg.PC += 3
			cpu.pushPC()
			cpu.Reg.PC = destination
		} else {
			cpu.Reg.PC += 3
		}
	case "NZ":
		if !cpu.getZFlag() {
			destination := cpu.a16Fetch()
			cpu.Reg.PC += 3
			cpu.pushPC()
			cpu.Reg.PC = destination
		} else {
			cpu.Reg.PC += 3
		}
	case "NC":
		if !cpu.getCFlag() {
			destination := cpu.a16Fetch()
			cpu.Reg.PC += 3
			cpu.pushPC()
			cpu.Reg.PC = destination
		} else {
			cpu.Reg.PC += 3
		}
	default:
		errMsg := fmt.Sprintf("Error: CALL %s %s", operand1, operand2)
		panic(errMsg)
	}
}

// DI Disable Interrupt
func (cpu *CPU) DI(operand1, operand2 string) {
	if operand1 == "*" && operand2 == "*" {
		cpu.Reg.IME = false
		cpu.Reg.PC++
	} else {
		errMsg := fmt.Sprintf("Error: DI %s %s", operand1, operand2)
		panic(errMsg)
	}
}

// EI Enable Interrupt
func (cpu *CPU) EI(operand1, operand2 string) {
	if operand1 == "*" && operand2 == "*" {
		cpu.Reg.IME = true
		cpu.Reg.PC++
	} else {
		errMsg := fmt.Sprintf("Error: EI %s %s", operand1, operand2)
		panic(errMsg)
	}
}

// CP Compare
func (cpu *CPU) CP(operand1, operand2 string) {
	var value, carryBits byte

	switch operand1 {
	case "A":
		value = cpu.getAReg() - cpu.getAReg()
		carryBits = cpu.getAReg() ^ cpu.getAReg() ^ value
		cpu.flagC8Sub(cpu.getAReg(), cpu.getAReg())
	case "B":
		value = cpu.getAReg() - cpu.getBReg()
		carryBits = cpu.getAReg() ^ cpu.getBReg() ^ value
		cpu.flagC8Sub(cpu.getAReg(), cpu.getBReg())
	case "C":
		value = cpu.getAReg() - cpu.getCReg()
		carryBits = cpu.getAReg() ^ cpu.getCReg() ^ value
		cpu.flagC8Sub(cpu.getAReg(), cpu.getCReg())
	case "D":
		value = cpu.getAReg() - cpu.getDReg()
		carryBits = cpu.getAReg() ^ cpu.getDReg() ^ value
		cpu.flagC8Sub(cpu.getAReg(), cpu.getDReg())
	case "E":
		value = cpu.getAReg() - cpu.getEReg()
		carryBits = cpu.getAReg() ^ cpu.getEReg() ^ value
		cpu.flagC8Sub(cpu.getAReg(), cpu.getEReg())
	case "H":
		value = cpu.getAReg() - cpu.getHReg()
		carryBits = cpu.getAReg() ^ cpu.getHReg() ^ value
		cpu.flagC8Sub(cpu.getAReg(), cpu.getHReg())
	case "L":
		value = cpu.getAReg() - cpu.getLReg()
		carryBits = cpu.getAReg() ^ cpu.getLReg() ^ value
		cpu.flagC8Sub(cpu.getAReg(), cpu.getLReg())
	case "d8":
		value = cpu.getAReg() - cpu.d8Fetch()
		carryBits = cpu.getAReg() ^ cpu.d8Fetch() ^ value
		cpu.flagC8Sub(cpu.getAReg(), cpu.d8Fetch())
		cpu.Reg.PC++
	case "(HL)":
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
func (cpu *CPU) AND(operand1, operand2 string) {
	var value byte
	switch operand1 {
	case "A":
		value = cpu.getAReg() & cpu.getAReg()
	case "B":
		value = cpu.getAReg() & cpu.getBReg()
	case "C":
		value = cpu.getAReg() & cpu.getCReg()
	case "D":
		value = cpu.getAReg() & cpu.getDReg()
	case "E":
		value = cpu.getAReg() & cpu.getEReg()
	case "H":
		value = cpu.getAReg() & cpu.getHReg()
	case "L":
		value = cpu.getAReg() & cpu.getLReg()
	case "(HL)":
		value = cpu.getAReg() & cpu.FetchMemory8(cpu.Reg.HL)
	case "d8":
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
func (cpu *CPU) OR(operand1, operand2 string) {
	switch operand1 {
	case "A":
		value := cpu.getAReg() | cpu.getAReg()
		cpu.setAReg(value)
		cpu.flagZ(value)
	case "B":
		value := cpu.getAReg() | cpu.getBReg()
		cpu.setAReg(value)
		cpu.flagZ(value)
	case "C":
		value := cpu.getAReg() | cpu.getCReg()
		cpu.setAReg(value)
		cpu.flagZ(value)
	case "D":
		value := cpu.getAReg() | cpu.getDReg()
		cpu.setAReg(value)
		cpu.flagZ(value)
	case "E":
		value := cpu.getAReg() | cpu.getEReg()
		cpu.setAReg(value)
		cpu.flagZ(value)
	case "H":
		value := cpu.getAReg() | cpu.getHReg()
		cpu.setAReg(value)
		cpu.flagZ(value)
	case "L":
		value := cpu.getAReg() | cpu.getLReg()
		cpu.setAReg(value)
		cpu.flagZ(value)
	case "d8":
		value := cpu.getAReg() | cpu.FetchMemory8(cpu.Reg.PC+1)
		cpu.setAReg(value)
		cpu.flagZ(value)
		cpu.Reg.PC++
	case "(HL)":
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
func (cpu *CPU) ADD(operand1, operand2 string) {
	switch operand1 {
	case "A":
		switch operand2 {
		case "A":
			value := uint16(cpu.getAReg()) + uint16(cpu.getAReg())
			carryBits := uint16(cpu.getAReg()) ^ uint16(cpu.getAReg()) ^ value
			cpu.setAReg(byte(value))
			cpu.flagZ(byte(value))
			cpu.flagN(false)
			cpu.flagH8(byte(carryBits))
			cpu.flagC8(carryBits)
			cpu.Reg.PC++
		case "B":
			value := uint16(cpu.getAReg()) + uint16(cpu.getBReg())
			carryBits := uint16(cpu.getAReg()) ^ uint16(cpu.getBReg()) ^ value
			cpu.setAReg(byte(value))
			cpu.flagZ(byte(value))
			cpu.flagN(false)
			cpu.flagH8(byte(carryBits))
			cpu.flagC8(carryBits)
			cpu.Reg.PC++
		case "C":
			value := uint16(cpu.getAReg()) + uint16(cpu.getCReg())
			carryBits := uint16(cpu.getAReg()) ^ uint16(cpu.getCReg()) ^ value
			cpu.setAReg(byte(value))
			cpu.flagZ(byte(value))
			cpu.flagN(false)
			cpu.flagH8(byte(carryBits))
			cpu.flagC8(carryBits)
			cpu.Reg.PC++
		case "D":
			value := uint16(cpu.getAReg()) + uint16(cpu.getDReg())
			carryBits := uint16(cpu.getAReg()) ^ uint16(cpu.getDReg()) ^ value
			cpu.setAReg(byte(value))
			cpu.flagZ(byte(value))
			cpu.flagN(false)
			cpu.flagH8(byte(carryBits))
			cpu.flagC8(carryBits)
			cpu.Reg.PC++
		case "E":
			value := uint16(cpu.getAReg()) + uint16(cpu.getEReg())
			carryBits := uint16(cpu.getAReg()) ^ uint16(cpu.getEReg()) ^ value
			cpu.setAReg(byte(value))
			cpu.flagZ(byte(value))
			cpu.flagN(false)
			cpu.flagH8(byte(carryBits))
			cpu.flagC8(carryBits)
			cpu.Reg.PC++
		case "H":
			value := uint16(cpu.getAReg()) + uint16(cpu.getHReg())
			carryBits := uint16(cpu.getAReg()) ^ uint16(cpu.getHReg()) ^ value
			cpu.setAReg(byte(value))
			cpu.flagZ(byte(value))
			cpu.flagN(false)
			cpu.flagH8(byte(carryBits))
			cpu.flagC8(carryBits)
			cpu.Reg.PC++
		case "L":
			value := uint16(cpu.getAReg()) + uint16(cpu.getLReg())
			carryBits := uint16(cpu.getAReg()) ^ uint16(cpu.getLReg()) ^ value
			cpu.setAReg(byte(value))
			cpu.flagZ(byte(value))
			cpu.flagN(false)
			cpu.flagH8(byte(carryBits))
			cpu.flagC8(carryBits)
			cpu.Reg.PC++
		case "d8":
			value := uint16(cpu.getAReg()) + uint16(cpu.d8Fetch())
			carryBits := uint16(cpu.getAReg()) ^ uint16(cpu.d8Fetch()) ^ value
			cpu.setAReg(byte(value))
			cpu.flagZ(byte(value))
			cpu.flagN(false)
			cpu.flagH8(byte(carryBits))
			cpu.flagC8(carryBits)
			cpu.Reg.PC += 2
		case "(HL)":
			value := uint16(cpu.getAReg()) + uint16(cpu.FetchMemory8(cpu.Reg.HL))
			carryBits := uint16(cpu.getAReg()) ^ uint16(cpu.FetchMemory8(cpu.Reg.HL)) ^ value
			cpu.setAReg(byte(value))
			cpu.flagZ(byte(value))
			cpu.flagN(false)
			cpu.flagH8(byte(carryBits))
			cpu.flagC8(carryBits)
			cpu.Reg.PC++
		}
	case "HL":
		switch operand2 {
		case "BC":
			value := uint32(cpu.Reg.HL) + uint32(cpu.Reg.BC)
			carryBits := uint32(cpu.Reg.HL) ^ uint32(cpu.Reg.BC) ^ value
			cpu.Reg.HL = uint16(value)
			cpu.flagN(false)
			cpu.flagH16(uint16(carryBits))
			cpu.flagC16(carryBits)
			cpu.Reg.PC++
		case "DE":
			value := uint32(cpu.Reg.HL) + uint32(cpu.Reg.DE)
			carryBits := uint32(cpu.Reg.HL) ^ uint32(cpu.Reg.DE) ^ value
			cpu.Reg.HL = uint16(value)
			cpu.flagN(false)
			cpu.flagH16(uint16(carryBits))
			cpu.flagC16(carryBits)
			cpu.Reg.PC++
		case "HL":
			value := uint32(cpu.Reg.HL) + uint32(cpu.Reg.HL)
			carryBits := uint32(cpu.Reg.HL) ^ uint32(cpu.Reg.HL) ^ value
			cpu.Reg.HL = uint16(value)
			cpu.flagN(false)
			cpu.flagH16(uint16(carryBits))
			cpu.flagC16(carryBits)
			cpu.Reg.PC++
		case "SP":
			value := uint32(cpu.Reg.HL) + uint32(cpu.Reg.SP)
			carryBits := uint32(cpu.Reg.HL) ^ uint32(cpu.Reg.SP) ^ value
			cpu.Reg.HL = uint16(value)
			cpu.flagN(false)
			cpu.flagH16(uint16(carryBits))
			cpu.flagC16(carryBits)
			cpu.Reg.PC++
		}
	case "SP":
		switch operand2 {
		case "r8":
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
func (cpu *CPU) CPL(operand1, operand2 string) {
	if operand1 == "*" && operand2 == "*" {
		A := ^cpu.getAReg()
		cpu.setAReg(A)
		cpu.flagN(true)
		cpu.setHFlag()
		cpu.Reg.PC++
	} else {
		errMsg := fmt.Sprintf("Error: CPL %s %s", operand1, operand2)
		panic(errMsg)
	}
}

// PREFIXCB 拡張命令
func (cpu *CPU) PREFIXCB(operand1, operand2 string) {
	if operand1 == "*" && operand2 == "*" {
		cpu.Reg.PC++
		opcode := cpu.FetchMemory8(cpu.Reg.PC)
		instruction, operand1, operand2 := prefixCB[opcode][0], prefixCB[opcode][1], prefixCB[opcode][2]

		// cpu.pushHistory(cpu.Reg.PC, opcode, instruction, operand1, operand2)

		switch instruction {
		case "RLC":
			cpu.RLC(operand1, operand2)
		case "RRC":
			cpu.RRC(operand1, operand2)
		case "RL":
			cpu.RL(operand1, operand2)
		case "RR":
			cpu.RR(operand1, operand2)
		case "SLA":
			cpu.SLA(operand1, operand2)
		case "SRA":
			cpu.SRA(operand1, operand2)
		case "SWAP":
			cpu.SWAP(operand1, operand2)
		case "SRL":
			cpu.SRL(operand1, operand2)
		case "BIT":
			cpu.BIT(operand1, operand2)
		case "RES":
			cpu.RES(operand1, operand2)
		case "SET":
			cpu.SET(operand1, operand2)
		default:
			errMsg := fmt.Sprintf("eip: 0x%04x opcode: 0x%02x", cpu.Reg.PC, opcode)
			panic(errMsg)
		}
	} else {
		errMsg := fmt.Sprintf("Error: PREFIXCB %s %s", operand1, operand2)
		panic(errMsg)
	}
}

// RLC Rotate n left carry => bit0
func (cpu *CPU) RLC(operand1, operand2 string) {
	var value byte
	var bit7 byte
	if operand1 == "B" && operand2 == "*" {
		value = cpu.getBReg()
		bit7 = value >> 7
		value = (value << 1)
		if bit7 != 0 {
			value |= 1
		} else {
			value &= 0xfe
		}
		cpu.setBReg(value)
	} else if operand1 == "C" && operand2 == "*" {
		value = cpu.getCReg()
		bit7 = value >> 7
		value = (value << 1)
		if bit7 != 0 {
			value |= 1
		} else {
			value &= 0xfe
		}
		cpu.setCReg(value)
	} else if operand1 == "D" && operand2 == "*" {
		value = cpu.getDReg()
		bit7 = value >> 7
		value = (value << 1)
		if bit7 != 0 {
			value |= 1
		} else {
			value &= 0xfe
		}
		cpu.setDReg(value)
	} else if operand1 == "E" && operand2 == "*" {
		value = cpu.getEReg()
		bit7 = value >> 7
		value = (value << 1)
		if bit7 != 0 {
			value |= 1
		} else {
			value &= 0xfe
		}
		cpu.setEReg(value)
	} else if operand1 == "H" && operand2 == "*" {
		value = cpu.getHReg()
		bit7 = value >> 7
		value = (value << 1)
		if bit7 != 0 {
			value |= 1
		} else {
			value &= 0xfe
		}
		cpu.setHReg(value)
	} else if operand1 == "L" && operand2 == "*" {
		value = cpu.getLReg()
		bit7 = value >> 7
		value = (value << 1)
		if bit7 != 0 {
			value |= 1
		} else {
			value &= 0xfe
		}
		cpu.setLReg(value)
	} else if operand1 == "(HL)" && operand2 == "*" {
		value = cpu.FetchMemory8(cpu.Reg.HL)
		bit7 = value >> 7
		value = (value << 1)
		if bit7 != 0 {
			value |= 1
		} else {
			value &= 0xfe
		}
		cpu.SetMemory8(cpu.Reg.HL, value)
	} else if operand1 == "A" && operand2 == "*" {
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
func (cpu *CPU) RLCA(operand1, operand2 string) {
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
func (cpu *CPU) RRC(operand1, operand2 string) {
	var value byte
	var bit0 byte
	if operand1 == "B" && operand2 == "*" {
		value = cpu.getBReg()
		bit0 = value % 2
		value = (value >> 1)
		if bit0 != 0 {
			value |= 0x80
		} else {
			value &= 0x7f
		}
		cpu.setBReg(value)
	} else if operand1 == "C" && operand2 == "*" {
		value = cpu.getCReg()
		bit0 = value % 2
		value = (value >> 1)
		if bit0 != 0 {
			value |= 0x80
		} else {
			value &= 0x7f
		}
		cpu.setCReg(value)
	} else if operand1 == "D" && operand2 == "*" {
		value = cpu.getDReg()
		bit0 = value % 2
		value = (value >> 1)
		if bit0 != 0 {
			value |= 0x80
		} else {
			value &= 0x7f
		}
		cpu.setDReg(value)
	} else if operand1 == "E" && operand2 == "*" {
		value = cpu.getEReg()
		bit0 = value % 2
		value = (value >> 1)
		if bit0 != 0 {
			value |= 0x80
		} else {
			value &= 0x7f
		}
		cpu.setEReg(value)
	} else if operand1 == "H" && operand2 == "*" {
		value = cpu.getHReg()
		bit0 = value % 2
		value = (value >> 1)
		if bit0 != 0 {
			value |= 0x80
		} else {
			value &= 0x7f
		}
		cpu.setHReg(value)
	} else if operand1 == "L" && operand2 == "*" {
		value = cpu.getLReg()
		bit0 = value % 2
		value = (value >> 1)
		if bit0 != 0 {
			value |= 0x80
		} else {
			value &= 0x7f
		}
		cpu.setLReg(value)
	} else if operand1 == "(HL)" && operand2 == "*" {
		value = cpu.FetchMemory8(cpu.Reg.HL)
		bit0 = value % 2
		value = (value >> 1)
		if bit0 != 0 {
			value |= 0x80
		} else {
			value &= 0x7f
		}
		cpu.SetMemory8(cpu.Reg.HL, value)
	} else if operand1 == "A" && operand2 == "*" {
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
func (cpu *CPU) RRCA(operand1, operand2 string) {
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
func (cpu *CPU) RL(operand1, operand2 string) {
	var value byte
	var bit7 byte
	carry := cpu.getCFlag()

	switch operand1 {
	case "A":
		value = cpu.getAReg()
		bit7 = value >> 7
		value = (value << 1)
		if carry {
			value |= 1
		} else {
			value &= 0xfe
		}
		cpu.setAReg(value)
	case "B":
		value = cpu.getBReg()
		bit7 = value >> 7
		value = (value << 1)
		if carry {
			value |= 1
		} else {
			value &= 0xfe
		}
		cpu.setBReg(value)
	case "C":
		value = cpu.getCReg()
		bit7 = value >> 7
		value = (value << 1)
		if carry {
			value |= 1
		} else {
			value &= 0xfe
		}
		cpu.setCReg(value)
	case "D":
		value = cpu.getDReg()
		bit7 = value >> 7
		value = (value << 1)
		if carry {
			value |= 1
		} else {
			value &= 0xfe
		}
		cpu.setDReg(value)
	case "E":
		value = cpu.getEReg()
		bit7 = value >> 7
		value = (value << 1)
		if carry {
			value |= 1
		} else {
			value &= 0xfe
		}
		cpu.setEReg(value)
	case "H":
		value = cpu.getHReg()
		bit7 = value >> 7
		value = (value << 1)
		if carry {
			value |= 1
		} else {
			value &= 0xfe
		}
		cpu.setHReg(value)
	case "L":
		value = cpu.getLReg()
		bit7 = value >> 7
		value = (value << 1)
		if carry {
			value |= 1
		} else {
			value &= 0xfe
		}
		cpu.setLReg(value)
	case "(HL)":
		value = cpu.FetchMemory8(cpu.Reg.HL)
		bit7 = value >> 7
		value = (value << 1)
		if carry {
			value |= 1
		} else {
			value &= 0xfe
		}
		cpu.SetMemory8(cpu.Reg.HL, value)
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
func (cpu *CPU) RLA(operand1, operand2 string) {
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
func (cpu *CPU) RR(operand1, operand2 string) {
	var value byte
	var bit0 byte
	carry := cpu.getCFlag()

	switch operand1 {
	case "A":
		value = cpu.getAReg()
		bit0 = value % 2
		value = (value >> 1)
		if carry {
			value |= 0x80
		} else {
			value &= 0x7f
		}
		cpu.setAReg(value)
	case "B":
		value = cpu.getBReg()
		bit0 = value % 2
		value = (value >> 1)
		if carry {
			value |= 0x80
		} else {
			value &= 0x7f
		}
		cpu.setBReg(value)
	case "C":
		value = cpu.getCReg()
		bit0 = value % 2
		value = (value >> 1)
		if carry {
			value |= 0x80
		} else {
			value &= 0x7f
		}
		cpu.setCReg(value)
	case "D":
		value = cpu.getDReg()
		bit0 = value % 2
		value = (value >> 1)
		if carry {
			value |= 0x80
		} else {
			value &= 0x7f
		}
		cpu.setDReg(value)
	case "E":
		value = cpu.getEReg()
		bit0 = value % 2
		value = (value >> 1)
		if carry {
			value |= 0x80
		} else {
			value &= 0x7f
		}
		cpu.setEReg(value)
	case "H":
		value = cpu.getHReg()
		bit0 = value % 2
		value = (value >> 1)
		if carry {
			value |= 0x80
		} else {
			value &= 0x7f
		}
		cpu.setHReg(value)
	case "L":
		value = cpu.getLReg()
		bit0 = value % 2
		value = (value >> 1)
		if carry {
			value |= 0x80
		} else {
			value &= 0x7f
		}
		cpu.setLReg(value)
	case "(HL)":
		value = cpu.FetchMemory8(cpu.Reg.HL)
		bit0 = value % 2
		value = (value >> 1)
		if carry {
			value |= 0x80
		} else {
			value &= 0x7f
		}
		cpu.SetMemory8(cpu.Reg.HL, value)
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
func (cpu *CPU) SLA(operand1, operand2 string) {
	var value byte
	var bit7 byte
	if operand1 == "B" && operand2 == "*" {
		value = cpu.getBReg()
		bit7 = value >> 7
		value = (value << 1)
		cpu.setBReg(value)
	} else if operand1 == "C" && operand2 == "*" {
		value = cpu.getCReg()
		bit7 = value >> 7
		value = (value << 1)
		cpu.setCReg(value)
	} else if operand1 == "D" && operand2 == "*" {
		value = cpu.getDReg()
		bit7 = value >> 7
		value = (value << 1)
		cpu.setDReg(value)
	} else if operand1 == "E" && operand2 == "*" {
		value = cpu.getEReg()
		bit7 = value >> 7
		value = (value << 1)
		cpu.setEReg(value)
	} else if operand1 == "H" && operand2 == "*" {
		value = cpu.getHReg()
		bit7 = value >> 7
		value = (value << 1)
		cpu.setHReg(value)
	} else if operand1 == "L" && operand2 == "*" {
		value = cpu.getLReg()
		bit7 = value >> 7
		value = (value << 1)
		cpu.setLReg(value)
	} else if operand1 == "(HL)" && operand2 == "*" {
		value = cpu.FetchMemory8(cpu.Reg.HL)
		bit7 = value >> 7
		value = (value << 1)
		cpu.SetMemory8(cpu.Reg.HL, value)
	} else if operand1 == "A" && operand2 == "*" {
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
func (cpu *CPU) SRA(operand1, operand2 string) {
	var value byte
	var bit0 byte
	if operand1 == "B" && operand2 == "*" {
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
	} else if operand1 == "C" && operand2 == "*" {
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
	} else if operand1 == "D" && operand2 == "*" {
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
	} else if operand1 == "E" && operand2 == "*" {
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
	} else if operand1 == "H" && operand2 == "*" {
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
	} else if operand1 == "L" && operand2 == "*" {
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
	} else if operand1 == "(HL)" && operand2 == "*" {
		value = cpu.FetchMemory8(cpu.Reg.HL)
		bit0 = value % 2
		bit7 := value >> 7
		value = (value >> 1)
		if bit7 != 0 {
			value |= 0x80
		} else {
			value &= 0x7f
		}
		cpu.SetMemory8(cpu.Reg.HL, value)
	} else if operand1 == "A" && operand2 == "*" {
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

// SWAP Swap n[4:8] and n[0:4]
func (cpu *CPU) SWAP(operand1, operand2 string) {
	var value byte

	switch operand1 {
	case "B":
		B := cpu.getBReg()
		B03 := B & 0x0f
		B47 := B >> 4
		value = (B03 << 4) | B47
		cpu.setBReg(value)
	case "C":
		C := cpu.getCReg()
		C03 := C & 0x0f
		C47 := C >> 4
		value = (C03 << 4) | C47
		cpu.setCReg(value)
	case "D":
		D := cpu.getDReg()
		D03 := D & 0x0f
		D47 := D >> 4
		value = (D03 << 4) | D47
		cpu.setDReg(value)
	case "E":
		E := cpu.getEReg()
		E03 := E & 0x0f
		E47 := E >> 4
		value = (E03 << 4) | E47
		cpu.setEReg(value)
	case "H":
		H := cpu.getHReg()
		H03 := H & 0x0f
		H47 := H >> 4
		value = (H03 << 4) | H47
		cpu.setHReg(value)
	case "L":
		L := cpu.getLReg()
		L03 := L & 0x0f
		L47 := L >> 4
		value = (L03 << 4) | L47
		cpu.setLReg(value)
	case "(HL)":
		data := cpu.FetchMemory8(cpu.Reg.HL)
		data03 := data & 0x0f
		data47 := data >> 4
		value = (data03 << 4) | data47
		cpu.SetMemory8(cpu.Reg.HL, value)
	case "A":
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
func (cpu *CPU) SRL(operand1, operand2 string) {
	var value byte
	var bit0 byte

	switch operand1 {
	case "B":
		value = cpu.getBReg()
		bit0 = value % 2
		value = (value >> 1)
		cpu.setBReg(value)
	case "C":
		value = cpu.getCReg()
		bit0 = value % 2
		value = (value >> 1)
		cpu.setCReg(value)
	case "D":
		value = cpu.getDReg()
		bit0 = value % 2
		value = (value >> 1)
		cpu.setDReg(value)
	case "E":
		value = cpu.getEReg()
		bit0 = value % 2
		value = (value >> 1)
		cpu.setEReg(value)
	case "H":
		value = cpu.getHReg()
		bit0 = value % 2
		value = (value >> 1)
		cpu.setHReg(value)
	case "L":
		value = cpu.getLReg()
		bit0 = value % 2
		value = (value >> 1)
		cpu.setLReg(value)
	case "(HL)":
		value = cpu.FetchMemory8(cpu.Reg.HL)
		bit0 = value % 2
		value = (value >> 1)
		cpu.SetMemory8(cpu.Reg.HL, value)
	case "A":
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
func (cpu *CPU) BIT(operand1, operand2 string) {
	var value byte
	tmp, _ := strconv.Atoi(operand1)
	targetBit := uint(tmp) // ターゲットのbit
	switch operand2 {
	case "B":
		value = (cpu.getBReg() >> targetBit) % 2
	case "C":
		value = (cpu.getCReg() >> targetBit) % 2
	case "D":
		value = (cpu.getDReg() >> targetBit) % 2
	case "E":
		value = (cpu.getEReg() >> targetBit) % 2
	case "H":
		value = (cpu.getHReg() >> targetBit) % 2
	case "L":
		value = (cpu.getLReg() >> targetBit) % 2
	case "(HL)":
		value = (cpu.FetchMemory8(cpu.Reg.HL) >> targetBit) % 2
	case "A":
		value = (cpu.getAReg() >> targetBit) % 2
	}

	cpu.flagZ(value)
	cpu.flagN(false)
	cpu.setHFlag()
	cpu.Reg.PC++
}

// RES Clear bit n
func (cpu *CPU) RES(operand1, operand2 string) {
	tmp, _ := strconv.Atoi(operand1)
	targetBit := uint(tmp) // ターゲットのbit
	switch operand2 {
	case "B":
		mask := ^(byte(1) << targetBit)
		B := cpu.getBReg() & mask
		cpu.setBReg(B)
	case "C":
		mask := ^(byte(1) << targetBit)
		C := cpu.getCReg() & mask
		cpu.setCReg(C)
	case "D":
		mask := ^(byte(1) << targetBit)
		D := cpu.getDReg() & mask
		cpu.setDReg(D)
	case "E":
		mask := ^(byte(1) << targetBit)
		E := cpu.getEReg() & mask
		cpu.setEReg(E)
	case "H":
		mask := ^(byte(1) << targetBit)
		H := cpu.getHReg() & mask
		cpu.setHReg(H)
	case "L":
		mask := ^(byte(1) << targetBit)
		L := cpu.getLReg() & mask
		cpu.setLReg(L)
	case "(HL)":
		mask := ^(byte(1) << targetBit)
		value := cpu.FetchMemory8(cpu.Reg.HL) & mask
		cpu.SetMemory8(cpu.Reg.HL, value)
	case "A":
		mask := ^(byte(1) << targetBit)
		A := cpu.getAReg() & mask
		cpu.setAReg(A)
	}
	cpu.Reg.PC++
}

// SET Clear bit n
func (cpu *CPU) SET(operand1, operand2 string) {
	tmp, _ := strconv.Atoi(operand1)
	targetBit := uint(tmp) // ターゲットのbit
	switch operand2 {
	case "B":
		mask := byte(1) << targetBit
		B := cpu.getBReg() | mask
		cpu.setBReg(B)
	case "C":
		mask := byte(1) << targetBit
		C := cpu.getCReg() | mask
		cpu.setCReg(C)
	case "D":
		mask := byte(1) << targetBit
		D := cpu.getDReg() | mask
		cpu.setDReg(D)
	case "E":
		mask := byte(1) << targetBit
		E := cpu.getEReg() | mask
		cpu.setEReg(E)
	case "H":
		mask := byte(1) << targetBit
		H := cpu.getHReg() | mask
		cpu.setHReg(H)
	case "L":
		mask := byte(1) << targetBit
		L := cpu.getLReg() | mask
		cpu.setLReg(L)
	case "(HL)":
		mask := byte(1) << targetBit
		value := cpu.FetchMemory8(cpu.Reg.HL) | mask
		cpu.SetMemory8(cpu.Reg.HL, value)
	case "A":
		mask := byte(1) << targetBit
		A := cpu.getAReg() | mask
		cpu.setAReg(A)
	}
	cpu.Reg.PC++
}

// PUSH スタックにPUSH
func (cpu *CPU) PUSH(operand1, operand2 string) {
	switch operand1 {
	case "BC":
		cpu.pushBC()
	case "DE":
		cpu.pushDE()
	case "HL":
		cpu.pushHL()
	case "AF":
		cpu.pushAF()
	default:
		errMsg := fmt.Sprintf("Error: PUSH %s %s", operand1, operand2)
		panic(errMsg)
	}
	cpu.Reg.PC++
}

// POP スタックからPOP
func (cpu *CPU) POP(operand1, operand2 string) {
	switch operand1 {
	case "BC":
		cpu.popBC()
	case "DE":
		cpu.popDE()
	case "HL":
		cpu.popHL()
	case "AF":
		cpu.popAF()
	default:
		errMsg := fmt.Sprintf("Error: POP %s %s", operand1, operand2)
		panic(errMsg)
	}
	cpu.Reg.PC++
}

// SUB 減算
func (cpu *CPU) SUB(operand1, operand2 string) {
	switch operand1 {
	case "A":
		value := cpu.getAReg() - cpu.getAReg()
		carryBits := cpu.getAReg() ^ cpu.getAReg() ^ value
		cpu.flagC8Sub(cpu.getAReg(), cpu.getAReg())
		cpu.setAReg(value)
		cpu.flagZ(value)
		cpu.flagN(true)
		cpu.flagH8(carryBits)
		cpu.Reg.PC++
	case "B":
		value := cpu.getAReg() - cpu.getBReg()
		carryBits := cpu.getAReg() ^ cpu.getBReg() ^ value
		cpu.flagC8Sub(cpu.getAReg(), cpu.getBReg())
		cpu.setAReg(value)
		cpu.flagZ(value)
		cpu.flagN(true)
		cpu.flagH8(carryBits)
		cpu.Reg.PC++
	case "C":
		value := cpu.getAReg() - cpu.getCReg()
		carryBits := cpu.getAReg() ^ cpu.getCReg() ^ value
		cpu.flagC8Sub(cpu.getAReg(), cpu.getCReg())
		cpu.setAReg(value)
		cpu.flagZ(value)
		cpu.flagN(true)
		cpu.flagH8(carryBits)
		cpu.Reg.PC++
	case "D":
		value := cpu.getAReg() - cpu.getDReg()
		carryBits := cpu.getAReg() ^ cpu.getDReg() ^ value
		cpu.flagC8Sub(cpu.getAReg(), cpu.getDReg())
		cpu.setAReg(value)
		cpu.flagZ(value)
		cpu.flagN(true)
		cpu.flagH8(carryBits)
		cpu.Reg.PC++
	case "E":
		value := cpu.getAReg() - cpu.getEReg()
		carryBits := cpu.getAReg() ^ cpu.getEReg() ^ value
		cpu.flagC8Sub(cpu.getAReg(), cpu.getEReg())
		cpu.setAReg(value)
		cpu.flagZ(value)
		cpu.flagN(true)
		cpu.flagH8(carryBits)
		cpu.Reg.PC++
	case "H":
		value := cpu.getAReg() - cpu.getHReg()
		carryBits := cpu.getAReg() ^ cpu.getHReg() ^ value
		cpu.flagC8Sub(cpu.getAReg(), cpu.getHReg())
		cpu.setAReg(value)
		cpu.flagZ(value)
		cpu.flagN(true)
		cpu.flagH8(carryBits)
		cpu.Reg.PC++
	case "L":
		value := cpu.getAReg() - cpu.getLReg()
		carryBits := cpu.getAReg() ^ cpu.getLReg() ^ value
		cpu.flagC8Sub(cpu.getAReg(), cpu.getLReg())
		cpu.setAReg(value)
		cpu.flagZ(value)
		cpu.flagN(true)
		cpu.flagH8(carryBits)
		cpu.Reg.PC++
	case "d8":
		value := cpu.getAReg() - cpu.d8Fetch()
		carryBits := cpu.getAReg() ^ cpu.d8Fetch() ^ value
		cpu.flagC8Sub(cpu.getAReg(), cpu.d8Fetch())
		cpu.setAReg(value)
		cpu.flagZ(value)
		cpu.flagN(true)
		cpu.flagH8(carryBits)
		cpu.Reg.PC += 2
	case "(HL)":
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
func (cpu *CPU) RRA(operand1, operand2 string) {
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
func (cpu *CPU) ADC(operand1, operand2 string) {
	var carry, value, value4 byte
	var value16 uint16
	if cpu.getCFlag() {
		carry = 1
	} else {
		carry = 0
	}

	switch operand1 {
	case "A":
		switch operand2 {
		case "A":
			value = cpu.getAReg() + carry + cpu.getAReg()
			value4 = cpu.getARegLower4() + carry + cpu.getARegLower4()
			value16 = uint16(cpu.getAReg()) + uint16(cpu.getAReg()) + uint16(carry)
		case "B":
			value = cpu.getBReg() + carry + cpu.getAReg()
			value4 = cpu.getBRegLower4() + carry + cpu.getARegLower4()
			value16 = uint16(cpu.getBReg()) + uint16(carry) + uint16(cpu.getAReg())
		case "C":
			value = cpu.getCReg() + carry + cpu.getAReg()
			value4 = cpu.getCRegLower4() + carry + cpu.getARegLower4()
			value16 = uint16(cpu.getCReg()) + uint16(carry) + uint16(cpu.getAReg())
		case "D":
			value = cpu.getDReg() + carry + cpu.getAReg()
			value4 = cpu.getDRegLower4() + carry + cpu.getARegLower4()
			value16 = uint16(cpu.getDReg()) + uint16(carry) + uint16(cpu.getAReg())
		case "E":
			value = cpu.getEReg() + carry + cpu.getAReg()
			value4 = cpu.getERegLower4() + carry + cpu.getARegLower4()
			value16 = uint16(cpu.getEReg()) + uint16(carry) + uint16(cpu.getAReg())
		case "H":
			value = cpu.getHReg() + carry + cpu.getAReg()
			value4 = cpu.getHRegLower4() + carry + cpu.getARegLower4()
			value16 = uint16(cpu.getHReg()) + uint16(carry) + uint16(cpu.getAReg())
		case "L":
			value = cpu.getLReg() + carry + cpu.getAReg()
			value4 = cpu.getLRegLower4() + carry + cpu.getARegLower4()
			value16 = uint16(cpu.getLReg()) + uint16(carry) + uint16(cpu.getAReg())
		case "(HL)":
			data := cpu.FetchMemory8(cpu.Reg.HL)
			value = data + carry + cpu.getAReg()
			value4 = (data & 0x0f) + carry + cpu.getARegLower4()
			value16 = uint16(data) + uint16(cpu.getAReg()) + uint16(carry)
		case "d8":
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
func (cpu *CPU) SBC(operand1, operand2 string) {
	var carry, value, value4 byte
	var value16 uint16
	if cpu.getCFlag() {
		carry = 1
	} else {
		carry = 0
	}

	switch operand1 {
	case "A":
		switch operand2 {
		case "A":
			value = cpu.getAReg() - (cpu.getAReg() + carry)
			value4 = cpu.getARegLower4() - (cpu.getARegLower4() + carry)
			value16 = uint16(cpu.getAReg()) - (uint16(cpu.getAReg()) + uint16(carry))
		case "B":
			value = cpu.getAReg() - (cpu.getBReg() + carry)
			value4 = cpu.getARegLower4() - (cpu.getBRegLower4() + carry)
			value16 = uint16(cpu.getAReg()) - (uint16(cpu.getBReg()) + uint16(carry))
		case "C":
			value = cpu.getAReg() - (cpu.getCReg() + carry)
			value4 = cpu.getARegLower4() - (cpu.getCRegLower4() + carry)
			value16 = uint16(cpu.getAReg()) - (uint16(cpu.getCReg()) + uint16(carry))
		case "D":
			value = cpu.getAReg() - (cpu.getDReg() + carry)
			value4 = cpu.getARegLower4() - (cpu.getDRegLower4() + carry)
			value16 = uint16(cpu.getAReg()) - (uint16(cpu.getDReg()) + uint16(carry))
		case "E":
			value = cpu.getAReg() - (cpu.getEReg() + carry)
			value4 = cpu.getARegLower4() - (cpu.getERegLower4() + carry)
			value16 = uint16(cpu.getAReg()) - (uint16(cpu.getEReg()) + uint16(carry))
		case "H":
			value = cpu.getAReg() - (cpu.getHReg() + carry)
			value4 = cpu.getARegLower4() - (cpu.getHRegLower4() + carry)
			value16 = uint16(cpu.getAReg()) - (uint16(cpu.getHReg()) + uint16(carry))
		case "L":
			value = cpu.getAReg() - (cpu.getLReg() + carry)
			value4 = cpu.getARegLower4() - (cpu.getLRegLower4() + carry)
			value16 = uint16(cpu.getAReg()) - (uint16(cpu.getLReg()) + uint16(carry))
		case "(HL)":
			data := cpu.FetchMemory8(cpu.Reg.HL)
			value = cpu.getAReg() - (data + carry)
			value4 = cpu.getARegLower4() - ((data & 0x0f) + carry)
			value16 = uint16(cpu.getAReg()) - (uint16(data) + uint16(carry))
		case "d8":
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
func (cpu *CPU) DAA(operand1, operand2 string) {
	A := uint8(cpu.getAReg())
	// 参考: https://forums.nesdev.com/viewtopic.php?f=20&t=15944
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
func (cpu *CPU) RST(operand1, operand2 string) {
	vector, err := strconv.Atoi(operand1)
	if err != nil {
		panic(err)
	}
	destination := uint16(vector)
	cpu.Reg.PC++
	cpu.pushPC()
	cpu.Reg.PC = destination
}

// SCF Set Carry Flag
func (cpu *CPU) SCF(operand1, operand2 string) {
	cpu.flagN(false)
	cpu.clearHFlag()
	cpu.setCFlag()
	cpu.Reg.PC++
}

// CCF Complement Carry Flag
func (cpu *CPU) CCF(operand1, operand2 string) {
	cpu.flagN(false)
	cpu.clearHFlag()
	if cpu.getCFlag() {
		cpu.clearCFlag()
	} else {
		cpu.setCFlag()
	}
	cpu.Reg.PC++
}
