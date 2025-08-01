package core

import (
	"os"
	"strings"
	"testing"
	
	hdl "github.com/SoulPancake/HFT/types"
)

func TestNewModule(t *testing.T) {
	name := "TestModule"
	m := NewModule(name)

	if m.Name != name {
		t.Errorf("NewModule name = %v, want %v", m.Name, name)
	}
	if m.Inputs == nil || len(m.Inputs) != 0 {
		t.Errorf("NewModule inputs should be empty slice, got %v", m.Inputs)
	}
	if m.Outputs == nil || len(m.Outputs) != 0 {
		t.Errorf("NewModule outputs should be empty slice, got %v", m.Outputs)
	}
	if m.Wires == nil || len(m.Wires) != 0 {
		t.Errorf("NewModule wires should be empty slice, got %v", m.Wires)
	}
	if m.Regs == nil || len(m.Regs) != 0 {
		t.Errorf("NewModule regs should be empty slice, got %v", m.Regs)
	}
	if m.Assigns == nil || len(m.Assigns) != 0 {
		t.Errorf("NewModule assigns should be empty slice, got %v", m.Assigns)
	}
	if m.Always == nil || len(m.Always) != 0 {
		t.Errorf("NewModule always should be empty slice, got %v", m.Always)
	}
}

func TestEmitVerilogBasic(t *testing.T) {
	// Create a simple adder module
	m := NewModule("Adder")
	a := m.Input("a", 8)
	b := m.Input("b", 8)
	sum := m.Output("sum", 8)

	m.Assign(sum, a.Add(b))

	// Emit Verilog to temporary file
	originalFile := "out.v"
	testFile := "test_out.v"
	
	// Backup original if it exists
	if _, err := os.Stat(originalFile); err == nil {
		content, _ := os.ReadFile(originalFile)
		defer os.WriteFile(originalFile, content, 0644)
	}
	
	// Rename output file for testing
	defer func() {
		if _, err := os.Stat(originalFile); err == nil {
			os.Rename(originalFile, testFile)
		}
	}()

	EmitVerilog(m)

	// Read and verify output
	content, err := os.ReadFile(originalFile)
	if err != nil {
		t.Fatalf("Failed to read generated Verilog file: %v", err)
	}

	verilog := string(content)
	
	// Check module declaration
	if !strings.Contains(verilog, "module Adder(") {
		t.Errorf("Module declaration not found in generated Verilog")
	}
	
	// Check inputs and outputs
	if !strings.Contains(verilog, "input [7:0] a") {
		t.Errorf("Input 'a' declaration not found")
	}
	if !strings.Contains(verilog, "input [7:0] b") {
		t.Errorf("Input 'b' declaration not found")
	}
	if !strings.Contains(verilog, "output [7:0] sum") {
		t.Errorf("Output 'sum' declaration not found")
	}
	
	// Check assign statement
	if !strings.Contains(verilog, "assign sum = a + b;") {
		t.Errorf("Assign statement not found")
	}
	
	// Check endmodule
	if !strings.Contains(verilog, "endmodule") {
		t.Errorf("endmodule not found")
	}

	// Clean up test file
	os.Remove(testFile)
}

func TestEmitVerilogWithRegisters(t *testing.T) {
	// Create a module with registers
	m := NewModule("Counter")
	clk := m.Input("clk", 1)
	rst := m.Input("rst", 1)
	count := m.Output("count", 8)
	counter_reg := m.Reg("counter_reg", 8)

	m.SetClock(clk)
	m.SetReset(rst)
	m.Assign(count, counter_reg)

	// Emit Verilog to temporary file
	originalFile := "out.v"
	testFile := "test_counter.v"
	
	// Backup original if it exists
	if _, err := os.Stat(originalFile); err == nil {
		content, _ := os.ReadFile(originalFile)
		defer os.WriteFile(originalFile, content, 0644)
	}
	
	// Rename output file for testing
	defer func() {
		if _, err := os.Stat(originalFile); err == nil {
			os.Rename(originalFile, testFile)
		}
	}()

	EmitVerilog(m)

	// Read and verify output
	content, err := os.ReadFile(originalFile)
	if err != nil {
		t.Fatalf("Failed to read generated Verilog file: %v", err)
	}

	verilog := string(content)
	
	// Check module declaration
	if !strings.Contains(verilog, "module Counter(") {
		t.Errorf("Module declaration not found in generated Verilog")
	}
	
	// Check register declaration
	if !strings.Contains(verilog, "reg [7:0] counter_reg;") {
		t.Errorf("Register declaration not found")
	}
	
	// Check assign statement
	if !strings.Contains(verilog, "assign count = counter_reg;") {
		t.Errorf("Assign statement not found")
	}

	// Clean up test file
	os.Remove(testFile)
}

func TestEmitVerilogWithWires(t *testing.T) {
	// Create a module with wires
	m := NewModule("Logic")
	a := m.Input("a", 4)
	b := m.Input("b", 4)
	result := m.Output("result", 4)
	temp_wire := m.Wire("temp", 4)

	m.Assign(temp_wire, a.And(b))
	m.Assign(result, temp_wire)

	// Emit Verilog to temporary file
	originalFile := "out.v"
	testFile := "test_logic.v"
	
	// Backup original if it exists
	if _, err := os.Stat(originalFile); err == nil {
		content, _ := os.ReadFile(originalFile)
		defer os.WriteFile(originalFile, content, 0644)
	}
	
	// Rename output file for testing
	defer func() {
		if _, err := os.Stat(originalFile); err == nil {
			os.Rename(originalFile, testFile)
		}
	}()

	EmitVerilog(m)

	// Read and verify output
	content, err := os.ReadFile(originalFile)
	if err != nil {
		t.Fatalf("Failed to read generated Verilog file: %v", err)
	}

	verilog := string(content)
	
	// Check wire declaration
	if !strings.Contains(verilog, "wire [3:0] temp;") {
		t.Errorf("Wire declaration not found")
	}
	
	// Check assign statements
	if !strings.Contains(verilog, "assign temp = a & b;") {
		t.Errorf("First assign statement not found")
	}
	if !strings.Contains(verilog, "assign result = temp;") {
		t.Errorf("Second assign statement not found")
	}

	// Clean up test file
	os.Remove(testFile)
}

func TestEmitVerilogWithBundles(t *testing.T) {
	// Create a module with bundles
	m := NewModule("BundleTest")
	_ = m.Input("clk", 1)
	
	// Create bundles
	bundle1 := m.Bundle("io_bundle")
	bundle1.AddField("data", &hdl.Signal{Width: 8, Kind: "wire"})
	bundle1.AddField("valid", &hdl.Signal{Width: 1, Kind: "wire"})
	
	// Emit Verilog to temporary file
	originalFile := "out.v"
	testFile := "test_bundle.v"
	
	// Backup original if it exists
	if _, err := os.Stat(originalFile); err == nil {
		content, _ := os.ReadFile(originalFile)
		defer os.WriteFile(originalFile, content, 0644)
	}
	
	// Rename output file for testing
	defer func() {
		if _, err := os.Stat(originalFile); err == nil {
			os.Rename(originalFile, testFile)
		}
	}()

	EmitVerilog(m)

	// Read and verify output
	content, err := os.ReadFile(originalFile)
	if err != nil {
		t.Fatalf("Failed to read generated Verilog file: %v", err)
	}

	verilog := string(content)
	
	// Check bundle field declarations
	if !strings.Contains(verilog, "wire [7:0] io_bundle_data;") {
		t.Errorf("Bundle data field declaration not found")
	}
	if !strings.Contains(verilog, "wire [0:0] io_bundle_valid;") {
		t.Errorf("Bundle valid field declaration not found")
	}

	// Clean up test file
	os.Remove(testFile)
}

func TestEmitVerilogWithVecs(t *testing.T) {
	// Create a module with vectors
	m := NewModule("VecTest")
	_ = m.Input("clk", 1)
	
	// Create vectors
	_ = m.Vec("data_vec", 4, 8)
	
	// Emit Verilog to temporary file
	originalFile := "out.v"
	testFile := "test_vec.v"
	
	// Backup original if it exists
	if _, err := os.Stat(originalFile); err == nil {
		content, _ := os.ReadFile(originalFile)
		defer os.WriteFile(originalFile, content, 0644)
	}
	
	// Rename output file for testing
	defer func() {
		if _, err := os.Stat(originalFile); err == nil {
			os.Rename(originalFile, testFile)
		}
	}()

	EmitVerilog(m)

	// Read and verify output
	content, err := os.ReadFile(originalFile)
	if err != nil {
		t.Fatalf("Failed to read generated Verilog file: %v", err)
	}

	verilog := string(content)
	
	// Check vector element declarations
	if !strings.Contains(verilog, "wire [7:0] data_vec[0];") {
		t.Errorf("Vec element 0 declaration not found")
	}
	if !strings.Contains(verilog, "wire [7:0] data_vec[3];") {
		t.Errorf("Vec element 3 declaration not found")
	}

	// Clean up test file
	os.Remove(testFile)
}

func TestEmitVerilogWithMemory(t *testing.T) {
	// Create a module with memory
	m := NewModule("MemoryTest")
	_ = m.Input("clk", 1)
	
	// Create memory
	_ = m.SyncMem("data_mem", 16, 64)
	
	// Emit Verilog to temporary file
	originalFile := "out.v"
	testFile := "test_memory.v"
	
	// Backup original if it exists
	if _, err := os.Stat(originalFile); err == nil {
		content, _ := os.ReadFile(originalFile)
		defer os.WriteFile(originalFile, content, 0644)
	}
	
	// Rename output file for testing
	defer func() {
		if _, err := os.Stat(originalFile); err == nil {
			os.Rename(originalFile, testFile)
		}
	}()

	EmitVerilog(m)

	// Read and verify output
	content, err := os.ReadFile(originalFile)
	if err != nil {
		t.Fatalf("Failed to read generated Verilog file: %v", err)
	}

	verilog := string(content)
	
	// Check memory declaration
	if !strings.Contains(verilog, "reg [15:0] data_mem [0:63];") {
		t.Errorf("Memory declaration not found")
	}

	// Clean up test file
	os.Remove(testFile)
}

func TestEmitVerilogWithParameters(t *testing.T) {
	// Create a module with parameters
	m := NewModule("ParameterTest")
	m.SetParameter("WIDTH", 8)
	m.SetParameter("DEPTH", 256)
	
	clk := m.Input("clk", 1)
	_ = clk // Use the variable
	
	// Emit Verilog to temporary file
	originalFile := "out.v"
	testFile := "test_params.v"
	
	// Backup original if it exists
	if _, err := os.Stat(originalFile); err == nil {
		content, _ := os.ReadFile(originalFile)
		defer os.WriteFile(originalFile, content, 0644)
	}
	
	// Rename output file for testing
	defer func() {
		if _, err := os.Stat(originalFile); err == nil {
			os.Rename(originalFile, testFile)
		}
	}()

	EmitVerilog(m)

	// Read and verify output
	content, err := os.ReadFile(originalFile)
	if err != nil {
		t.Fatalf("Failed to read generated Verilog file: %v", err)
	}

	verilog := string(content)
	
	// Check parameter declarations
	if !strings.Contains(verilog, "module ParameterTest #(") {
		t.Errorf("Parameterized module declaration not found")
	}
	if !strings.Contains(verilog, "parameter WIDTH = 8") {
		t.Errorf("WIDTH parameter not found")
	}
	if !strings.Contains(verilog, "parameter DEPTH = 256") {
		t.Errorf("DEPTH parameter not found")
	}

	// Clean up test file
	os.Remove(testFile)
}

func TestEmitVerilogWithInstances(t *testing.T) {
	// Create a module with instances
	m := NewModule("TopModule")
	clk := m.Input("clk", 1)
	
	// Create an instance
	inst := m.Instance("SubModule", "sub_inst")
	inst.SetParameter("WIDTH", 16)
	inst.Connect("clk", clk)
	
	data_out := inst.IO("data_out", 16)
	output_port := m.Output("output_data", 16)
	m.Assign(output_port, data_out)
	
	// Emit Verilog to temporary file
	originalFile := "out.v"
	testFile := "test_instance.v"
	
	// Backup original if it exists
	if _, err := os.Stat(originalFile); err == nil {
		content, _ := os.ReadFile(originalFile)
		defer os.WriteFile(originalFile, content, 0644)
	}
	
	// Rename output file for testing
	defer func() {
		if _, err := os.Stat(originalFile); err == nil {
			os.Rename(originalFile, testFile)
		}
	}()

	EmitVerilog(m)

	// Read and verify output
	content, err := os.ReadFile(originalFile)
	if err != nil {
		t.Fatalf("Failed to read generated Verilog file: %v", err)
	}

	verilog := string(content)
	
	// Check instance declaration
	if !strings.Contains(verilog, "SubModule #(") {
		t.Errorf("Parameterized instance not found")
	}
	if !strings.Contains(verilog, ".WIDTH(16)") {
		t.Errorf("Instance parameter not found")
	}
	if !strings.Contains(verilog, ") sub_inst (") {
		t.Errorf("Instance name not found")
	}
	if !strings.Contains(verilog, ".clk(clk)") {
		t.Errorf("Instance connection not found")
	}

	// Clean up test file
	os.Remove(testFile)
}