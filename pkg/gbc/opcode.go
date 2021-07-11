package gbc

import (
	"fmt"
	"gbc/pkg/gbc/scheduler"
	"gbc/pkg/util"
)

func (g *GBC) a16Fetch() uint16 {
	return g.d16Fetch()
}

func (g *GBC) a16FetchJP() uint16 {
	lower := uint16(g.Load8(g.Reg.PC + 1)) // M = 1: nn read: memory access for low byte
	g.timer.tick(1 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))
	upper := uint16(g.Load8(g.Reg.PC + 2)) // M = 2: nn read: memory access for high byte
	g.timer.tick(1 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))
	value := (upper << 8) | lower
	return value
}

func (g *GBC) d8Fetch() byte {
	return g.Load8(g.Reg.PC + 1)
}

func (g *GBC) d16Fetch() uint16 {
	lower, upper := uint16(g.Load8(g.Reg.PC+1)), uint16(g.Load8(g.Reg.PC+2))
	return (upper << 8) | lower
}

// LD R8,R8
func ld8r(g *GBC, op1, op2 int) {
	g.Reg.R[op1] = g.Reg.R[op2]
	g.Reg.PC++
}

// ------ LD A, *

// ld r8, mem[r16]
func ld8m(g *GBC, r8, r16 int) {
	g.Reg.R[r8] = g.Load8(g.Reg.R16(r16))
	g.Reg.PC++
}

// ld r8, mem[imm]
func ld8i(g *GBC, r8, _ int) {
	g.Reg.R[r8] = g.d8Fetch()
	g.Reg.PC += 2
}

// LD A, (u16)
func op0xfa(g *GBC, operand1, operand2 int) {
	g.Reg.R[A] = g.Load8(g.a16FetchJP())
	g.Reg.PC += 3
	g.timer.tick(2 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))
}

// LD A,(FF00+C)
func op0xf2(g *GBC, operand1, operand2 int) {
	addr := 0xff00 + uint16(g.Reg.R[C])
	g.Reg.R[A] = g.loadIO(byte(addr))
	g.Reg.PC++ // mistake?(https://www.pastraiser.com/g/gameboy/gameboy_opcodes.html)
}

// ------ LD (HL), *

// LD (HL),u8
func op0x36(g *GBC, operand1, operand2 int) {
	value := g.d8Fetch()
	g.timer.tick(1 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))
	g.Store8(g.Reg.HL(), value)
	g.Reg.PC += 2
	g.timer.tick(2 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))
}

// LD (HL),R8
func ldHLR8(g *GBC, unused, op int) {
	g.Store8(g.Reg.HL(), g.Reg.R[op])
	g.Reg.PC++
}

// ------ others ld

// LD (u16),SP
func op0x08(g *GBC, operand1, operand2 int) {
	// Store SP into addresses n16 (LSB) and n16 + 1 (MSB).
	addr := g.a16Fetch()
	upper, lower := byte(g.Reg.SP>>8), byte(g.Reg.SP) // MSB
	g.Store8(addr, lower)
	g.Store8(addr+1, upper)
	g.Reg.PC += 3
	g.timer.tick(5 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))
}

// LD (u16),A
func op0xea(g *GBC, operand1, operand2 int) {
	g.Store8(g.a16FetchJP(), g.Reg.R[A])
	g.Reg.PC += 3
	g.timer.tick(2 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))
}

// ld r16, u16
func ld16i(g *GBC, r16, _ int) {
	g.Reg.setR16(r16, g.d16Fetch())
	g.Reg.PC += 3
}

// LD HL,SP+i8
func op0xf8(g *GBC, operand1, operand2 int) {
	delta := int8(g.Load8(g.Reg.PC + 1))
	value := int32(g.Reg.SP) + int32(delta)
	carryBits := uint32(g.Reg.SP) ^ uint32(delta) ^ uint32(value)
	g.Reg.setHL(uint16(value))
	g.setF(flagZ, false)
	g.setF(flagN, false)
	g.setF(flagC, util.Bit(carryBits, 8))
	g.setF(flagH, util.Bit(carryBits, 4))
	g.Reg.PC += 2
}

// LD SP,HL
func op0xf9(g *GBC, operand1, operand2 int) {
	g.Reg.SP = g.Reg.HL()
	g.Reg.PC++
}

// LD (FF00+C),A
func op0xe2(g *GBC, operand1, operand2 int) {
	addr := 0xff00 + uint16(g.Reg.R[C])
	g.Store8(addr, g.Reg.R[A])
	g.Reg.PC++ // mistake?(https://www.pastraiser.com/g/gameboy/gameboy_opcodes.html)
}

func ldm16r(g *GBC, r16, r8 int) {
	g.Store8(g.Reg.R16(r16), g.Reg.R[r8])
	g.Reg.PC++
}

// LD ($FF00+a8),A
func op0xe0(g *GBC, operand1, operand2 int) {
	addr := 0xff00 + uint16(g.Load8(g.Reg.PC+1))
	g.timer.tick(1 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))
	g.storeIO(byte(addr), g.Reg.R[A])
	g.Reg.PC += 2
	g.timer.tick(2 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))
}

// LD A,($FF00+a8)
func op0xf0(g *GBC, operand1, operand2 int) {
	addr := 0xff00 + uint16(g.Load8(g.Reg.PC+1))
	g.timer.tick(1 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))
	value := g.loadIO(byte(addr))

	g.Reg.R[A] = value
	g.Reg.PC += 2
	g.timer.tick(2 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))
}

// No operation
func nop(g *GBC, operand1, operand2 int) { g.Reg.PC++ }

// INC Increment

func inc8(g *GBC, r8, _ int) {
	value := g.Reg.R[r8] + 1
	carryBits := g.Reg.R[r8] ^ 1 ^ value
	g.Reg.R[r8] = value

	g.setF(flagZ, value == 0)
	g.setF(flagN, false)
	g.setF(flagH, util.Bit(carryBits, 4))
	g.Reg.PC++
}

func inc16(g *GBC, r16, _ int) {
	g.Reg.setR16(r16, g.Reg.R16(r16)+1)
	g.Reg.PC++
}

func incHL(g *GBC, _, _ int) {
	value := g.Load8(g.Reg.HL()) + 1
	g.timer.tick(1 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))
	carryBits := g.Load8(g.Reg.HL()) ^ 1 ^ value
	g.Store8(g.Reg.HL(), value)
	g.timer.tick(2 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))

	g.setF(flagZ, value == 0)
	g.setF(flagN, false)
	g.setF(flagH, util.Bit(carryBits, 4))
	g.Reg.PC++
}

// DEC Decrement

func dec8(g *GBC, r8, _ int) {
	value := g.Reg.R[r8] - 1
	carryBits := g.Reg.R[r8] ^ 1 ^ value
	g.Reg.R[r8] = value

	g.setF(flagZ, value == 0)
	g.setF(flagN, true)
	g.setF(flagH, util.Bit(carryBits, 4))
	g.Reg.PC++
}

func dec16(g *GBC, r16, _ int) {
	g.Reg.setR16(r16, g.Reg.R16(r16)-1)
	g.Reg.PC++
}

func decHL(g *GBC, _, _ int) {
	value := g.Load8(g.Reg.HL()) - 1
	g.timer.tick(1 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))
	carryBits := g.Load8(g.Reg.HL()) ^ 1 ^ value
	g.Store8(g.Reg.HL(), value)
	g.timer.tick(2 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))

	g.setF(flagZ, value == 0)
	g.setF(flagN, true)
	g.setF(flagH, util.Bit(carryBits, 4))
	g.Reg.PC++
}

// --------- JR ---------

// jr i8
func jr(g *GBC, _, _ int) {
	delta := int8(g.Load8(g.Reg.PC + 1))
	g.Reg.PC = uint16(int32(g.Reg.PC+2) + int32(delta)) // PC+2 because of time after fetch(pc is incremented)
	g.timer.tick(3 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))
}

// jr cc,i8
func jrcc(g *GBC, cc, _ int) {
	if g.f(cc) {
		jr(g, 0, 0)
	} else {
		g.Reg.PC += 2
		g.timer.tick(2 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))
	}
}

// jr ncc,i8 (ncc = not cc)
func jrncc(g *GBC, cc, _ int) {
	if !g.f(cc) {
		jr(g, 0, 0)
	} else {
		g.Reg.PC += 2
		g.timer.tick(2 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))
	}
}

func halt(g *GBC, _, _ int) {
	g.Reg.PC++
	if g.IO[IEIO]&g.IO[IFIO]&0x1f == 0 {
		g.halt = true
	}
}

// stop GBC
func stop(g *GBC, _, _ int) {
	g.Reg.PC += 2
	if g.model >= util.GB_MODEL_CGB && util.Bit(g.IO[KEY1IO], 0) {
		g.DoubleSpeed = !g.DoubleSpeed
		g.IO[KEY1IO] = byte(util.Bool2Int(g.DoubleSpeed)) << 7
	} else {
		sleep := ^(g.loadIO(JOYPIO) & 0x30)
		if sleep > 0 {
			fmt.Println("TODO: impl sleep on stop")
		}
	}
}

// XOR xor
func xor8(g *GBC, _, r8 int) {
	value := g.Reg.R[A] ^ g.Reg.R[r8]
	g.Reg.R[A] = value
	g.setF(flagZ, value == 0)
	g.setF(flagN, false)
	g.setF(flagH, false)
	g.setF(flagC, false)
	g.Reg.PC++
}

// XOR A,(HL)
func xoraHL(g *GBC, _, _ int) {
	value := g.Reg.R[A] ^ g.Load8(g.Reg.HL())
	g.Reg.R[A] = value
	g.setF(flagZ, value == 0)
	g.setF(flagN, false)
	g.setF(flagH, false)
	g.setF(flagC, false)
	g.Reg.PC++
}

// XOR A,u8
func xoru8(g *GBC, _, _ int) {
	value := g.Reg.R[A] ^ g.Load8(g.Reg.PC+1)
	g.Reg.PC++
	g.Reg.R[A] = value
	g.setF(flagZ, value == 0)
	g.setF(flagN, false)
	g.setF(flagH, false)
	g.setF(flagC, false)
	g.Reg.PC++
}

// jp u16
func jp(g *GBC, _, _ int) {
	g.Reg.PC = g.a16FetchJP()
	g.timer.tick(2 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))
}

func jpcc(g *GBC, cc, _ int) {
	dst := g.a16FetchJP()
	if g.f(cc) {
		g.Reg.PC = dst
		g.timer.tick(2 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))
	} else {
		g.Reg.PC += 3
		g.timer.tick(1 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))
	}
}

func jpncc(g *GBC, cc, _ int) {
	dst := g.a16FetchJP()
	if !g.f(cc) {
		g.Reg.PC = dst
		g.timer.tick(2 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))
	} else {
		g.Reg.PC += 3
		g.timer.tick(1 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))
	}
}

func jpHL(g *GBC, _, _ int) {
	g.Reg.PC = g.Reg.HL()
	g.timer.tick(1 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))
}

// Return
func ret(g *GBC, _, _ int) {
	g.popPC()
}

func retcc(g *GBC, cc, _ int) {
	if g.f(cc) {
		g.popPC()
		g.timer.tick(5 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))
	} else {
		g.Reg.PC++
		g.timer.tick(2 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))
	}
}

// not retcc
func retncc(g *GBC, cc, _ int) {
	if !g.f(cc) {
		g.popPC()
		g.timer.tick(5 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))
	} else {
		g.Reg.PC++
		g.timer.tick(2 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))
	}
}

// Return Interrupt
func reti(g *GBC, operand1, operand2 int) {
	g.popPC()
	g.scheduler.ScheduleEvent(scheduler.EiPending, func(cyclesLate uint64) {
		g.Reg.IME = true
		g.updateIRQs()
	}, 4>>util.Bool2U64(g.DoubleSpeed))
}

func call(g *GBC, _, _ int) {
	dst := g.a16FetchJP()
	g.Reg.PC += 3
	g.timer.tick(1 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))
	g.pushPCCALL()
	g.timer.tick(1 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))
	g.Reg.PC = dst
}

func callcc(g *GBC, cc, _ int) {
	if g.f(cc) {
		call(g, 0, 0)
		return
	}
	g.Reg.PC += 3
	g.timer.tick(3 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))
}

func callncc(g *GBC, cc, _ int) {
	if !g.f(cc) {
		call(g, 0, 0)
		return
	}
	g.Reg.PC += 3
	g.timer.tick(3 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))
}

// DI Disable Interrupt
func di(g *GBC, _, _ int) {
	g.Reg.IME = false
	g.updateIRQs()
	g.Reg.PC++
}

// EI Enable Interrupt
func ei(g *GBC, _, _ int) {
	// TODO ref: https://github.com/Gekkio/mooneye-gb/blob/master/tests/acceptance/halt_ime0_ei.s#L23
	g.Reg.PC++
	g.scheduler.DescheduleEvent(scheduler.EiPending)
	g.scheduler.ScheduleEvent(scheduler.EiPending, func(cyclesLate uint64) {
		g.Reg.IME = true
		g.updateIRQs()
	}, 4>>util.Bool2U64(g.DoubleSpeed))
}

// CP Compare
func cp(g *GBC, _, r8 int) {
	value := g.Reg.R[A] - g.Reg.R[r8]
	carryBits := g.Reg.R[A] ^ g.Reg.R[r8] ^ value
	g.setCSub(g.Reg.R[A], g.Reg.R[r8])

	g.setF(flagZ, value == 0)
	g.setF(flagN, true)
	g.setF(flagH, util.Bit(carryBits, 4))
	g.Reg.PC++
}

// CP A,(HL)
func cpaHL(g *GBC, _, _ int) {
	value := g.Reg.R[A] - g.Load8(g.Reg.HL())
	carryBits := g.Reg.R[A] ^ g.Load8(g.Reg.HL()) ^ value
	g.setCSub(g.Reg.R[A], g.Load8(g.Reg.HL()))
	g.setF(flagZ, value == 0)
	g.setF(flagN, true)
	g.setF(flagH, util.Bit(carryBits, 4))
	g.Reg.PC++
}

// CP A,u8
func cpu8(g *GBC, _, _ int) {
	value := g.Reg.R[A] - g.d8Fetch()
	carryBits := g.Reg.R[A] ^ g.d8Fetch() ^ value
	g.setCSub(g.Reg.R[A], g.d8Fetch())
	g.Reg.PC++
	g.setF(flagZ, value == 0)
	g.setF(flagN, true)
	g.setF(flagH, util.Bit(carryBits, 4))
	g.Reg.PC++
}

// AND And instruction

func and8(g *GBC, _, r8 int) {
	value := g.Reg.R[A] & g.Reg.R[r8]
	g.Reg.R[A] = value
	g.setF(flagZ, value == 0)
	g.setF(flagN, false)
	g.setF(flagH, true)
	g.setF(flagC, false)
	g.Reg.PC++
}

// AND A,(HL)
func op0xa6(g *GBC, _, _ int) {
	value := g.Reg.R[A] & g.Load8(g.Reg.HL())
	g.Reg.R[A] = value
	g.setF(flagZ, value == 0)
	g.setF(flagN, false)
	g.setF(flagH, true)
	g.setF(flagC, false)
	g.Reg.PC++
}

// AND A,u8
func andu8(g *GBC, _, _ int) {
	value := g.Reg.R[A] & g.d8Fetch()
	g.Reg.PC++

	g.Reg.R[A] = value
	g.setF(flagZ, value == 0)
	g.setF(flagN, false)
	g.setF(flagH, true)
	g.setF(flagC, false)
	g.Reg.PC++
}

// OR or
func orR8(g *GBC, _, r8 int) {
	value := g.Reg.R[A] | g.Reg.R[r8]
	g.Reg.R[A] = value
	g.setF(flagZ, value == 0)

	g.setF(flagN, false)
	g.setF(flagH, false)
	g.setF(flagC, false)
	g.Reg.PC++
}

// OR A,(HL)
func oraHL(g *GBC, _, _ int) {
	value := g.Reg.R[A] | g.Load8(g.Reg.HL())
	g.Reg.R[A] = value
	g.setF(flagZ, value == 0)
	g.setF(flagN, false)
	g.setF(flagH, false)
	g.setF(flagC, false)
	g.Reg.PC++
}

// OR A,u8
func oru8(g *GBC, _, _ int) {
	value := g.Reg.R[A] | g.Load8(g.Reg.PC+1)
	g.Reg.R[A] = value
	g.setF(flagZ, value == 0)
	g.Reg.PC++
	g.setF(flagN, false)
	g.setF(flagH, false)
	g.setF(flagC, false)
	g.Reg.PC++
}

// ADD Addition
func add8(g *GBC, _, r8 int) {
	value := uint16(g.Reg.R[A]) + uint16(g.Reg.R[r8])
	carryBits := uint16(g.Reg.R[A]) ^ uint16(g.Reg.R[r8]) ^ value
	g.Reg.R[A] = byte(value)
	g.setF(flagZ, byte(value) == 0)
	g.setF(flagN, false)
	g.setF(flagH, util.Bit(carryBits, 4))
	g.setF(flagC, util.Bit(carryBits, 8))
	g.Reg.PC++
}

// add hl,r16
func addHL(g *GBC, _, r16 int) {
	value := uint32(g.Reg.HL()) + uint32(g.Reg.R16(r16))
	carryBits := uint32(g.Reg.HL()) ^ uint32(g.Reg.R16(r16)) ^ value
	g.Reg.setHL(uint16(value))
	g.setF(flagN, false)
	g.setF(flagH, util.Bit(carryBits, 12))
	g.setF(flagC, util.Bit(carryBits, 16))
	g.Reg.PC++
}

// ADD A,u8
func addu8(g *GBC, _, _ int) {
	value := uint16(g.Reg.R[A]) + uint16(g.d8Fetch())
	carryBits := uint16(g.Reg.R[A]) ^ uint16(g.d8Fetch()) ^ value
	g.Reg.R[A] = byte(value)
	g.setF(flagZ, byte(value) == 0)
	g.setF(flagN, false)
	g.setF(flagH, util.Bit(carryBits, 4))
	g.setF(flagC, util.Bit(carryBits, 8))
	g.Reg.PC += 2
}

// ADD A,(HL)
func addaHL(g *GBC, _, _ int) {
	value := uint16(g.Reg.R[A]) + uint16(g.Load8(g.Reg.HL()))
	carryBits := uint16(g.Reg.R[A]) ^ uint16(g.Load8(g.Reg.HL())) ^ value
	g.Reg.R[A] = byte(value)
	g.setF(flagZ, byte(value) == 0)
	g.setF(flagN, false)
	g.setF(flagH, util.Bit(carryBits, 4))
	g.setF(flagC, util.Bit(carryBits, 8))
	g.Reg.PC++
}

// ADD SP,i8
func addSPi8(g *GBC, _, _ int) {
	delta := int8(g.Load8(g.Reg.PC + 1))
	value := int32(g.Reg.SP) + int32(delta)
	carryBits := uint32(g.Reg.SP) ^ uint32(delta) ^ uint32(value)
	g.Reg.SP = uint16(value)
	g.setF(flagZ, false)
	g.setF(flagN, false)
	g.setF(flagH, util.Bit(carryBits, 4))
	g.setF(flagC, util.Bit(carryBits, 8))
	g.Reg.PC += 2
}

// complement A Register
func cpl(g *GBC, _, _ int) {
	g.Reg.R[A] = ^g.Reg.R[A]
	g.setF(flagN, true)
	g.setF(flagH, true)
	g.Reg.PC++
}

// extend instruction
func prefixCB(g *GBC, _, _ int) {
	g.Reg.PC++
	g.timer.tick(1 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))
	op := prefixCBs[g.Load8(g.Reg.PC)]
	_, op1, op2, cycle, handler := op.Ins, op.Operand1, op.Operand2, op.Cycle1, op.Handler
	handler(g, op1, op2)

	if cycle > 1 {
		g.timer.tick(uint32(cycle-1) * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))
	}
}

// RLC Rotate n left carry => bit0
func rlc(g *GBC, r8, _ int) {
	value := g.Reg.R[r8]
	bit7 := value >> 7
	value = (value << 1)
	value = util.SetLSB(value, bit7 != 0)
	g.Reg.R[r8] = value

	g.setF(flagZ, value == 0)
	g.setF(flagN, false)
	g.setF(flagH, false)
	g.setF(flagC, bit7 != 0)
	g.Reg.PC++
}

func rlcHL(g *GBC, _, _ int) {
	value := g.Load8(g.Reg.HL())
	g.timer.tick(1 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))
	bit7 := value >> 7
	value = (value << 1)
	value = util.SetLSB(value, bit7 != 0)
	g.Store8(g.Reg.HL(), value)
	g.timer.tick(2 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))

	g.setF(flagZ, value == 0)
	g.setF(flagN, false)
	g.setF(flagH, false)
	g.setF(flagC, bit7 != 0)
	g.Reg.PC++
}

// Rotate register A left.
func rlca(g *GBC, _, _ int) {
	value := g.Reg.R[A]
	bit7 := value >> 7
	value = (value << 1)
	value = util.SetLSB(value, bit7 != 0)
	g.Reg.R[A] = value

	g.setF(flagZ, false)
	g.setF(flagN, false)
	g.setF(flagH, false)
	g.setF(flagC, bit7 != 0)
	g.Reg.PC++
}

// Rotate n right carry => bit7
func rrc(g *GBC, r8, _ int) {
	value := g.Reg.R[r8]
	bit0 := value % 2
	value = (value >> 1)
	value = util.SetMSB(value, bit0 != 0)
	g.Reg.R[r8] = value

	g.setF(flagZ, value == 0)
	g.setF(flagN, false)
	g.setF(flagH, false)
	g.setF(flagC, bit0 != 0)
	g.Reg.PC++
}

func rrcHL(g *GBC, _, _ int) {
	value := g.Load8(g.Reg.HL())
	g.timer.tick(1 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))
	bit0 := value % 2
	value = (value >> 1)
	value = util.SetMSB(value, bit0 != 0)
	g.Store8(g.Reg.HL(), value)
	g.timer.tick(2 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))

	g.setF(flagZ, value == 0)
	g.setF(flagN, false)
	g.setF(flagH, false)
	g.setF(flagC, bit0 != 0)
	g.Reg.PC++
}

// rotate register A right.
func rrca(g *GBC, _, _ int) {
	value, lsb := g.Reg.R[A], util.Bit(g.Reg.R[A], 0)
	value = (value >> 1)
	value = util.SetMSB(value, lsb)
	g.Reg.R[A] = value

	g.setF(flagZ, false)
	g.setF(flagN, false)
	g.setF(flagH, false)
	g.setF(flagC, lsb)
	g.Reg.PC++
}

// RL Rotate n rigth through carry bit7 => bit0
func rl(g *GBC, _, r8 int) {
	carry, value := g.f(flagC), g.Reg.R[r8]
	bit7 := value >> 7
	value = (value << 1)
	value = util.SetLSB(value, carry)
	g.Reg.R[r8] = value

	g.setF(flagZ, value == 0)
	g.setF(flagN, false)
	g.setF(flagH, false)
	g.setF(flagC, bit7 != 0)
	g.Reg.PC++
}

func rlHL(g *GBC, _, _ int) {
	var value, bit7 byte
	carry := g.f(flagC)
	value = g.Load8(g.Reg.HL())
	g.timer.tick(1 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))
	bit7 = value >> 7
	value = (value << 1)
	value = util.SetLSB(value, carry)
	g.Store8(g.Reg.HL(), value)
	g.timer.tick(2 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))

	g.setF(flagZ, value == 0)
	g.setF(flagN, false)
	g.setF(flagH, false)
	g.setF(flagC, bit7 != 0)
	g.Reg.PC++
}

// Rotate register A left through carry.
func rla(g *GBC, _, _ int) {
	carry := g.f(flagC)

	value := g.Reg.R[A]
	bit7 := value >> 7
	value = (value << 1)
	value = util.SetLSB(value, carry)
	g.Reg.R[A] = value

	g.setF(flagZ, false)
	g.setF(flagN, false)
	g.setF(flagH, false)
	g.setF(flagC, bit7 != 0)
	g.Reg.PC++
}

// rr Rotate n right through carry bit0 => bit7
func rr(g *GBC, r8, _ int) {
	value, lsb, carry := g.Reg.R[r8], util.Bit(g.Reg.R[r8], 0), g.f(flagC)
	value >>= 1
	value = util.SetMSB(value, carry)
	g.Reg.R[r8] = value

	g.setF(flagZ, value == 0)
	g.setF(flagN, false)
	g.setF(flagH, false)
	g.setF(flagC, lsb)
	g.Reg.PC++
}

func rrHL(g *GBC, _, _ int) {
	carry := g.f(flagC)
	value := g.Load8(g.Reg.HL())
	g.timer.tick(1 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))
	lsb := util.Bit(value, 0)
	value >>= 1
	value = util.SetMSB(value, carry)
	g.Store8(g.Reg.HL(), value)
	g.timer.tick(2 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))

	g.setF(flagZ, value == 0)
	g.setF(flagN, false)
	g.setF(flagH, false)
	g.setF(flagC, lsb)
	g.Reg.PC++
}

// Shift Left
func sla(g *GBC, r8, _ int) {
	value := g.Reg.R[r8]
	bit7 := value >> 7
	value = (value << 1)
	g.Reg.R[r8] = value

	g.setF(flagZ, value == 0)
	g.setF(flagN, false)
	g.setF(flagH, false)
	g.setF(flagC, bit7 != 0)
	g.Reg.PC++
}

func slaHL(g *GBC, _, _ int) {
	value := g.Load8(g.Reg.HL())
	g.timer.tick(1 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))
	bit7 := value >> 7
	value = (value << 1)
	g.Store8(g.Reg.HL(), value)
	g.timer.tick(2 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))

	g.setF(flagZ, value == 0)
	g.setF(flagN, false)
	g.setF(flagH, false)
	g.setF(flagC, bit7 != 0)
	g.Reg.PC++
}

// Shift Right MSBit dosen't change
func sra(g *GBC, r8, _ int) {
	value, lsb, msb := g.Reg.R[r8], util.Bit(g.Reg.R[r8], 0), util.Bit(g.Reg.R[r8], 7)
	value = (value >> 1)
	value = util.SetMSB(value, msb)
	g.Reg.R[r8] = value

	g.setF(flagZ, value == 0)
	g.setF(flagN, false)
	g.setF(flagH, false)
	g.setF(flagC, lsb)
	g.Reg.PC++
}

func sraHL(g *GBC, operand1, operand2 int) {
	value := g.Load8(g.Reg.HL())
	g.timer.tick(1 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))
	lsb, msb := util.Bit(value, 0), util.Bit(value, 7)
	value = (value >> 1)
	value = util.SetMSB(value, msb)
	g.Store8(g.Reg.HL(), value)
	g.timer.tick(2 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))

	g.setF(flagZ, value == 0)
	g.setF(flagN, false)
	g.setF(flagH, false)
	g.setF(flagC, lsb)
	g.Reg.PC++
}

// SWAP Swap n[5:8] and n[0:4]
func swap(g *GBC, _, r8 int) {
	b := g.Reg.R[r8]
	lower := b & 0b1111
	upper := b >> 4
	g.Reg.R[r8] = (lower << 4) | upper

	g.setF(flagZ, g.Reg.R[r8] == 0)
	g.setF(flagN, false)
	g.setF(flagH, false)
	g.setF(flagC, false)
	g.Reg.PC++
}

func swapHL(g *GBC, _, _ int) {
	data := g.Load8(g.Reg.HL())
	g.timer.tick(1 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))
	data03 := data & 0x0f
	data47 := data >> 4
	value := (data03 << 4) | data47
	g.Store8(g.Reg.HL(), value)
	g.timer.tick(2 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))

	g.setF(flagZ, value == 0)
	g.setF(flagN, false)
	g.setF(flagH, false)
	g.setF(flagC, false)
	g.Reg.PC++
}

// SRL Shift Right MSBit = 0
func srl(g *GBC, r8, _ int) {
	value := g.Reg.R[r8]
	bit0 := value % 2
	value = (value >> 1)
	g.Reg.R[r8] = value

	g.setF(flagZ, value == 0)
	g.setF(flagN, false)
	g.setF(flagH, false)
	g.setF(flagC, bit0 == 1)
	g.Reg.PC++
}

func srlHL(g *GBC, _, _ int) {
	value := g.Load8(g.Reg.HL())
	g.timer.tick(1 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))
	bit0 := value % 2
	value = (value >> 1)
	g.Store8(g.Reg.HL(), value)
	g.timer.tick(2 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))

	g.setF(flagZ, value == 0)
	g.setF(flagN, false)
	g.setF(flagH, false)
	g.setF(flagC, bit0 == 1)
	g.Reg.PC++
}

// BIT Test bit n
func bit(g *GBC, bit, r8 int) {
	value := util.Bit(g.Reg.R[r8], bit)
	g.setF(flagZ, !value)
	g.setF(flagN, false)
	g.setF(flagH, true)
	g.Reg.PC++
}

func bitHL(g *GBC, bit, _ int) {
	value := util.Bit(g.Load8(g.Reg.HL()), bit)
	g.setF(flagZ, !value)
	g.setF(flagN, false)
	g.setF(flagH, true)
	g.Reg.PC++
}

func res(g *GBC, bit, r8 int) {
	mask := ^(byte(1) << bit)
	g.Reg.R[r8] &= mask
	g.Reg.PC++
}

func resHL(g *GBC, bit, _ int) {
	mask := ^(byte(1) << bit)
	value := g.Load8(g.Reg.HL()) & mask
	g.timer.tick(1 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))
	g.Store8(g.Reg.HL(), value)
	g.timer.tick(2 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))
	g.Reg.PC++
}

func set(g *GBC, bit, r8 int) {
	mask := byte(1) << bit
	g.Reg.R[r8] |= mask
	g.Reg.PC++
}

func setHL(g *GBC, bit, _ int) {
	mask := byte(1) << bit
	value := g.Load8(g.Reg.HL()) | mask
	g.timer.tick(1 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))
	g.Store8(g.Reg.HL(), value)
	g.timer.tick(2 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))
	g.Reg.PC++
}

// push af
func pushAF(g *GBC, _, _ int) {
	g.timer.tick(1 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))
	g.push(g.Reg.R[A])
	g.timer.tick(1 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))
	g.push(g.Reg.R[F] & 0xf0)
	g.Reg.PC++
	g.timer.tick(2 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))
}

// push r16
func push(g *GBC, r0, r1 int) {
	g.timer.tick(1 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))
	g.push(g.Reg.R[r0])
	g.timer.tick(1 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))
	g.push(g.Reg.R[r1])
	g.Reg.PC++
	g.timer.tick(2 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))
}

func popAF(g *GBC, _, _ int) {
	g.Reg.R[F] = g.pop() & 0xf0
	g.timer.tick(1 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))
	g.Reg.R[A] = g.pop()
	g.Reg.PC++
	g.timer.tick(2 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))
}

func pop(g *GBC, r0, r1 int) {
	g.Reg.R[r0] = g.pop()
	g.timer.tick(1 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))
	g.Reg.R[r1] = g.pop()
	g.Reg.PC++
	g.timer.tick(2 * 4 >> uint32(util.Bool2U64(g.DoubleSpeed)))
}

// SUB subtract
func sub8(g *GBC, _, r8 int) {
	value := g.Reg.R[A] - g.Reg.R[r8]
	carryBits := g.Reg.R[A] ^ g.Reg.R[r8] ^ value
	g.setCSub(g.Reg.R[A], g.Reg.R[r8])
	g.Reg.R[A] = value
	g.setF(flagZ, value == 0)
	g.setF(flagN, true)
	g.setF(flagH, util.Bit(carryBits, 4))
	g.Reg.PC++
}

// SUB A,(HL)
func subaHL(g *GBC, _, _ int) {
	value := g.Reg.R[A] - g.Load8(g.Reg.HL())
	carryBits := g.Reg.R[A] ^ g.Load8(g.Reg.HL()) ^ value
	g.setCSub(g.Reg.R[A], g.Load8(g.Reg.HL()))
	g.Reg.R[A] = value
	g.setF(flagZ, value == 0)
	g.setF(flagN, true)
	g.setF(flagH, util.Bit(carryBits, 4))
	g.Reg.PC++
}

// SUB A,u8
func subu8(g *GBC, _, _ int) {
	value := g.Reg.R[A] - g.d8Fetch()
	carryBits := g.Reg.R[A] ^ g.d8Fetch() ^ value
	g.setCSub(g.Reg.R[A], g.d8Fetch())
	g.Reg.R[A] = value
	g.setF(flagZ, value == 0)
	g.setF(flagN, true)
	g.setF(flagH, util.Bit(carryBits, 4))
	g.Reg.PC += 2
}

// Rotate register A right through carry.
func rra(g *GBC, _, _ int) {
	carry := g.f(flagC)
	regA := g.Reg.R[A]
	g.setF(flagC, util.Bit(regA, 0))
	if carry {
		regA = (1 << 7) | (regA >> 1)
	} else {
		regA = (0 << 7) | (regA >> 1)
	}
	g.Reg.R[A] = regA
	g.setF(flagZ, false)
	g.setF(flagN, false)
	g.setF(flagH, false)
	g.Reg.PC++
}

// ADC Add the value n8 plus the carry flag to A
func adc8(g *GBC, _, op int) {
	var carry, value, value4 byte
	var value16 uint16
	if g.f(flagC) {
		carry = 1
	}

	value = g.Reg.R[op] + carry + g.Reg.R[A]
	value4 = (g.Reg.R[op] & 0b1111) + carry + (g.Reg.R[A] & 0b1111)
	value16 = uint16(g.Reg.R[op]) + uint16(carry) + uint16(g.Reg.R[A])
	g.Reg.R[A] = value

	g.setF(flagZ, value == 0)
	g.setF(flagN, false)
	g.setF(flagH, util.Bit(value4, 4))
	g.setF(flagC, util.Bit(value16, 8))
	g.Reg.PC++
}

// ADC A,(HL)
func adcaHL(g *GBC, _, _ int) {
	var carry, value, value4 byte
	var value16 uint16
	if g.f(flagC) {
		carry = 1
	}

	data := g.Load8(g.Reg.HL())
	value = data + carry + g.Reg.R[A]
	value4 = (data & 0x0f) + carry + (g.Reg.R[A] & 0b1111)
	value16 = uint16(data) + uint16(g.Reg.R[A]) + uint16(carry)

	g.Reg.R[A] = value
	g.setF(flagZ, value == 0)
	g.setF(flagN, false)
	g.setF(flagH, util.Bit(value4, 4))
	g.setF(flagC, util.Bit(value16, 8))
	g.Reg.PC++
}

// ADC A,u8
func adcu8(g *GBC, _, _ int) {
	var carry, value, value4 byte
	var value16 uint16
	if g.f(flagC) {
		carry = 1
	}

	data := g.d8Fetch()
	value = data + carry + g.Reg.R[A]
	value4 = (data & 0x0f) + carry + (g.Reg.R[A] & 0b1111)
	value16 = uint16(data) + uint16(g.Reg.R[A]) + uint16(carry)
	g.Reg.PC++

	g.Reg.R[A] = value
	g.setF(flagZ, value == 0)
	g.setF(flagN, false)
	g.setF(flagH, util.Bit(value4, 4))
	g.setF(flagC, util.Bit(value16, 8))
	g.Reg.PC++
}

// SBC Subtract the value n8 and the carry flag from A

func sbc8(g *GBC, _, op int) {
	var carry, value, value4 byte
	var value16 uint16
	if g.f(flagC) {
		carry = 1
	}

	value = g.Reg.R[A] - (g.Reg.R[op] + carry)
	value4 = (g.Reg.R[A] & 0b1111) - ((g.Reg.R[op] & 0b1111) + carry)
	value16 = uint16(g.Reg.R[A]) - (uint16(g.Reg.R[op]) + uint16(carry))
	g.Reg.R[A] = value

	g.setF(flagZ, value == 0)
	g.setF(flagN, true)
	g.setF(flagH, util.Bit(value4, 4))
	g.setF(flagC, util.Bit(value16, 8))
	g.Reg.PC++
}

func (g *GBC) SBC(_, op2 int) {
	var carry, value, value4 byte
	var value16 uint16
	if g.f(flagC) {
		carry = 1
	}

	switch op2 {
	case OP_HL_PAREN:
		data := g.Load8(g.Reg.HL())
		value = g.Reg.R[A] - (data + carry)
		value4 = (g.Reg.R[A] & 0b1111) - ((data & 0x0f) + carry)
		value16 = uint16(g.Reg.R[A]) - (uint16(data) + uint16(carry))
	case OP_d8:
		data := g.d8Fetch()
		value = g.Reg.R[A] - (data + carry)
		value4 = (g.Reg.R[A] & 0b1111) - ((data & 0x0f) + carry)
		value16 = uint16(g.Reg.R[A]) - (uint16(data) + uint16(carry))
		g.Reg.PC++
	}
	g.Reg.R[A] = value

	g.setF(flagZ, value == 0)
	g.setF(flagN, true)
	g.setF(flagH, util.Bit(value4, 4))
	g.setF(flagC, util.Bit(value16, 8))
	g.Reg.PC++
}

// DAA Decimal adjust
func daa(g *GBC, _, _ int) {
	a := uint8(g.Reg.R[A])
	// ref: https://forums.nesdev.com/viewtopic.php?f=20&t=15944
	if !g.f(flagN) {
		if g.f(flagC) || a > 0x99 {
			a += 0x60
			g.setF(flagC, true)
		}
		if g.f(flagH) || (a&0x0f) > 0x09 {
			a += 0x06
		}
	} else {
		if g.f(flagC) {
			a -= 0x60
		}
		if g.f(flagH) {
			a -= 0x06
		}
	}

	g.Reg.R[A] = a
	g.setF(flagZ, a == 0)
	g.setF(flagH, false)
	g.Reg.PC++
}

// push present address and jump to vector address
func rst(g *GBC, addr, _ int) {
	g.Reg.PC++
	g.pushPC()
	g.Reg.PC = uint16(addr)
}

func scf(g *GBC, _, _ int) {
	g.setF(flagN, false)
	g.setF(flagH, false)
	g.setF(flagC, true)
	g.Reg.PC++
}

// CCF Complement Carry Flag
func ccf(g *GBC, _, _ int) {
	g.setF(flagN, false)
	g.setF(flagH, false)
	g.setF(flagC, !g.f(flagC))
	g.Reg.PC++
}
