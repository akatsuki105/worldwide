package debug

type CPU struct {
	all  int
	halt int
}

type Monitor struct {
	CPU
}

func (c *CPU) Usage() float64 {
	all, halt := float64(c.all), float64(c.halt)
	return (all - halt) / all
}

func (c *CPU) Increment(halt bool) {
	c.all++
	if halt {
		c.halt++
	}
}
