package core

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/robfig/cron"
)

type Schedule struct {
	Name        string        `json:"name"`
	ProgramName string        `json:"program"`
	Program     *Program      `json:"-"`
	Spec        string        `json:"spec"`
	Sched       cron.Schedule `json:"-"`
	Enabled     bool          `json:"enabled"`
	m           sync.Mutex
	ctx         context.Context
	cancel      context.CancelFunc
}

type Schedules map[string]*Schedule

func NewSchedules() *Schedules {
	s := make(Schedules)
	return &s
}

func (s *Schedules) Add(sched *Schedule) error {
	if _, exists := (*s)[sched.Name]; exists {
		return AlreadyExists
	}

	sched.SetProgram(sched.Program)
	err := sched.SetSpec(sched.Spec)
	if err != nil {
		return err
	}

	(*s)[sched.Name] = sched

	return nil
}

func (s *Schedules) Get(name string) (*Schedule, error) {
	if sched, exists := (*s)[name]; exists {
		return sched, nil
	}

	return nil, NotFound
}

func (s *Schedules) Del(name string) error {
	if _, exists := (*s)[name]; exists {
		delete(*s, name)
		return nil
	}

	return NotFound
}

func (s *Schedules) Set(name string, newSch *Schedule) error {
	if sch, exists := (*s)[name]; exists {
		sch.SetProgram(newSch.Program)
		err := sch.SetSpec(newSch.Spec)
		if err != nil {
			return err
		}
		if newSch.Enabled {
			sch.Enable()
		} else {
			sch.Disable()
		}
		return nil
	}

	return NotFound
}

func (s *Schedules) DisableAll() {
	for _, sc := range *s {
		sc.Disable()
	}
}

func (s *Schedule) SetProgram(prog *Program) {
	s.Program = prog
	if prog != nil {
		s.ProgramName = prog.Name
	}
}

func (s *Schedule) SetSpec(spec string) error {
	sc, err := cron.ParseStandard(spec)
	if err != nil {
		return err
	}
	s.Spec = spec
	s.Sched = sc
	return nil
}

func (s *Schedule) GetNext() time.Time {
	return s.Sched.Next(time.Now())
}

func (s *Schedule) Enable() {
	s.m.Lock()
	s.Enabled = true
	s.m.Unlock()

	s.kill()
	go s.run()
}

func (s *Schedule) Disable() {
	s.m.Lock()
	s.Enabled = false
	s.m.Unlock()

	s.kill()
}

func (s *Schedule) kill() {
	s.m.Lock()
	if s.cancel != nil {
		s.cancel()
		s.cancel = nil
	}
	s.m.Unlock()
}

func (s *Schedule) run() {
	s.m.Lock()
	next := s.GetNext()
	prog := s.Program
	s.ctx, s.cancel = context.WithCancel(context.Background())
	s.m.Unlock()

	log.Printf("schedule %s will start program %s at %s", s.Name, s.Program.Name, next)
	t := time.NewTimer(next.Sub(time.Now()))
	select {
	case <-s.ctx.Done():
		log.Printf("schedule %s is canceled", s.Name)
		if !t.Stop() {
			<-t.C
		}
		return
	case <-t.C:
		log.Printf("schedule %s is starting program %s now", s.Name, s.Program.Name)
		prog.Start()
		go s.run()
	}
}
