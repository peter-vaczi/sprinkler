package core

import (
	"context"
	"errors"
	"log"
	"time"
)

type ProgramElement struct {
	DeviceName string        `json:"device"`
	Device     *Device       `json:"-"`
	Duration   time.Duration `json:"duration"`
}

type Program struct {
	Name     string            `json:"name"`
	Elements []*ProgramElement `json:"elements"`
	ctx      context.Context
	cancel   context.CancelFunc
}

type Programs map[string]*Program

func NewPrograms() *Programs {
	progs := make(Programs)
	return &progs
}

func (p *Programs) Add(prog *Program) error {
	if _, exists := (*p)[prog.Name]; exists {
		return AlreadyExists
	}

	(*p)[prog.Name] = prog

	return nil
}

func (p *Programs) Get(name string) (*Program, error) {
	if prg, exists := (*p)[name]; exists {
		return prg, nil
	}

	return nil, NotFound
}

func (p *Programs) Del(name string) error {
	if _, exists := (*p)[name]; exists {
		delete(*p, name)
		return nil
	}

	return NotFound
}

func (p *Program) AddDevice(device *Device, duration time.Duration) error {
	p.Elements = append(p.Elements, &ProgramElement{DeviceName: device.Name, Device: device, Duration: duration})

	return nil
}

func (p *Program) DelDevice(idx int) error {
	if idx >= len(p.Elements) {
		return errors.New("Element index out of range")
	}
	p.Elements = append(p.Elements[:idx], p.Elements[idx+1:]...)
	return nil
}

func (p *Program) Reset() error {
	p.cancel()
	for _, elem := range p.Elements {
		elem.Device.TurnOff()
	}

	return nil
}

func (p *Program) Start() error {
	p.ctx, p.cancel = context.WithCancel(context.Background())
	go func() {
		log.Printf("program %s is started", p.Name)

		for _, elem := range p.Elements {
			elem.Device.TurnOn()

			t := time.NewTimer(elem.Duration)
			select {
			case <-p.ctx.Done():
				log.Printf("program %s is canceled", p.Name)
				return
			case <-t.C:
				// do nothing
			}
			elem.Device.TurnOff()
			time.Sleep(1 * time.Second)
		}
		log.Printf("program %s is finished", p.Name)
	}()

	return nil
}

func (p *Program) Stop() error {
	return p.Reset()
}
