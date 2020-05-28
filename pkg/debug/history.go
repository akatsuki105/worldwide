package debug

import (
	"fmt"
	"gbc/pkg/util"
)

// History - CPU instruction log
type History struct {
	ptr    uint
	buffer [10]string
}

func (h *History) SetHistory(bank byte, PC uint16, opcode byte) {
	if PC <= 0x4000 {
		bank = 0
	}
	bankPC := fmt.Sprintf("%02x:%04x: ", bank, PC)

	instruction, operand1, operand2 := util.OpcodeToString[opcode][0], util.OpcodeToString[opcode][1], util.OpcodeToString[opcode][2]
	switch {
	case operand1 == "*" && operand2 == "*":
		h.buffer[h.ptr] = bankPC + instruction
	case operand2 == "*":
		h.buffer[h.ptr] = bankPC + instruction + " " + operand1
	default:
		h.buffer[h.ptr] = bankPC + instruction + " " + operand1 + ", " + operand2
	}
	h.ptr = (h.ptr + 1) % 10
}

func (h *History) History() string {
	result := "History\n"
	for i := 0; i < 10; i++ {
		index := (h.ptr + uint(i)) % 10
		log := h.buffer[index]
		if i < 9 {
			result += fmt.Sprintf("%d:    %0s\n", -(9 - i), log)
		} else if i == 9 {
			result += fmt.Sprintf(" %d:    %0s\n", 9-i, log)
		}
	}
	return result
}
