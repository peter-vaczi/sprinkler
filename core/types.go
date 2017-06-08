package core

import "errors"

var (
	AlreadyExists = errors.New("Already exists")
	NotFound      = errors.New("Not found")
)

type Device struct {
	Name string `json:"name"`
	On   bool   `json:"on"`
}

type Devices map[string]*Device

func New() *Devices {
	devs := make(Devices)
	return &devs
}

func (d *Devices) Add(dev *Device) error {
	if _, exists := (*d)[dev.Name]; exists {
		return AlreadyExists
	}

	(*d)[dev.Name] = dev
	return nil
}

func (d *Devices) Del(name string) error {
	if _, exists := (*d)[name]; exists {
		delete(*d, name)
		return nil
	}

	return NotFound
}

func (d *Device) TurnOn() {
	d.On = true
}

func (d *Device) TurnOff() {
	d.On = false
}

func (d *Device) IsOn() bool {
	return d.On
}
