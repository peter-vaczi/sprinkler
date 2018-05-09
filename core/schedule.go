package core

import (
	"time"

	"github.com/robfig/cron"
)

type Schedule struct {
	Name        string        `json:"name"`
	ProgramName string        `json:"program"`
	Program     *Program      `json:"-"`
	Spec        string        `json:"spec"`
	Sched       cron.Schedule `json:"-"`
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
		sch.SetSpec(newSch.Spec)
		return nil
	}

	return NotFound
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
