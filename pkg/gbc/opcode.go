package gbc

import (
	"fmt"

	"github.com/pokemium/worldwide/pkg/gbc/scheduler"
	"github.com/pokemium/worldwide/pkg/util"
)

func (g *GBC) fixCycles(cycles uint32) uint32 {
	return cycles * 4 >> util.Bool2U32(g.DoubleSpeed)
}

func (g *GBC) a16Fetch() uint16 {
	lower, upper := uint16(g.Load8(g.Inst.PC+1)), uint16(g.Load8(g.Inst.PC+2))
	g.Reg.PC += 2
	return (upper << 8) | lower
}

func (g *GBC) a16FetchJP() uint16 {
	lower := uint16(g.Load8(g.Inst.PC + 1)) // M = 1: nn read: memory access for low byte
	g.timer.tick(g.fixCycles(1))
	upper := uint16(g.Load8(g.Inst.PC + 2)) // M = 2: nn read: memory access for high byte
	g.timer.tick(g.fixCycles(1))
	g.Reg.PC += 2
	value := (upper << 8) | lower
	return value
}

func (g *GBC) d8Fetch() byte {
	value := g.Load8(g.Inst.PC + 1)
	g.Reg.PC++
	return value
}

func (g *GBC) fetchAtHL() byte {
	value := g.Load8(g.Reg.HL())
	g.timer.tick(g.fixCycles(1))
	return value
}

// LD R8,R8
func ld8r(g *GBC, op1, op2 int) {
	g.Reg.R[op1] = g.Reg.R[op2]
}

// LD r8, mem[r16]
func ld8m(g *GBC, r8, r16 int) {
	g.Reg.R[r8] = g.Load8(g.Reg.R16(r16))
}

// LD r8, mem[imm]
func ld8i(g *GBC, r8, _ int) {
	g.Reg.R[r8] = g.d8Fetch()
}

// LD A, (u16)
func ldau16(g *GBC, operand1, operand2 int) {
	g.Reg.R[A] = g.Load8(g.a16FetchJP())
}

// LD A,(FF00+C)
func op0xf2(g *GBC, operand1, operand2 int) {
	addr := 0xff00 + uint16(g.Reg.R[C])
	g.Reg.R[A] = g.loadIO(byte(addr))
}

// LD (HL),u8
func op0x36(g *GBC, operand1, operand2 int) {
	value := g.d8Fetch()
	g.timer.tick(g.fixCycles(1))
	g.Store8(g.Reg.HL(), value)
}

// LD (HL),R8
func ldHLR8(g *GBC, unused, op int) {
	g.Store8(g.Reg.HL(), g.Reg.R[op])
}

// LD (u16),SP
func op0x08(g *GBC, operand1, operand2 int) {
	addr := g.a16Fetch()
	upper, lower := byte(g.Reg.SP>>8), byte(g.Reg.SP)
	g.Store8(addr, lower)
	g.Store8(addr+1, upper)
}

// LD (u16),A
func op0xea(g *GBC, operand1, operand2 int) {
	g.Store8(g.a16FetchJP(), g.Reg.R[A])
}

// LD r16, u16
func ld16i(g *GBC, r16, _ int) {
	g.Reg.setR16(r16, g.a16Fetch())
}

// LD HL,SP+i8
func op0xf8(g *GBC, operand1, operand2 int) {
	lhs, rhs := g.Reg.SP, int8(g.d8Fetch())
	value := int32(lhs) + int32(rhs)
	carryBits := uint32(lhs) ^ uint32(rhs) ^ uint32(value)
	g.Reg.setHL(uint16(value))
	g.setZNHC(false, false, util.Bit(carryBits, 4), util.Bit(carryBits, 8))
}

// LD SP,HL
func op0xf9(g *GBC, operand1, operand2 int) {
	g.Reg.SP = g.Reg.HL()
}

// LD (FF00+C),A
func op0xe2(g *GBC, operand1, operand2 int) {
	addr := 0xff00 + uint16(g.Reg.R[C])
	g.Store8(addr, g.Reg.R[A])
}

func ldm16r(g *GBC, r16, r8 int) {
	g.Store8(g.Reg.R16(r16), g.Reg.R[r8])
}

// LD ($FF00+a8),A
func op0xe0(g *GBC, operand1, operand2 int) {
	addr := 0xff00 + uint16(g.d8Fetch())
	g.timer.tick(g.fixCycles(1))
	g.storeIO(byte(addr), g.Reg.R[A])
}

// LD A,($FF00+a8)
func op0xf0(g *GBC, operand1, operand2 int) {
	addr := 0xff00 + uint16(g.d8Fetch())
	g.timer.tick(g.fixCycles(1))
	g.Reg.R[A] = g.loadIO(byte(addr))
}

// nop No operation
func nop(g *GBC, _, _ int) {}

func inc8(g *GBC, r8, _ int) {
	value := g.Reg.R[r8] + 1
	carryBits := g.Reg.R[r8] ^ 1 ^ value
	g.Reg.R[r8] = value
	g.setZNH(value == 0, false, util.Bit(carryBits, 4))
}

func inc16(g *GBC, r16, _ int) {
	g.Reg.setR16(r16, g.Reg.R16(r16)+1)
}

func incHL(g *GBC, _, _ int) {
	hl := g.fetchAtHL()
	value := hl + 1
	carryBits := hl ^ 1 ^ value
	g.Store8(g.Reg.HL(), value)
	g.setZNH(value == 0, false, util.Bit(carryBits, 4))
}

func dec8(g *GBC, r8, _ int) {
	value := g.Reg.R[r8] - 1
	carryBits := g.Reg.R[r8] ^ 1 ^ value
	g.Reg.R[r8] = value
	g.setZNH(value == 0, true, util.Bit(carryBits, 4))
}

func dec16(g *GBC, r16, _ int) {
	g.Reg.setR16(r16, g.Reg.R16(r16)-1)
}

func decHL(g *GBC, _, _ int) {
	hl := g.fetchAtHL()
	value := hl - 1
	carryBits := hl ^ 1 ^ value
	g.Store8(g.Reg.HL(), value)
	g.setZNH(value == 0, true, util.Bit(carryBits, 4))
}

// --------- JR ---------

// jr i8
func jr(g *GBC, _, _ int) {
	_jr(g, int8(g.d8Fetch()))
}

func _jr(g *GBC, delta int8) {
	g.Reg.PC = uint16(int32(g.Reg.PC) + int32(delta))
	g.timer.tick(g.fixCycles(3))
}

// jr cc,i8
func jrcc(g *GBC, cc, _ int) {
	delta := int8(g.d8Fetch())
	if g.f(cc) {
		_jr(g, delta)
	} else {
		g.timer.tick(g.fixCycles(2))
	}
}

// jr ncc,i8 (ncc = not cc)
func jrncc(g *GBC, cc, _ int) {
	delta := int8(g.d8Fetch())
	if !g.f(cc) {
		_jr(g, delta)
	} else {
		g.timer.tick(g.fixCycles(2))
	}
}

func halt(g *GBC, _, _ int) {
	if g.IO[IEIO]&g.IO[IFIO]&0x1f == 0 {
		g.Halt = true
	}
}

// stop GBC
func stop(g *GBC, _, _ int) {
	g.Reg.PC++
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
	g.Reg.R[A] ^= g.Reg.R[r8]
	g.setZNHC(g.Reg.R[A] == 0, false, false, false)
}

// XOR A,(HL)
func xoraHL(g *GBC, _, _ int) {
	g.Reg.R[A] ^= g.Load8(g.Reg.HL())
	g.setZNHC(g.Reg.R[A] == 0, false, false, false)
}

// XOR A,u8
func xoru8(g *GBC, _, _ int) {
	g.Reg.R[A] ^= g.d8Fetch()
	g.setZNHC(g.Reg.R[A] == 0, false, false, false)
}

// jp u16
func jp(g *GBC, _, _ int) {
	g.Reg.PC = g.a16FetchJP()
}

func jpcc(g *GBC, cc, _ int) {
	dst := g.a16FetchJP()
	if g.f(cc) {
		g.Reg.PC = dst
		g.timer.tick(g.fixCycles(1))
	}
}

func jpncc(g *GBC, cc, _ int) {
	dst := g.a16FetchJP()
	if !g.f(cc) {
		g.Reg.PC = dst
		g.timer.tick(g.fixCycles(1))
	}
}

func jpHL(g *GBC, _, _ int) {
	g.Reg.PC = g.Reg.HL()
}

// Return
func ret(g *GBC, _, _ int) {
	g.popPC()
}

func retcc(g *GBC, cc, _ int) {
	if g.f(cc) {
		g.popPC()
		g.timer.tick(g.fixCycles(3))
	}
}

// not retcc
func retncc(g *GBC, cc, _ int) {
	if !g.f(cc) {
		g.popPC()
		g.timer.tick(g.fixCycles(3))
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
	dest := g.a16FetchJP()
	_call(g, dest)
}

func _call(g *GBC, dest uint16) {
	g.timer.tick(g.fixCycles(1))
	g.pushPCCALL()
	g.timer.tick(g.fixCycles(1))
	g.Reg.PC = dest
}

func callcc(g *GBC, cc, _ int) {
	dest := g.a16FetchJP()
	if g.f(cc) {
		_call(g, dest)
		return
	}
	g.timer.tick(g.fixCycles(1))
}

func callncc(g *GBC, cc, _ int) {
	dest := g.a16FetchJP()
	if !g.f(cc) {
		_call(g, dest)
		return
	}
	g.timer.tick(g.fixCycles(1))
}

// DI Disable Interrupt
func di(g *GBC, _, _ int) {
	g.setInterrupts(false)
}

// EI Enable Interrupt
func ei(g *GBC, _, _ int) {
	g.setInterrupts(true)
}

// CP Compare
func cp(g *GBC, _, r8 int) {
	lhs, rhs, value := g.Reg.R[A], g.Reg.R[r8], g.Reg.R[A]-g.Reg.R[r8]
	carryBits := lhs ^ rhs ^ value
	g.setZNHC(value == 0, true, util.Bit(carryBits, 4), subC(lhs, rhs))
}

// CP A,(HL)
func cpaHL(g *GBC, _, _ int) {
	lhs, rhs := g.Reg.R[A], g.Load8(g.Reg.HL())
	value := lhs - rhs
	carryBits := lhs ^ rhs ^ value
	g.setZNHC(value == 0, true, util.Bit(carryBits, 4), subC(lhs, rhs))
}

// CP A,u8
func cpu8(g *GBC, _, _ int) {
	lhs, rhs := g.Reg.R[A], g.d8Fetch()
	value := lhs - rhs
	carryBits := lhs ^ rhs ^ value
	g.setZNHC(value == 0, true, util.Bit(carryBits, 4), subC(lhs, rhs))
}

// AND A,r8
func and8(g *GBC, _, r8 int) {
	g.Reg.R[A] &= g.Reg.R[r8]
	g.setZNHC(g.Reg.R[A] == 0, false, true, false)
}

// AND A,(HL)
func andaHL(g *GBC, _, _ int) {
	g.Reg.R[A] &= g.Load8(g.Reg.HL())
	g.setZNHC(g.Reg.R[A] == 0, false, true, false)
}

// AND A,u8
func andu8(g *GBC, _, _ int) {
	g.Reg.R[A] &= g.d8Fetch()
	g.setZNHC(g.Reg.R[A] == 0, false, true, false)
}

// OR A,r8
func or8(g *GBC, _, r8 int) {
	g.Reg.R[A] |= g.Reg.R[r8]
	g.setZNHC(g.Reg.R[A] == 0, false, false, false)
}

// OR A,(HL)
func oraHL(g *GBC, _, _ int) {
	g.Reg.R[A] |= g.Load8(g.Reg.HL())
	g.setZNHC(g.Reg.R[A] == 0, false, false, false)
}

// OR A,u8
func oru8(g *GBC, _, _ int) {
	g.Reg.R[A] |= g.d8Fetch()
	g.setZNHC(g.Reg.R[A] == 0, false, false, false)
}

// ADD A,r8
func add8(g *GBC, _, r8 int) {
	lhs, rhs := g.Reg.R[A], g.Reg.R[r8]
	value := uint16(lhs) + uint16(rhs)
	carryBits := uint16(lhs) ^ uint16(rhs) ^ value
	g.Reg.R[A] = byte(value)
	g.setZNHC(byte(value) == 0, false, util.Bit(carryBits, 4), util.Bit(carryBits, 8))
}

// ADD HL,r16
func addHL(g *GBC, _, r16 int) {
	lhs, rhs := g.Reg.HL(), g.Reg.R16(r16)
	value := uint32(lhs) + uint32(rhs)
	carryBits := uint32(lhs) ^ uint32(rhs) ^ value
	g.Reg.setHL(uint16(value))
	g.setNHC(false, util.Bit(carryBits, 12), util.Bit(carryBits, 16))
}

// ADD A,u8
func addu8(g *GBC, _, _ int) {
	lhs, rhs := g.Reg.R[A], g.d8Fetch()
	value := uint16(lhs) + uint16(rhs)
	carryBits := uint16(lhs) ^ uint16(rhs) ^ value
	g.Reg.R[A] = byte(value)
	g.setZNHC(byte(value) == 0, false, util.Bit(carryBits, 4), util.Bit(carryBits, 8))
}

// ADD A,(HL)
func addaHL(g *GBC, _, _ int) {
	lhs, rhs := g.Reg.R[A], g.Load8(g.Reg.HL())
	value := uint16(lhs) + uint16(rhs)
	carryBits := uint16(lhs) ^ uint16(rhs) ^ value
	g.Reg.R[A] = byte(value)
	g.setZNHC(byte(value) == 0, false, util.Bit(carryBits, 4), util.Bit(carryBits, 8))
}

// ADD SP,i8
func addSPi8(g *GBC, _, _ int) {
	rhs := int8(g.d8Fetch())
	value := int32(g.Reg.SP) + int32(rhs)
	carryBits := uint32(g.Reg.SP) ^ uint32(rhs) ^ uint32(value)
	g.Reg.SP = uint16(value)
	g.setZNHC(false, false, util.Bit(carryBits, 4), util.Bit(carryBits, 8))
}

// complement A Register
func cpl(g *GBC, _, _ int) {
	g.Reg.R[A] = ^g.Reg.R[A]
	g.setNH(true, true)
}

// extend instruction
func prefixCB(g *GBC, _, _ int) {
	g.timer.tick(g.fixCycles(1))
	inst := gbz80instsCb[g.d8Fetch()]
	op1, op2, cycle, handler := inst.Operand1, inst.Operand2, inst.Cycle1, inst.Handler
	handler(g, op1, op2)

	if cycle > 1 {
		g.timer.tick(g.fixCycles(uint32(cycle - 1)))
	}
}

// RLC Rotate n left carry => bit0
func rlc(g *GBC, r8, _ int) {
	value, bit7 := g.Reg.R[r8], util.Bit(g.Reg.R[r8], 7)
	g.Reg.R[r8] = util.SetLSB(value<<1, bit7)
	g.setZNHC(g.Reg.R[r8] == 0, false, false, bit7)
}

func rlcHL(g *GBC, _, _ int) {
	value := g.fetchAtHL()
	bit7 := util.Bit(value, 7)
	value = util.SetLSB(value<<1, bit7)
	g.Store8(g.Reg.HL(), value)
	g.setZNHC(value == 0, false, false, bit7)
}

// Rotate register A left.
func rlca(g *GBC, _, _ int) {
	value, bit7 := g.Reg.R[A], util.Bit(g.Reg.R[A], 7)
	g.Reg.R[A] = util.SetLSB(value<<1, bit7)
	g.setZNHC(false, false, false, bit7)
}

// Rotate n right carry => bit7
func rrc(g *GBC, r8, _ int) {
	value, bit0 := g.Reg.R[r8], util.Bit(g.Reg.R[r8], 0)
	g.Reg.R[r8] = util.SetMSB(value>>1, bit0)
	g.setZNHC(g.Reg.R[r8] == 0, false, false, bit0)
}

func rrcHL(g *GBC, _, _ int) {
	value := g.fetchAtHL()
	bit0 := util.Bit(value, 0)
	value = util.SetMSB(value>>1, bit0)
	g.Store8(g.Reg.HL(), value)
	g.setZNHC(value == 0, false, false, bit0)
}

// rotate register A right.
func rrca(g *GBC, _, _ int) {
	value, bit0 := g.Reg.R[A], util.Bit(g.Reg.R[A], 0)
	g.Reg.R[A] = util.SetMSB(value>>1, bit0)
	g.setZNHC(false, false, false, bit0)
}

// RL Rotate n rigth through carry bit7 => bit0
func rl(g *GBC, _, r8 int) {
	carry, value, bit7 := g.f(flagC), g.Reg.R[r8], util.Bit(g.Reg.R[r8], 7)
	g.Reg.R[r8] = util.SetLSB(value<<1, carry)
	g.setZNHC(g.Reg.R[r8] == 0, false, false, bit7)
}

func rlHL(g *GBC, _, _ int) {
	carry, value := g.f(flagC), g.fetchAtHL()
	bit7 := util.Bit(value, 7)
	value = util.SetLSB(value<<1, carry)
	g.Store8(g.Reg.HL(), value)
	g.setZNHC(value == 0, false, false, bit7)
}

// Rotate register A left through carry.
func rla(g *GBC, _, _ int) {
	carry, value, bit7 := g.f(flagC), g.Reg.R[A], util.Bit(g.Reg.R[A], 7)
	g.Reg.R[A] = util.SetLSB(value<<1, carry)
	g.setZNHC(false, false, false, bit7)
}

// rr Rotate n right through carry bit0 => bit7
func rr(g *GBC, r8, _ int) {
	value, bit0, carry := g.Reg.R[r8], util.Bit(g.Reg.R[r8], 0), g.f(flagC)
	g.Reg.R[r8] = util.SetMSB(value>>1, carry)
	g.setZNHC(g.Reg.R[r8] == 0, false, false, bit0)
}

func rrHL(g *GBC, _, _ int) {
	carry, value := g.f(flagC), g.fetchAtHL()
	bit0 := util.Bit(value, 0)
	value = util.SetMSB(value>>1, carry)
	g.Store8(g.Reg.HL(), value)
	g.setZNHC(value == 0, false, false, bit0)
}

// Shift Left
func sla(g *GBC, r8, _ int) {
	value, bit7 := g.Reg.R[r8], util.Bit(g.Reg.R[r8], 7)
	g.Reg.R[r8] = value << 1
	g.setZNHC(g.Reg.R[r8] == 0, false, false, bit7)
}

func slaHL(g *GBC, _, _ int) {
	value := g.fetchAtHL()
	bit7 := util.Bit(value, 7)
	hl := (value << 1)
	g.Store8(g.Reg.HL(), hl)
	g.setZNHC(hl == 0, false, false, bit7)
}

// Shift Right MSBit dosen't change
func sra(g *GBC, r8, _ int) {
	value, bit0, bit7 := g.Reg.R[r8], util.Bit(g.Reg.R[r8], 0), util.Bit(g.Reg.R[r8], 7)
	g.Reg.R[r8] = util.SetMSB(value>>1, bit7)
	g.setZNHC(g.Reg.R[r8] == 0, false, false, bit0)
}

func sraHL(g *GBC, operand1, operand2 int) {
	value := g.fetchAtHL()
	bit0, bit7 := util.Bit(value, 0), util.Bit(value, 7)
	value = util.SetMSB(value>>1, bit7)
	g.Store8(g.Reg.HL(), value)
	g.setZNHC(value == 0, false, false, bit0)
}

// SWAP Swap b[5:8] and b[0:4]
func swap(g *GBC, _, r8 int) {
	b := g.Reg.R[r8]
	upper, lower := b>>4, b&0b1111
	g.Reg.R[r8] = (lower << 4) | upper
	g.setZNHC(g.Reg.R[r8] == 0, false, false, false)
}

func swapHL(g *GBC, _, _ int) {
	b := g.fetchAtHL()
	upper, lower := b>>4, b&0b1111
	value := (lower << 4) | upper
	g.Store8(g.Reg.HL(), value)
	g.setZNHC(value == 0, false, false, false)
}

// SRL Shift Right MSBit = 0
func srl(g *GBC, r8, _ int) {
	value, bit0 := g.Reg.R[r8], util.Bit(g.Reg.R[r8], 0)
	g.Reg.R[r8] = value >> 1
	g.setZNHC(g.Reg.R[r8] == 0, false, false, bit0)
}

func srlHL(g *GBC, _, _ int) {
	value := g.fetchAtHL()
	bit0 := util.Bit(value, 0)
	value = (value >> 1)
	g.Store8(g.Reg.HL(), value)
	g.setZNHC(value == 0, false, false, bit0)
}

// BIT Test bit n
func bit(g *GBC, bit, r8 int) {
	isSet := util.Bit(g.Reg.R[r8], bit)
	g.setZNH(!isSet, false, true)
}

func bitHL(g *GBC, bit, _ int) {
	isSet := util.Bit(g.Load8(g.Reg.HL()), bit)
	g.setZNH(!isSet, false, true)
}

func res(g *GBC, bit, r8 int) {
	g.Reg.R[r8] &= ^(byte(1) << bit)
}

func resHL(g *GBC, bit, _ int) {
	mask := ^(byte(1) << bit)
	g.Store8(g.Reg.HL(), g.fetchAtHL()&mask)
}

func set(g *GBC, bit, r8 int) {
	g.Reg.R[r8] |= byte(1) << bit
}

func setHL(g *GBC, bit, _ int) {
	mask := byte(1) << bit
	g.Store8(g.Reg.HL(), g.fetchAtHL()|mask)
}

// push af
func pushAF(g *GBC, _, _ int) {
	g.timer.tick(g.fixCycles(1))
	g.push(g.Reg.R[A])
	g.timer.tick(g.fixCycles(1))
	g.push(g.Reg.R[F] & 0xf0)
}

// push r16
func push(g *GBC, r0, r1 int) {
	g.timer.tick(g.fixCycles(1))
	g.push(g.Reg.R[r0])
	g.timer.tick(g.fixCycles(1))
	g.push(g.Reg.R[r1])
}

func popAF(g *GBC, _, _ int) {
	g.Reg.R[F] = g.pop() & 0xf0
	g.timer.tick(g.fixCycles(1))
	g.Reg.R[A] = g.pop()
}

func pop(g *GBC, r0, r1 int) {
	g.Reg.R[r0] = g.pop()
	g.timer.tick(g.fixCycles(1))
	g.Reg.R[r1] = g.pop()
}

// SUB subtract
func sub8(g *GBC, _, r8 int) {
	lhs, rhs := g.Reg.R[A], g.Reg.R[r8]
	value := lhs - rhs
	carryBits := lhs ^ rhs ^ value
	g.Reg.R[A] = value
	g.setZNHC(value == 0, true, util.Bit(carryBits, 4), subC(lhs, rhs))
}

// SUB A,(HL)
func subaHL(g *GBC, _, _ int) {
	lhs, rhs := g.Reg.R[A], g.Load8(g.Reg.HL())
	value := lhs - rhs
	carryBits := lhs ^ rhs ^ value
	g.Reg.R[A] = value
	g.setZNHC(value == 0, true, util.Bit(carryBits, 4), subC(lhs, rhs))
}

// SUB A,u8
func subu8(g *GBC, _, _ int) {
	lhs, rhs := g.Reg.R[A], g.d8Fetch()
	value := lhs - rhs
	carryBits := lhs ^ rhs ^ value
	g.Reg.R[A] = value
	g.setZNHC(value == 0, true, util.Bit(carryBits, 4), subC(lhs, rhs))
}

// Rotate register A right through carry.
func rra(g *GBC, _, _ int) {
	oldC, a := g.f(flagC), g.Reg.R[A]
	newC := util.Bit(a, 0)
	if oldC {
		a = (1 << 7) | (a >> 1)
	} else {
		a = (0 << 7) | (a >> 1)
	}
	g.Reg.R[A] = a
	g.setZNHC(false, false, false, newC)
}

// ADC Add the value n8 plus the carry flag to A
// ADC A, r8
func adc8(g *GBC, _, op int) {
	lhs, rhs, carry := g.Reg.R[A], g.Reg.R[op], util.Bool2U8(g.f(flagC))
	value := lhs + rhs + carry
	value4, value16 := (lhs&0b1111)+(rhs&0b1111)+carry, uint16(lhs)+uint16(rhs)+uint16(carry)
	g.Reg.R[A] = value
	g.setZNHC(value == 0, false, util.Bit(value4, 4), util.Bit(value16, 8))
}

// ADC A,(HL)
func adcaHL(g *GBC, _, _ int) {
	lhs, rhs, carry := g.Load8(g.Reg.HL()), g.Reg.R[A], util.Bool2U8(g.f(flagC))
	value := lhs + carry + rhs
	value4, value16 := (lhs&0x0f)+carry+(rhs&0b1111), uint16(lhs)+uint16(rhs)+uint16(carry)
	g.Reg.R[A] = value
	g.setZNHC(value == 0, false, util.Bit(value4, 4), util.Bit(value16, 8))
}

// ADC A,u8
func adcu8(g *GBC, _, _ int) {
	lhs, rhs, carry := g.d8Fetch(), g.Reg.R[A], util.Bool2U8(g.f(flagC))
	value := lhs + carry + rhs
	value4, value16 := (lhs&0x0f)+carry+(rhs&0b1111), uint16(lhs)+uint16(rhs)+uint16(carry)
	g.Reg.R[A] = value
	g.setZNHC(value == 0, false, util.Bit(value4, 4), util.Bit(value16, 8))
}

// SBC Subtract the value n8 and the carry flag from A

func sbc8(g *GBC, _, op int) {
	lhs, rhs, carry := g.Reg.R[A], g.Reg.R[op], util.Bool2U8(g.f(flagC))
	value := lhs - (rhs + carry)
	value4, value16 := (lhs&0b1111)-((rhs&0b1111)+carry), uint16(lhs)-(uint16(rhs)+uint16(carry))
	g.Reg.R[A] = value
	g.setZNHC(value == 0, true, util.Bit(value4, 4), util.Bit(value16, 8))
}

// SBC A,(HL)
func sbcaHL(g *GBC, _, _ int) {
	lhs, rhs, carry := g.Reg.R[A], g.Load8(g.Reg.HL()), util.Bool2U8(g.f(flagC))
	value := lhs - (rhs + carry)
	value4, value16 := (lhs&0b1111)-((rhs&0x0f)+carry), uint16(lhs)-(uint16(rhs)+uint16(carry))
	g.Reg.R[A] = value
	g.setZNHC(value == 0, true, util.Bit(value4, 4), util.Bit(value16, 8))
}

// SBC A,u8
func sbcu8(g *GBC, _, _ int) {
	lhs, rhs, carry := g.Reg.R[A], g.d8Fetch(), util.Bool2U8(g.f(flagC))
	value := lhs - (rhs + carry)
	value4, value16 := (lhs&0b1111)-((rhs&0x0f)+carry), uint16(lhs)-(uint16(rhs)+uint16(carry))
	g.Reg.R[A] = value
	g.setZNHC(value == 0, true, util.Bit(value4, 4), util.Bit(value16, 8))
}

// DAA Decimal adjust (ref: https://forums.nesdev.com/viewtopic.php?f=20&t=15944)
func daa(g *GBC, _, _ int) {
	a := g.Reg.R[A]
	if !g.f(flagN) {
		if g.f(flagC) || a > 0x99 {
			a += 0x60
			g.setF(flagC, true)
		}
		a += 0x06 * util.Bool2U8(g.f(flagH) || (a&0x0f) > 0x09)
	} else {
		a -= 0x60 * util.Bool2U8(g.f(flagC))
		a -= 0x06 * util.Bool2U8(g.f(flagH))
	}
	g.Reg.R[A] = a

	g.setF(flagZ, a == 0)
	g.setF(flagH, false)
}

// push present address and jump to vector address
func rst(g *GBC, addr, _ int) {
	g.pushPC()
	g.Reg.PC = uint16(addr)
}

func scf(g *GBC, _, _ int) {
	g.setNHC(false, false, true)
}

// CCF Complement Carry Flag
func ccf(g *GBC, _, _ int) {
	g.setNHC(false, false, !g.f(flagC))
}
