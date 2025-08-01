package hdl

import (
	"testing"
)

func TestWidthHelpers(t *testing.T) {
	tests := []struct {
		value    int
		expected Width
	}{
		{0, 1},
		{1, 1},
		{2, 2},
		{3, 2},
		{4, 3},
		{7, 3},
		{8, 4},
		{255, 8},
		{256, 9},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := WidthOf(tt.value); got != tt.expected {
				t.Errorf("WidthOf(%d) = %v, want %v", tt.value, got, tt.expected)
			}
		})
	}
}

func TestMaxWidth(t *testing.T) {
	tests := []struct {
		a, b     Width
		expected Width
	}{
		{4, 8, 8},
		{8, 4, 8},
		{1, 1, 1},
		{16, 8, 16},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := MaxWidth(tt.a, tt.b); got != tt.expected {
				t.Errorf("MaxWidth(%d, %d) = %v, want %v", tt.a, tt.b, got, tt.expected)
			}
		})
	}
}

func TestWidthBits(t *testing.T) {
	tests := []struct {
		width    Width
		expected string
	}{
		{1, "[0:0]"},
		{8, "[7:0]"},
		{16, "[15:0]"},
		{32, "[31:0]"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := tt.width.Bits(); got != tt.expected {
				t.Errorf("Width(%d).Bits() = %v, want %v", tt.width, got, tt.expected)
			}
		})
	}
}

func TestSignalWithWidth(t *testing.T) {
	original := &Signal{Name: "test", Width: 8, Kind: "wire"}
	modified := original.WithWidth(16)
	
	if modified.Width != 16 {
		t.Errorf("WithWidth failed: got width %d, want 16", modified.Width)
	}
	if modified.Name != "test" {
		t.Errorf("WithWidth changed name: got %s, want test", modified.Name)
	}
	if original.Width != 8 {
		t.Errorf("WithWidth modified original: got width %d, want 8", original.Width)
	}
}

func TestSignalCheckWidth(t *testing.T) {
	sig1 := &Signal{Name: "a", Width: 8, Kind: "wire"}
	sig2 := &Signal{Name: "b", Width: 8, Kind: "wire"}
	sig3 := &Signal{Name: "c", Width: 16, Kind: "wire"}
	
	// Same width should pass
	if err := sig1.CheckWidth(sig2); err != nil {
		t.Errorf("CheckWidth failed for same widths: %v", err)
	}
	
	// Different width should fail
	if err := sig1.CheckWidth(sig3); err == nil {
		t.Errorf("CheckWidth should have failed for different widths")
	}
}

func TestWidthAwareArithmetic(t *testing.T) {
	a := &Signal{Name: "a", Width: 8, Kind: "wire"}
	b := &Signal{Name: "b", Width: 16, Kind: "wire"}
	
	// Addition should take max width
	add_result := a.Add(b)
	if add_result.Width != 16 {
		t.Errorf("Add width inference failed: got %d, want 16", add_result.Width)
	}
	
	// Multiplication should add widths
	mul_result := a.Mul(b)
	if mul_result.Width != 24 {
		t.Errorf("Mul width inference failed: got %d, want 24", mul_result.Width)
	}
	
	// Comparison should return 1-bit result
	eq_result := a.Eq(b)
	if eq_result.Width != 1 {
		t.Errorf("Eq width inference failed: got %d, want 1", eq_result.Width)
	}
}

func TestShiftOperations(t *testing.T) {
	sig := &Signal{Name: "data", Width: 8, Kind: "wire"}
	
	// Left shift should increase width
	shl_result := sig.Shl(2)
	if shl_result.Width != 10 {
		t.Errorf("Shl width inference failed: got %d, want 10", shl_result.Width)
	}
	
	// Right shift should decrease width
	shr_result := sig.Shr(2)
	if shr_result.Width != 6 {
		t.Errorf("Shr width inference failed: got %d, want 6", shr_result.Width)
	}
	
	// Right shift should not go below 1
	shr_result2 := sig.Shr(10)
	if shr_result2.Width != 1 {
		t.Errorf("Shr width inference failed: got %d, want 1", shr_result2.Width)
	}
}

func TestBitSelection(t *testing.T) {
	sig := &Signal{Name: "data", Width: 16, Kind: "wire"}
	
	// Single bit selection
	bit_result := sig.Bits(7, 7)
	if bit_result.Width != 1 {
		t.Errorf("Single bit selection width failed: got %d, want 1", bit_result.Width)
	}
	
	// Range selection
	range_result := sig.Bits(15, 8)
	if range_result.Width != 8 {
		t.Errorf("Range selection width failed: got %d, want 8", range_result.Width)
	}
}

func TestCatWidthInference(t *testing.T) {
	sig1 := &Signal{Name: "a", Width: 8, Kind: "wire"}
	sig2 := &Signal{Name: "b", Width: 16, Kind: "wire"}
	sig3 := &Signal{Name: "c", Width: 4, Kind: "wire"}
	
	cat_result := Cat(sig1, sig2, sig3)
	expected_width := Width(8 + 16 + 4)
	
	if cat_result.Width != expected_width {
		t.Errorf("Cat width inference failed: got %d, want %d", cat_result.Width, expected_width)
	}
}

func TestMuxWidthInference(t *testing.T) {
	sel := &Signal{Name: "sel", Width: 2, Kind: "wire"}
	input1 := &Signal{Name: "in1", Width: 8, Kind: "wire"}
	input2 := &Signal{Name: "in2", Width: 16, Kind: "wire"}
	input3 := &Signal{Name: "in3", Width: 4, Kind: "wire"}
	
	mux_result := Mux(sel, input1, input2, input3)
	
	// Should take the width of the widest input
	if mux_result.Width != 16 {
		t.Errorf("Mux width inference failed: got %d, want 16", mux_result.Width)
	}
}

func TestLiteralCreation(t *testing.T) {
	lit8 := Lit(255, 8)
	if lit8.Width != 8 {
		t.Errorf("Lit width failed: got %d, want 8", lit8.Width)
	}
	if lit8.Name != "8'hff" {
		t.Errorf("Lit name failed: got %s, want 8'hff", lit8.Name)
	}
	
	bool_true := Bool(true)
	if bool_true.Width != 1 {
		t.Errorf("Bool width failed: got %d, want 1", bool_true.Width)
	}
	if bool_true.Name != "1'b1" {
		t.Errorf("Bool name failed: got %s, want 1'b1", bool_true.Name)
	}
	
	bool_false := Bool(false)
	if bool_false.Name != "1'b0" {
		t.Errorf("Bool name failed: got %s, want 1'b0", bool_false.Name)
	}
}

func TestModuleInstance(t *testing.T) {
	m := &Module{
		Name:      "TestModule",
		Instances: []*ModuleInstance{},
	}
	
	inst := m.Instance("SubModule", "sub_inst")
	
	if inst.ModuleName != "SubModule" {
		t.Errorf("Instance module name failed: got %s, want SubModule", inst.ModuleName)
	}
	if inst.InstanceName != "sub_inst" {
		t.Errorf("Instance name failed: got %s, want sub_inst", inst.InstanceName)
	}
	if len(m.Instances) != 1 {
		t.Errorf("Instance not added to module: got %d instances", len(m.Instances))
	}
}

func TestModuleInstanceConnections(t *testing.T) {
	m := &Module{Name: "TestModule", Instances: []*ModuleInstance{}}
	inst := m.Instance("SubModule", "sub_inst")
	
	// Test parameter setting
	inst.SetParameter("WIDTH", 8)
	if inst.Parameters["WIDTH"] != 8 {
		t.Errorf("Parameter setting failed: got %v, want 8", inst.Parameters["WIDTH"])
	}
	
	// Test signal connection
	clk := &Signal{Name: "clk", Width: 1, Kind: "input"}
	inst.Connect("clk", clk)
	if inst.Connections["clk"] != clk {
		t.Errorf("Signal connection failed")
	}
	
	// Test IO creation
	data_out := inst.IO("data_out", 16)
	if data_out.Width != 16 {
		t.Errorf("IO signal width failed: got %d, want 16", data_out.Width)
	}
	if data_out.Name != "sub_inst_data_out" {
		t.Errorf("IO signal name failed: got %s, want sub_inst_data_out", data_out.Name)
	}
}

func TestModuleParameters(t *testing.T) {
	m := &Module{Name: "TestModule"}
	
	m.SetParameter("WIDTH", 8)
	m.SetParameter("DEPTH", 256)
	
	if m.Parameters["WIDTH"] != 8 {
		t.Errorf("Module parameter setting failed: got %v, want 8", m.Parameters["WIDTH"])
	}
	if m.Parameters["DEPTH"] != 256 {
		t.Errorf("Module parameter setting failed: got %v, want 256", m.Parameters["DEPTH"])
	}
}