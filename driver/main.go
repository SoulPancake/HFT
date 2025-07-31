package main

import (
	"github.com/SoulPancake/HFT/core"
)

func main() {
	m := core.NewModule("Adder")
	a := m.Input("a", 8)
	b := m.Input("b", 8)
	sum := m.Output("sum", 8)

	m.Assign(sum.Name, a.Add(b))
	core.EmitVerilog(m)
}
