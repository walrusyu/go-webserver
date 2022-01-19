package main

type People struct {
	ID      int
	Name    string
	Address string
}

type Option func(*People)

func WithAddress(address string) Option {
	return func(p *People) {
		p.Address = address
	}
}

func NewUser(id int, name string, opts ...Option) *People {
	p := &People{
		ID:   id,
		Name: name,
	}
	for _, o := range opts {
		o(p)
	}
	return p
}
