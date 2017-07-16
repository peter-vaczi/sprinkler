package core

type Program struct {
	Name string `json:"name"`
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
