package gpio

import rpio "github.com/stianeikeland/go-rpio"

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
