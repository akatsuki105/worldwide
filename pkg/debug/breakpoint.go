package debug

import (
	"fmt"
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

func (b *Break) On() bool {
	return b.flag == BreakOn
}

func (b *Break) Off() bool {
	return b.flag == BreakOff
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
	Cond Cond
}

type Cond struct {
	On      bool
	LHS     string
	Operand string
	RHS     uint16
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

	bank, PC := parseBankPC(slice[0])
	if bank == 0 && PC == 0 {
		ok = false
		return bk, ok
	}
	bk.Bank = bank
	bk.PC = PC

	bk.Cond = parseCond(slice[1])
	fmt.Println(bk.Cond)
	return bk, true
}

func parseBankPC(s string) (bank byte, PC uint16) {
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

func parseCond(s string) (cond Cond) {
	cond = Cond{
		On: false,
	}

	const (
		equal  = "=="
		nEqual = "!="
		gte    = ">="
		lte    = "<="
		gt     = ">"
		lt     = "<"
	)

	var slice []string
	switch {
	case strings.Index(s, equal) >= 0:
		slice = strings.Split(s, equal)
		cond.Operand = equal
	case strings.Index(s, nEqual) >= 0:
		slice = strings.Split(s, nEqual)
		cond.Operand = nEqual
	case strings.Index(s, gte) >= 0:
		slice = strings.Split(s, gte)
		cond.Operand = gte
	case strings.Index(s, lte) >= 0:
		slice = strings.Split(s, lte)
		cond.Operand = lte
	case strings.Index(s, gt) >= 0:
		slice = strings.Split(s, gt)
		cond.Operand = gt
	case strings.Index(s, lt) >= 0:
		slice = strings.Split(s, lt)
		cond.Operand = lt
	}

	if len(slice) != 2 {
		return cond
	}

	bitsize := 8
	switch slice[0] {
	case "A", "F", "B", "C", "D", "E", "H", "L":
		cond.LHS = slice[0]
	case "AF", "BC", "DE", "HL", "SP":
		cond.LHS = slice[0]
		bitsize = 16
	default:
		return cond
	}

	rhs, err := strconv.ParseUint(slice[1], 16, bitsize)
	if err != nil {
		return cond
	}
	cond.RHS = uint16(rhs)

	cond.On = true
	return cond
}
