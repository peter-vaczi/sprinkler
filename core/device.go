package core

import (
	"errors"
	"log"

	rpio "github.com/stianeikeland/go-rpio"
)

var (
	AlreadyExists = errors.New("Already exists")
	NotFound      = errors.New("Not found")
)

func init() {
	err := rpio.Open()
	if err != nil {
		log.Fatal(err)
	}
}

type Device struct {
	Name string `json:"name"`
	On   bool   `json:"on"`
	Pin  int    `json:"pin"`
	pin  rpio.Pin
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
	dev.pin = rpio.Pin(dev.Pin)
	dev.pin.Output()
	if dev.On {
		dev.pin.High()
	} else {
		dev.pin.Low()
	}

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
	d.pin.High()
}

func (d *Device) TurnOff() {
	d.On = false
	d.pin.Low()
}

func (d *Device) IsOn() bool {
	return d.On
}
