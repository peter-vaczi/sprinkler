package core

import (
	"errors"
	"log"

	"github.com/peter.vaczi/sprinklerd/gpio"
)

var (
	AlreadyExists = errors.New("Already exists")
	NotFound      = errors.New("Not found")
	gpioLib       gpio.Gpio
)

func InitGpio(g gpio.Gpio) {
	gpioLib = g
}

type Device struct {
	Name string `json:"name"`
	On   bool   `json:"on"`
	Pin  int    `json:"pin"`
	pin  gpio.Pin
}

type Devices map[string]*Device

func NewDevices() *Devices {
	devs := make(Devices)
	return &devs
}

func (d *Devices) Add(dev *Device) error {
	if _, exists := (*d)[dev.Name]; exists {
		return AlreadyExists
	}

	(*d)[dev.Name] = dev
	dev.SetState(dev.Pin, dev.On)

	return nil
}

func (d *Devices) Get(name string) (*Device, error) {
	if dev, exists := (*d)[name]; exists {
		return dev, nil
	}

	return nil, NotFound
}

func (d *Devices) Del(name string) error {
	if _, exists := (*d)[name]; exists {
		delete(*d, name)
		return nil
	}

	return NotFound
}

func (d *Devices) Set(name string, newDev *Device) error {
	if dev, exists := (*d)[name]; exists {
		dev.SetState(newDev.Pin, newDev.On)
		return nil
	}

	return NotFound
}

func (d *Device) SetPin(pin int) {
	d.Pin = pin
	d.pin = gpioLib.NewPin(pin)
	d.pin.Output()
}

func (d *Device) TurnOn() {
	d.On = true
	d.pin.High()
	log.Printf("device %s is on", d.Name)
}

func (d *Device) TurnOff() {
	d.On = false
	d.pin.Low()
	log.Printf("device %s is off", d.Name)
}

func (d *Device) SetState(pin int, on bool) {
	d.SetPin(pin)
	if on {
		d.TurnOn()
	} else {
		d.TurnOff()
	}
}

func (d *Device) Init() {
	d.SetState(d.Pin, d.On)
}

func (d *Device) IsOn() bool {
	return d.On
}
