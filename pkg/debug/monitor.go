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
	return (all - halt) * 100 / all
}

func (c *CPU) Add(halt bool, count int) {
	c.all += count
	if halt {
		c.halt += count
	}
}

func (c *CPU) Reset() {
	c.all, c.halt = 1, 1
}
