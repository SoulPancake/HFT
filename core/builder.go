package core

import (
	"github.com/SoulPancake/HFT/types"
)

func NewModule(name string) *hdl.Module {
	return &hdl.Module{
		Name:         name,
		Inputs:       []*hdl.Signal{},
		Outputs:      []*hdl.Signal{},
		Wires:        []*hdl.Signal{},
		Regs:         []*hdl.Signal{},
		Bundles:      []*hdl.Bundle{},
		Vecs:         []*hdl.Vec{},
		Memories:     []*hdl.Memory{},
		Instances:    []*hdl.ModuleInstance{},
		Assigns:      []string{},
		Always:       []string{},
		Parameters:   make(map[string]interface{}),
		ClockDomains: []*hdl.ClockDomain{},
		Mutexes:      []*hdl.Mutex{},
		Templates:    []*hdl.ModuleTemplate{},
	}
}
