package debug

import (
	"strconv"
	"strings"
)

const (
	BreakOff = iota
	BreakOn
	BreakDelay
)

// Break - Breakpoint state info used in debug mode.
type Break struct {
	flag        int
	breakpoints []BreakPoint
}

func (b *Break) Flag() int {
	return b.flag
}

func (b *Break) SetFlag(flag int) {
	b.flag = flag
}

func (b *Break) BreakPoints() []BreakPoint {
	return b.breakpoints
}

// BreakPoint - A Breakpoint info used in debug mode
type BreakPoint struct {
	Bank byte
	PC   uint16
	Cond string
}

func (b *Break) ParseBreakpoints(breakpoints []string) {
	for _, s := range breakpoints {
		if bk, ok := newBreakPoint(s); ok {
			b.breakpoints = append(b.breakpoints, bk)
		}
	}
}

func newBreakPoint(s string) (bk BreakPoint, ok bool) {
	slice := strings.Split(s, ";") // [00:0460], [SP==c0f3]
	if len(slice) < 2 {
		ok = false
		return bk, ok
	}

	bank, PC := parseBreakpointsPC(slice[0])
	if bank == 0 && PC == 0 {
		ok = false
		return bk, ok
	}
	bk.Bank = bank
	bk.PC = PC

	return bk, true
}

func parseBreakpointsPC(s string) (bank byte, PC uint16) {
	bankPC := strings.Split(s, ":")
	if len(bankPC) < 2 {
		return 0, 0
	}

	var err error
	bankI64, err := strconv.ParseInt(bankPC[0], 16, 8)
	if err != nil {
		return 0, 0
	}
	pcI64, err := strconv.ParseInt(bankPC[1], 16, 16)
	if err != nil {
		return 0, 0
	}
	bank = byte(bankI64)
	PC = uint16(pcI64)
	return bank, PC
}
