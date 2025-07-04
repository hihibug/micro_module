package Framework

import (
	"sync"
)

type service struct {
	opts   Options
	extend map[string]OptionHandle
	once   sync.Once
	err    error
}

// Run implements Service.
func (*service) Run() error {
	panic("unimplemented")
}

func newService(opts ...Option) *service {
	return &service{
		opts:   newOptions(opts...),
		extend: make(map[string]OptionHandle),
	}
}
