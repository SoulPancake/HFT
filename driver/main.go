package main

import (
	"github.com/SoulPancake/HFT/core"
)

func main() {
	// Create a simple adder with new API
	m := core.NewModule("SimpleAdder")
	
	// Inputs
	a := m.Input("a", 8)
	b := m.Input("b", 8)
	
	// Outputs
	sum := m.Output("sum", 8)
	
	// Use the new width-aware addition
	add_result := a.Add(b)
	
	// Assign result
	m.Assign(sum, add_result)
	
	core.EmitVerilog(m)
}
