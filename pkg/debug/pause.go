package debug

type Pause struct {
	on    bool
	delay int
}

func (p *Pause) On() bool {
	return p.on
}

func (p *Pause) SetOn(delay int) {
	p.on = true
	p.delay = delay
}

func (p *Pause) SetOff(delay int) {
	p.on = false
	p.delay = delay
}

func (p *Pause) Delay() bool {
	return p.delay > 0
}

func (p *Pause) DecrementDelay() {
	p.delay--
}
