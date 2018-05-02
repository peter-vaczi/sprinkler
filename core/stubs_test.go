package core_test

import "github.com/peter-vaczi/sprinkler/gpio"

type GpioStub struct {
	pins map[int]*PinStub
}

func NewGpioStub() *GpioStub {
	return &GpioStub{pins: make(map[int]*PinStub)}
}

func (g *GpioStub) NewPin(p int) gpio.Pin {
	pin := &PinStub{pin: p}
	g.pins[p] = pin
	return pin
}

type PinStub struct {
	pin    int
	output bool
	high   bool
}

func (p *PinStub) Output() { p.output = true }
func (p *PinStub) High()   { p.high = true }
func (p *PinStub) Low()    { p.high = false }
