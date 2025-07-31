// builder.go
package core

import (
	"github.com/SoulPancake/HFT/types"
)

func NewModule(name string) *hdl.Module {
	return &hdl.Module{
		Name:    name,
		Inputs:  []*hdl.Signal{},
		Outputs: []*hdl.Signal{},
		Wires:   []*hdl.Signal{},
		Assigns: []string{},
	}
}
