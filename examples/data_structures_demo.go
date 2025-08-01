package main

import (
	"fmt"
	"github.com/SoulPancake/HFT/core"
	"github.com/SoulPancake/HFT/types"
)

func createBundleExample() *hdl.Module {
	m := core.NewModule("BundleExample")
	
	// Inputs (for clock and reset)
	_ = m.Input("clk", 1)
	_ = m.Input("rst", 1)
	
	// Create input and output bundles
	input_bundle := m.Bundle("input_bus")
	output_bundle := m.Bundle("output_bus")
	
	// Add fields to bundles
	input_bundle.AddField("data", &hdl.Signal{Width: 16, Kind: "wire"})
	input_bundle.AddField("valid", &hdl.Signal{Width: 1, Kind: "wire"})
	input_bundle.AddField("ready", &hdl.Signal{Width: 1, Kind: "wire"})
	
	output_bundle.AddField("data", &hdl.Signal{Width: 16, Kind: "wire"})
	output_bundle.AddField("valid", &hdl.Signal{Width: 1, Kind: "wire"})
	output_bundle.AddField("ready", &hdl.Signal{Width: 1, Kind: "wire"})
	
	// Connect bundles (pass-through example)
	connections := input_bundle.Connect(output_bundle)
	for _, conn := range connections {
		m.Assigns = append(m.Assigns, conn)
	}
	
	return m
}

func createVecExample() *hdl.Module {
	m := core.NewModule("VecExample")
	
	// Inputs
	_ = m.Input("clk", 1)
	sel := m.Input("sel", 2)
	
	// Create vector of registers
	reg_vec := m.Vec("registers", 4, 8)
	
	// Create output
	selected_data := m.Output("selected_data", 8)
	
	// Use multiplexer to select from vector
	vec_signals := make([]*hdl.Signal, 4)
	for i := 0; i < 4; i++ {
		vec_signals[i] = reg_vec.Get(i)
	}
	
	mux_result := hdl.Mux(sel, vec_signals...)
	m.Assign(selected_data, mux_result)
	
	return m
}

func createMemoryExample() *hdl.Module {
	m := core.NewModule("MemoryExample")
	
	// Inputs
	_ = m.Input("clk", 1)
	_ = m.Input("rst", 1)
	addr := m.Input("addr", 8)
	write_data := m.Input("write_data", 16)
	write_enable := m.Input("write_enable", 1)
	read_enable := m.Input("read_enable", 1)
	
	// Outputs
	read_data := m.Output("read_data", 16)
	
	// Create memory
	data_mem := m.SyncMem("data_memory", 16, 256)
	
	// Read operation
	mem_read := data_mem.Read(addr)
	read_mux := hdl.Mux(read_enable, 
		hdl.Lit(0, 16), 
		mem_read)
	m.Assign(read_data, read_mux)
	
	// Write operation would be handled in always block
	// For now, demonstrate the write statement generation
	write_stmt := data_mem.Write(addr, write_data, write_enable)
	fmt.Printf("Generated write statement: %s\n", write_stmt)
	
	return m
}

func createComplexExample() *hdl.Module {
	m := core.NewModule("ComplexExample")
	
	// Inputs
	_ = m.Input("clk", 1)
	_ = m.Input("rst", 1)
	
	// Create a processor-like interface using bundles
	cpu_bundle := m.Bundle("cpu_interface")
	cpu_bundle.AddField("addr", &hdl.Signal{Width: 16, Kind: "wire"})
	cpu_bundle.AddField("data_in", &hdl.Signal{Width: 32, Kind: "wire"})
	cpu_bundle.AddField("data_out", &hdl.Signal{Width: 32, Kind: "wire"})
	cpu_bundle.AddField("write_en", &hdl.Signal{Width: 1, Kind: "wire"})
	cpu_bundle.AddField("read_en", &hdl.Signal{Width: 1, Kind: "wire"})
	
	// Create register file using Vec
	register_file := m.Vec("register_file", 16, 32)
	
	// Create cache memory (demonstrate memory creation)
	_ = m.SyncMem("cache", 32, 1024)
	
	// Decode address to select register
	reg_select := m.Wire("reg_select", 4)
	addr_field := cpu_bundle.GetField("addr")
	addr_bits := addr_field.Bits(3, 0)
	m.Assign(reg_select, addr_bits)
	
	// Select from register file
	reg_signals := make([]*hdl.Signal, 16)
	for i := 0; i < 16; i++ {
		reg_signals[i] = register_file.Get(i)
	}
	
	selected_reg := hdl.Mux(reg_select, reg_signals...)
	
	// Output the selected register
	data_out_field := cpu_bundle.GetField("data_out")
	m.Assign(data_out_field, selected_reg)
	
	return m
}

func runDataStructuresDemo() {
	fmt.Println("=== Bundle Example ===")
	bundle_mod := createBundleExample()
	core.EmitVerilog(bundle_mod)
	fmt.Printf("Generated Bundle Example module\n\n")
	
	fmt.Println("=== Vec Example ===")
	vec_mod := createVecExample()
	core.EmitVerilog(vec_mod)
	fmt.Printf("Generated Vec Example module\n\n")
	
	fmt.Println("=== Memory Example ===")
	mem_mod := createMemoryExample()
	core.EmitVerilog(mem_mod)
	fmt.Printf("Generated Memory Example module\n\n")
	
	fmt.Println("=== Complex Example ===")
	complex_mod := createComplexExample()
	core.EmitVerilog(complex_mod)
	fmt.Printf("Generated Complex Example module\n\n")
}