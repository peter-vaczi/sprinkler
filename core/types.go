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

func (d *Devices) Add(name string) (*Device, error) {
	if _, exists := (*d)[name]; exists {
		return nil, AlreadyExists
	}

	newDev := &Device{Name: name}
	(*d)[name] = newDev
	return newDev, nil
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
