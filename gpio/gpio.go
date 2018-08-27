package gpio

import (
	rpio "github.com/stianeikeland/go-rpio"
)

//https://periph.io

type Gpio interface {
	NewPin(p int) Pin
}

type gpioImpl struct {
}

func (g *gpioImpl) NewPin(p int) Pin {
	return pin(p)
}

func New() (Gpio, error) {
	err := rpio.Open()
	if err != nil {
		return nil, err
	}
	return &gpioImpl{}, nil
}

type Pin interface {
	Output()
	High()
	Low()
}

type pin rpio.Pin

func (p pin) Output() {
	rpio.Pin(p).Output()
}

func (p pin) High() {
	rpio.Pin(p).High()
}

func (p pin) Low() {
	rpio.Pin(p).Low()
}

type dummyGpioImpl struct {
	pins map[int]*dummyPin
}

func NewDummy() Gpio {
	return &dummyGpioImpl{pins: make(map[int]*dummyPin)}
}

func (g *dummyGpioImpl) NewPin(p int) Pin {
	pin := &dummyPin{pin: p}
	g.pins[p] = pin
	return pin
}

type dummyPin struct {
	pin    int
	output bool
	high   bool
}

func (p *dummyPin) Output() { p.output = true }
func (p *dummyPin) High()   { p.high = true }
func (p *dummyPin) Low()    { p.high = false }
