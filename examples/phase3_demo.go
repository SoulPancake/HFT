package main

import (
	"fmt"
	"github.com/SoulPancake/HFT/core"
	"github.com/SoulPancake/HFT/types"
)

func createHierarchicalProcessor() {
	fmt.Println("=== Creating Hierarchical Processor Design ===")
	
	// Create ALU module
	alu := core.NewModule("ALU")
	alu.SetParameter("WIDTH", 32)
	
	a := alu.Input("a", 32)
	b := alu.Input("b", 32)
	op := alu.Input("op", 3)
	result := alu.Output("result", 32)
	zero := alu.Output("zero", 1)
	
	// Width-aware operations with automatic inference
	add_res := a.Add(b)
	sub_res := a.Sub(b)
	and_res := a.And(b)
	or_res := a.Or(b)
	shl_res := a.Shl(1)
	
	// Use mux with width inference
	mux_result := hdl.Mux(op, add_res, sub_res, and_res, or_res, shl_res)
	
	// Assign results
	alu.Assign(result, mux_result)
	zero_check := result.Eq(hdl.Lit(0, 32))
	alu.Assign(zero, zero_check)
	
	core.EmitVerilog(alu)
	fmt.Println("Generated ALU module with width inference")
	
	// Create Register File module
	regfile := core.NewModule("RegisterFile")
	regfile.SetParameter("ADDR_WIDTH", 5)
	regfile.SetParameter("DATA_WIDTH", 32)
	regfile.SetParameter("NUM_REGS", 32)
	
	clk := regfile.Input("clk", 1)
	_ = clk // Used in real implementation
	read_addr1 := regfile.Input("read_addr1", 5)
	read_addr2 := regfile.Input("read_addr2", 5)
	write_addr := regfile.Input("write_addr", 5)
	_ = write_addr // Used in real implementation
	write_data := regfile.Input("write_data", 32)
	_ = write_data // Used in real implementation
	write_enable := regfile.Input("write_enable", 1)
	_ = write_enable // Used in real implementation
	
	read_data1 := regfile.Output("read_data1", 32)
	read_data2 := regfile.Output("read_data2", 32)
	
	// Create register file using Vec
	registers := regfile.Vec("registers", 32, 32)
	
	// Read operations using multiplexers
	reg_signals1 := make([]*hdl.Signal, 32)
	reg_signals2 := make([]*hdl.Signal, 32)
	for i := 0; i < 32; i++ {
		reg_signals1[i] = registers.Get(i)
		reg_signals2[i] = registers.Get(i)
	}
	
	mux1_result := hdl.Mux(read_addr1, reg_signals1...)
	mux2_result := hdl.Mux(read_addr2, reg_signals2...)
	
	regfile.Assign(read_data1, mux1_result)
	regfile.Assign(read_data2, mux2_result)
	
	core.EmitVerilog(regfile)
	fmt.Println("Generated RegisterFile module with Vec")
	
	// Create top-level CPU module
	cpu := core.NewModule("SimpleCPU")
	cpu_clk := cpu.Input("clk", 1)
	cpu_rst := cpu.Input("rst", 1)
	_ = cpu_rst // Used in real implementation
	
	// Instruction and data interfaces using Bundles
	instr_bus := cpu.Bundle("instruction_bus")
	instr_bus.AddField("addr", &hdl.Signal{Width: 32, Kind: "wire"})
	instr_bus.AddField("data", &hdl.Signal{Width: 32, Kind: "wire"})
	instr_bus.AddField("valid", &hdl.Signal{Width: 1, Kind: "wire"})
	
	data_bus := cpu.Bundle("data_bus")
	data_bus.AddField("addr", &hdl.Signal{Width: 32, Kind: "wire"})
	data_bus.AddField("read_data", &hdl.Signal{Width: 32, Kind: "wire"})
	data_bus.AddField("write_data", &hdl.Signal{Width: 32, Kind: "wire"})
	data_bus.AddField("write_en", &hdl.Signal{Width: 1, Kind: "wire"})
	
	// Instantiate ALU
	alu_inst := cpu.Instance("ALU", "alu_inst")
	alu_inst.SetParameter("WIDTH", 32)
	
	alu_a := alu_inst.IO("a", 32)
	alu_b := alu_inst.IO("b", 32)
	alu_op := alu_inst.IO("op", 3)
	alu_result := alu_inst.IO("result", 32)
	alu_zero := alu_inst.IO("zero", 1)
	_ = alu_zero // Used in real implementation
	
	// Instantiate Register File
	rf_inst := cpu.Instance("RegisterFile", "regfile_inst")
	rf_inst.SetParameter("ADDR_WIDTH", 5)
	rf_inst.SetParameter("DATA_WIDTH", 32)
	rf_inst.SetParameter("NUM_REGS", 32)
	
	rf_inst.Connect("clk", cpu_clk)
	rf_read_addr1 := rf_inst.IO("read_addr1", 5)
	rf_read_addr2 := rf_inst.IO("read_addr2", 5)
	rf_write_addr := rf_inst.IO("write_addr", 5)
	rf_write_data := rf_inst.IO("write_data", 32)
	rf_write_enable := rf_inst.IO("write_enable", 1)
	rf_read_data1 := rf_inst.IO("read_data1", 32)
	rf_read_data2 := rf_inst.IO("read_data2", 32)
	
	// Connect ALU inputs to register file outputs
	cpu.Assign(alu_a, rf_read_data1)
	cpu.Assign(alu_b, rf_read_data2)
	
	// Simple control logic (decode instruction fields)
	instr_data := instr_bus.GetField("data")
	opcode := instr_data.Bits(31, 26)
	rs1 := instr_data.Bits(25, 21)
	rs2 := instr_data.Bits(20, 16)
	rd := instr_data.Bits(15, 11)
	
	// Convert opcode to ALU operation (simplified)
	alu_op_wire := cpu.Wire("alu_op_decode", 3)
	cpu.Assign(alu_op_wire, opcode.Bits(2, 0))
	cpu.Assign(alu_op, alu_op_wire)
	
	// Connect register addresses
	cpu.Assign(rf_read_addr1, rs1)
	cpu.Assign(rf_read_addr2, rs2)
	cpu.Assign(rf_write_addr, rd)
	cpu.Assign(rf_write_data, alu_result)
	
	// Simple write enable (always write for now)
	cpu.Assign(rf_write_enable, hdl.Bool(true))
	
	core.EmitVerilog(cpu)
	fmt.Println("Generated SimpleCPU module with hierarchical design")
}

func demonstrateWidthInference() {
	fmt.Println("\n=== Width Inference Demonstration ===")
	
	m := core.NewModule("WidthDemo")
	
	// Different width signals
	a8 := m.Input("a8", 8)
	b16 := m.Input("b16", 16)
	c4 := m.Input("c4", 4)
	
	// Arithmetic operations with automatic width inference
	add_result := a8.Add(b16)  // Result width: max(8, 16) = 16
	fmt.Printf("a8 + b16 -> width: %d\n", add_result.Width)
	
	mul_result := a8.Mul(c4)   // Result width: 8 + 4 = 12
	fmt.Printf("a8 * c4 -> width: %d\n", mul_result.Width)
	
	// Shift operations
	shl_result := a8.Shl(3)    // Result width: 8 + 3 = 11
	fmt.Printf("a8 << 3 -> width: %d\n", shl_result.Width)
	
	shr_result := b16.Shr(4)   // Result width: 16 - 4 = 12
	fmt.Printf("b16 >> 4 -> width: %d\n", shr_result.Width)
	
	// Comparison operations (always 1-bit result)
	eq_result := a8.Eq(b16)    // Result width: 1
	fmt.Printf("a8 == b16 -> width: %d\n", eq_result.Width)
	
	// Concatenation
	cat_result := hdl.Cat(a8, c4, b16) // Result width: 8 + 4 + 16 = 28
	fmt.Printf("Cat(a8, c4, b16) -> width: %d\n", cat_result.Width)
	
	// Bit selection
	bits_result := b16.Bits(11, 4)     // Result width: 11 - 4 + 1 = 8
	fmt.Printf("b16[11:4] -> width: %d\n", bits_result.Width)
	
	// Output the results
	out1 := m.Output("add_out", add_result.Width)
	out2 := m.Output("mul_out", mul_result.Width)
	out3 := m.Output("eq_out", eq_result.Width)
	out4 := m.Output("cat_out", cat_result.Width)
	
	m.Assign(out1, add_result)
	m.Assign(out2, mul_result)
	m.Assign(out3, eq_result)
	m.Assign(out4, cat_result)
	
	core.EmitVerilog(m)
	fmt.Println("Generated WidthDemo module demonstrating automatic width inference")
}

func demonstrateTypeSafety() {
	fmt.Println("\n=== Type Safety Demonstration ===")
	
	sig1 := &hdl.Signal{Name: "sig1", Width: 8, Kind: "wire"}
	sig2 := &hdl.Signal{Name: "sig2", Width: 16, Kind: "wire"}
	
	// Demonstrate width checking
	fmt.Printf("sig1 width: %d, sig2 width: %d\n", sig1.Width, sig2.Width)
	
	if err := sig1.CheckWidth(sig2); err != nil {
		fmt.Printf("Width check failed as expected: %v\n", err)
	}
	
	// Demonstrate literal creation with type safety
	lit8 := hdl.Lit(255, 8)
	lit16 := hdl.Lit(65535, 16)
	
	fmt.Printf("8-bit literal: %s (width: %d)\n", lit8.Name, lit8.Width)
	fmt.Printf("16-bit literal: %s (width: %d)\n", lit16.Name, lit16.Width)
	
	// Demonstrate boolean literals
	true_sig := hdl.Bool(true)
	false_sig := hdl.Bool(false)
	
	fmt.Printf("Boolean true: %s (width: %d)\n", true_sig.Name, true_sig.Width)
	fmt.Printf("Boolean false: %s (width: %d)\n", false_sig.Name, false_sig.Width)
}

func runPhase3Demo() {
	createHierarchicalProcessor()
	demonstrateWidthInference()
	demonstrateTypeSafety()
	
	fmt.Println("\n=== Summary ===")
	fmt.Println("✅ Hierarchical module design with instances")
	fmt.Println("✅ Automatic width inference for all operations")
	fmt.Println("✅ Type safety with width checking")
	fmt.Println("✅ Parameterized modules")
	fmt.Println("✅ Bundle and Vec data structures")
	fmt.Println("✅ Module instantiation with parameters and connections")
	fmt.Println("✅ Comprehensive Verilog emission")
}