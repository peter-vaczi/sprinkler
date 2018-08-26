package core

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"
)

type ProgramElement struct {
	DeviceName string        `json:"device"`
	Device     *Device       `json:"-"`
	Duration   time.Duration `json:"duration"`
}

type Program struct {
	Name     string            `json:"name"`
	Elements []*ProgramElement `json:"devices"`
	ctx      context.Context
	cancel   context.CancelFunc
	running  bool
	m        sync.Mutex
}

type Programs map[string]*Program

var (
	OutOfRange = errors.New("Element index out of range")
)

func NewPrograms() *Programs {
	p := make(Programs)
	return &p
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

func (p *Programs) IsDeviceInUse(name string) bool {
	for _, pr := range *p {
		for _, e := range pr.Elements {
			if e.DeviceName == name {
				return true
			}
		}
	}
	return false
}

func (p *Program) AddDevice(device *Device, duration time.Duration) error {
	p.m.Lock()
	defer p.m.Unlock()

	p.Elements = append(p.Elements, &ProgramElement{DeviceName: device.Name, Device: device, Duration: duration})

	return nil
}

func (p *Program) DelDevice(idx int) error {
	p.m.Lock()
	defer p.m.Unlock()

	if idx >= len(p.Elements) {
		return OutOfRange
	}
	p.Elements = append(p.Elements[:idx], p.Elements[idx+1:]...)
	return nil
}

func (p *Program) Start() {
	p.m.Lock()
	defer p.m.Unlock()

	if !p.running {
		p.ctx, p.cancel = context.WithCancel(context.Background())
		go p.run()
	}
}

func (p *Program) Stop() {
	p.m.Lock()
	running := p.running
	p.m.Unlock()

	if running {
		p.cancel()
		p.m.Lock()
		for _, elem := range p.Elements {
			elem.Device.TurnOff()
		}
		p.m.Unlock()
	}
}

func (p *Program) run() {
	p.m.Lock()
	p.running = true
	elements := p.Elements
	p.m.Unlock()

	defer func() {
		p.m.Lock()
		p.running = false
		p.m.Unlock()
	}()

	log.Printf("program %s is started", p.Name)
	for _, elem := range elements {
		if elem.Device.IsOn() {
			elem.Device.TurnOff()
		}
	}

	for _, elem := range elements {
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
}
