package gpio

import rpio "github.com/stianeikeland/go-rpio"

type Gpio struct {
}

func New() (*Gpio, error) {
	err := rpio.Open()
	if err != nil {
		return nil, err
	}
	return &Gpio{}, nil
}

type Pin interface {
	Output()
	High()
	Low()
}

type pin struct {
	pin rpio.Pin
}

func (g *Gpio) NewPin(p int) Pin {
	return &pin{pin: rpio.Pin(p)}
}

func (p *pin) Output() {
	p.pin.Output()
}

func (p *pin) High() {
	p.pin.High()
}

func (p *pin) Low() {
	p.pin.Low()
}
