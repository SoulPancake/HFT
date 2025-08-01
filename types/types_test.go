package hdl

import (
	"testing"
)

func TestSignalArithmeticOperations(t *testing.T) {
	a := &Signal{Name: "a", Width: 8, Kind: "wire"}
	b := &Signal{Name: "b", Width: 8, Kind: "wire"}

	tests := []struct {
		name     string
		operation func() *Signal
		expected string
	}{
		{"Add", func() *Signal { return a.Add(b) }, "a + b"},
		{"Sub", func() *Signal { return a.Sub(b) }, "a - b"},
		{"Mul", func() *Signal { return a.Mul(b) }, "a * b"},
		{"Div", func() *Signal { return a.Div(b) }, "a / b"},
		{"Mod", func() *Signal { return a.Mod(b) }, "a % b"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.operation(); got.Name != tt.expected {
				t.Errorf("%s operation = %v, want %v", tt.name, got.Name, tt.expected)
			}
		})
	}
}

func TestSignalBitwiseOperations(t *testing.T) {
	a := &Signal{Name: "a", Width: 8, Kind: "wire"}
	b := &Signal{Name: "b", Width: 8, Kind: "wire"}

	tests := []struct {
		name      string
		operation func() *Signal
		expected  string
	}{
		{"And", func() *Signal { return a.And(b) }, "a & b"},
		{"Or", func() *Signal { return a.Or(b) }, "a | b"},
		{"Xor", func() *Signal { return a.Xor(b) }, "a ^ b"},
		{"Not", func() *Signal { return a.Not() }, "~a"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.operation(); got.Name != tt.expected {
				t.Errorf("%s operation = %v, want %v", tt.name, got.Name, tt.expected)
			}
		})
	}
}

func TestSignalLogicalOperations(t *testing.T) {
	a := &Signal{Name: "a", Width: 8, Kind: "wire"}
	b := &Signal{Name: "b", Width: 8, Kind: "wire"}

	tests := []struct {
		name      string
		operation func() *Signal
		expected  string
	}{
		{"LogicAnd", func() *Signal { return a.LogicAnd(b) }, "a && b"},
		{"LogicOr", func() *Signal { return a.LogicOr(b) }, "a || b"},
		{"LogicNot", func() *Signal { return a.LogicNot() }, "!a"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.operation(); got.Name != tt.expected {
				t.Errorf("%s operation = %v, want %v", tt.name, got.Name, tt.expected)
			}
		})
	}
}

func TestSignalComparisonOperations(t *testing.T) {
	a := &Signal{Name: "a", Width: 8, Kind: "wire"}
	b := &Signal{Name: "b", Width: 8, Kind: "wire"}

	tests := []struct {
		name      string
		operation func() *Signal
		expected  string
	}{
		{"Eq", func() *Signal { return a.Eq(b) }, "a == b"},
		{"Neq", func() *Signal { return a.Neq(b) }, "a != b"},
		{"Lt", func() *Signal { return a.Lt(b) }, "a < b"},
		{"Lte", func() *Signal { return a.Lte(b) }, "a <= b"},
		{"Gt", func() *Signal { return a.Gt(b) }, "a > b"},
		{"Gte", func() *Signal { return a.Gte(b) }, "a >= b"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.operation(); got.Name != tt.expected {
				t.Errorf("%s operation = %v, want %v", tt.name, got.Name, tt.expected)
			}
		})
	}
}

func TestSignalShiftOperations(t *testing.T) {
	a := &Signal{Name: "a", Width: 8, Kind: "wire"}

	tests := []struct {
		name      string
		operation func() *Signal
		expected  string
	}{
		{"Shl", func() *Signal { return a.Shl(2) }, "a << 2"},
		{"Shr", func() *Signal { return a.Shr(3) }, "a >> 3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.operation(); got.Name != tt.expected {
				t.Errorf("%s operation = %v, want %v", tt.name, got.Name, tt.expected)
			}
		})
	}
}

func TestSignalBitSelection(t *testing.T) {
	a := &Signal{Name: "data", Width: 8, Kind: "wire"}

	tests := []struct {
		name     string
		high     int
		low      int
		expected string
	}{
		{"Single bit", 3, 3, "data[3]"},
		{"Bit range", 7, 4, "data[7:4]"},
		{"Full range", 7, 0, "data[7:0]"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := a.Bits(tt.high, tt.low); got.Name != tt.expected {
				t.Errorf("Bits(%d, %d) = %v, want %v", tt.high, tt.low, got.Name, tt.expected)
			}
		})
	}
}

func TestCatOperation(t *testing.T) {
	a := &Signal{Name: "a", Width: 4, Kind: "wire"}
	b := &Signal{Name: "b", Width: 4, Kind: "wire"}
	c := &Signal{Name: "c", Width: 4, Kind: "wire"}

	tests := []struct {
		name     string
		signals  []*Signal
		expected string
	}{
		{"Empty", []*Signal{}, ""},
		{"Single signal", []*Signal{a}, "{a}"},
		{"Two signals", []*Signal{a, b}, "{a, b}"},
		{"Three signals", []*Signal{a, b, c}, "{a, b, c}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Cat(tt.signals...); got.Name != tt.expected {
				t.Errorf("Cat() = %v, want %v", got.Name, tt.expected)
			}
		})
	}
}

func TestFillOperation(t *testing.T) {
	tests := []struct {
		name     string
		n        int
		value    string
		expected string
	}{
		{"Fill zeros", 4, "1'b0", "{4{1'b0}}"},
		{"Fill ones", 8, "1'b1", "{8{1'b1}}"},
		{"Fill pattern", 3, "2'b01", "{3{2'b01}}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Fill(tt.n, tt.value); got.Name != tt.expected {
				t.Errorf("Fill(%d, %s) = %v, want %v", tt.n, tt.value, got.Name, tt.expected)
			}
		})
	}
}

func TestMuxOperation(t *testing.T) {
	sel := &Signal{Name: "sel", Width: 2, Kind: "wire"}
	a := &Signal{Name: "a", Width: 8, Kind: "wire"}
	b := &Signal{Name: "b", Width: 8, Kind: "wire"}
	c := &Signal{Name: "c", Width: 8, Kind: "wire"}

	tests := []struct {
		name     string
		inputs   []*Signal
		expected string
	}{
		{"Empty inputs", []*Signal{}, ""},
		{"Single input", []*Signal{a}, "a"},
		{"Two inputs", []*Signal{a, b}, "(b : (sel == 0) ? a)"},
		{"Three inputs", []*Signal{a, b, c}, "(c : (sel == 1) ? b : (sel == 0) ? a)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Mux(sel, tt.inputs...); got.Name != tt.expected {
				t.Errorf("Mux() = %v, want %v", got.Name, tt.expected)
			}
		})
	}
}

func TestModuleCreation(t *testing.T) {
	m := &Module{
		Name:    "TestModule",
		Inputs:  []*Signal{},
		Outputs: []*Signal{},
		Wires:   []*Signal{},
		Regs:    []*Signal{},
		Assigns: []string{},
		Always:  []string{},
	}

	// Test input creation
	inp := m.Input("clk", 1)
	if inp.Name != "clk" || inp.Width != 1 || inp.Kind != "input" {
		t.Errorf("Input creation failed: got %+v", inp)
	}
	if len(m.Inputs) != 1 {
		t.Errorf("Input not added to module: got %d inputs", len(m.Inputs))
	}

	// Test output creation
	out := m.Output("result", 8)
	if out.Name != "result" || out.Width != 8 || out.Kind != "output" {
		t.Errorf("Output creation failed: got %+v", out)
	}
	if len(m.Outputs) != 1 {
		t.Errorf("Output not added to module: got %d outputs", len(m.Outputs))
	}

	// Test wire creation
	wire := m.Wire("temp", 4)
	if wire.Name != "temp" || wire.Width != 4 || wire.Kind != "wire" {
		t.Errorf("Wire creation failed: got %+v", wire)
	}
	if len(m.Wires) != 1 {
		t.Errorf("Wire not added to module: got %d wires", len(m.Wires))
	}

	// Test reg creation
	reg := m.Reg("counter", 16)
	if reg.Name != "counter" || reg.Width != 16 || reg.Kind != "reg" {
		t.Errorf("Reg creation failed: got %+v", reg)
	}
	if len(m.Regs) != 1 {
		t.Errorf("Reg not added to module: got %d regs", len(m.Regs))
	}
}

func TestModuleAssign(t *testing.T) {
	m := &Module{
		Name:    "TestModule",
		Assigns: []string{},
	}

	out := &Signal{Name: "out", Width: 8, Kind: "wire"}
	a := &Signal{Name: "a", Width: 8, Kind: "wire"}
	b := &Signal{Name: "b", Width: 8, Kind: "wire"}
	
	add_result := a.Add(b)
	m.Assign(out, add_result)
	
	expected := "assign out = a + b;"
	if len(m.Assigns) != 1 || m.Assigns[0] != expected {
		t.Errorf("Assign failed: got %v, want %v", m.Assigns, []string{expected})
	}
}

func TestModuleClockAndReset(t *testing.T) {
	m := &Module{Name: "TestModule"}
	
	clk := &Signal{Name: "clk", Width: 1, Kind: "input"}
	rst := &Signal{Name: "rst", Width: 1, Kind: "input"}
	
	m.SetClock(clk)
	m.SetReset(rst)
	
	if m.Clock != clk {
		t.Errorf("Clock not set correctly: got %+v, want %+v", m.Clock, clk)
	}
	if m.Reset != rst {
		t.Errorf("Reset not set correctly: got %+v, want %+v", m.Reset, rst)
	}
}

func TestBundleCreation(t *testing.T) {
	m := &Module{
		Name:    "TestModule",
		Bundles: []*Bundle{},
	}

	bundle := m.Bundle("test_bundle")
	
	if bundle.Name != "test_bundle" {
		t.Errorf("Bundle name incorrect: got %v, want %v", bundle.Name, "test_bundle")
	}
	if bundle.Fields == nil {
		t.Errorf("Bundle fields should be initialized")
	}
	if len(m.Bundles) != 1 {
		t.Errorf("Bundle not added to module: got %d bundles", len(m.Bundles))
	}
}

func TestBundleFields(t *testing.T) {
	bundle := &Bundle{
		Name:   "test_bundle",
		Fields: make(map[string]*Signal),
	}

	signal := &Signal{Name: "data", Width: 8, Kind: "wire"}
	bundle.AddField("data", signal)
	
	if len(bundle.Fields) != 1 {
		t.Errorf("Field not added: got %d fields", len(bundle.Fields))
	}
	
	field := bundle.GetField("data")
	if field == nil {
		t.Errorf("Field not found")
	}
	if field.Name != "test_bundle_data" {
		t.Errorf("Field name incorrect: got %v, want %v", field.Name, "test_bundle_data")
	}
}

func TestBundleConnect(t *testing.T) {
	bundle1 := &Bundle{
		Name:   "bundle1",
		Fields: make(map[string]*Signal),
	}
	bundle2 := &Bundle{
		Name:   "bundle2",
		Fields: make(map[string]*Signal),
	}

	signal1 := &Signal{Name: "data", Width: 8, Kind: "wire"}
	signal2 := &Signal{Name: "data", Width: 8, Kind: "wire"}
	
	bundle1.AddField("data", signal1)
	bundle2.AddField("data", signal2)
	
	connections := bundle1.Connect(bundle2)
	
	if len(connections) != 1 {
		t.Errorf("Expected 1 connection, got %d", len(connections))
	}
	
	expected := "assign bundle1_data = bundle2_data;"
	if connections[0] != expected {
		t.Errorf("Connection incorrect: got %v, want %v", connections[0], expected)
	}
}

func TestVecCreation(t *testing.T) {
	m := &Module{
		Name: "TestModule",
		Vecs: []*Vec{},
	}

	vec := m.Vec("test_vec", 4, 8)
	
	if vec.Name != "test_vec" {
		t.Errorf("Vec name incorrect: got %v, want %v", vec.Name, "test_vec")
	}
	if vec.Size != 4 {
		t.Errorf("Vec size incorrect: got %v, want %v", vec.Size, 4)
	}
	if vec.Width != 8 {
		t.Errorf("Vec width incorrect: got %v, want %v", vec.Width, 8)
	}
	if len(vec.Elements) != 4 {
		t.Errorf("Vec elements not created: got %d elements", len(vec.Elements))
	}
	if len(m.Vecs) != 1 {
		t.Errorf("Vec not added to module: got %d vecs", len(m.Vecs))
	}
}

func TestVecAccess(t *testing.T) {
	m := &Module{Name: "TestModule", Vecs: []*Vec{}}
	vec := m.Vec("test_vec", 4, 8)
	
	// Test valid access
	elem := vec.Get(2)
	if elem == nil {
		t.Errorf("Element should not be nil")
	}
	if elem.Name != "test_vec[2]" {
		t.Errorf("Element name incorrect: got %v, want %v", elem.Name, "test_vec[2]")
	}
	
	// Test out of bounds access
	elem = vec.Get(5)
	if elem != nil {
		t.Errorf("Out of bounds access should return nil")
	}
	
	// Test set
	newSignal := &Signal{Name: "new_signal", Width: 8, Kind: "wire"}
	vec.Set(1, newSignal)
	if vec.Elements[1] != newSignal {
		t.Errorf("Set operation failed")
	}
}

func TestVecConnect(t *testing.T) {
	m := &Module{Name: "TestModule", Vecs: []*Vec{}}
	vec1 := m.Vec("vec1", 3, 8)
	vec2 := m.Vec("vec2", 3, 8)
	
	connections := vec1.Connect(vec2)
	
	if len(connections) != 3 {
		t.Errorf("Expected 3 connections, got %d", len(connections))
	}
	
	expected := "assign vec1[0] = vec2[0];"
	if connections[0] != expected {
		t.Errorf("First connection incorrect: got %v, want %v", connections[0], expected)
	}
}

func TestMemoryCreation(t *testing.T) {
	m := &Module{
		Name:     "TestModule",
		Memories: []*Memory{},
	}

	mem := m.SyncMem("test_mem", 32, 256)
	
	if mem.Name != "test_mem" {
		t.Errorf("Memory name incorrect: got %v, want %v", mem.Name, "test_mem")
	}
	if mem.Width != 32 {
		t.Errorf("Memory width incorrect: got %v, want %v", mem.Width, 32)
	}
	if mem.Depth != 256 {
		t.Errorf("Memory depth incorrect: got %v, want %v", mem.Depth, 256)
	}
	if mem.Kind != "sync" {
		t.Errorf("Memory kind incorrect: got %v, want %v", mem.Kind, "sync")
	}
	if mem.AddrWidth != 8 {
		t.Errorf("Memory address width incorrect: got %v, want %v", mem.AddrWidth, 8)
	}
	if len(m.Memories) != 1 {
		t.Errorf("Memory not added to module: got %d memories", len(m.Memories))
	}
}

func TestMemoryOperations(t *testing.T) {
	m := &Module{Name: "TestModule", Memories: []*Memory{}}
	mem := m.SyncMem("test_mem", 16, 64)
	
	addr := &Signal{Name: "addr", Width: 6, Kind: "wire"}
	data := &Signal{Name: "data", Width: 16, Kind: "wire"}
	enable := &Signal{Name: "we", Width: 1, Kind: "wire"}
	
	// Test read
	readData := mem.Read(addr)
	expected := "test_mem[addr]"
	if readData.Name != expected {
		t.Errorf("Memory read incorrect: got %v, want %v", readData.Name, expected)
	}
	
	// Test write
	writeStmt := mem.Write(addr, data, enable)
	expectedWrite := "if (we) test_mem[addr] <= data;"
	if writeStmt != expectedWrite {
		t.Errorf("Memory write incorrect: got %v, want %v", writeStmt, expectedWrite)
	}
}

func TestAsyncMemory(t *testing.T) {
	m := &Module{Name: "TestModule", Memories: []*Memory{}}
	mem := m.AsyncMem("async_mem", 8, 32)
	
	if mem.Kind != "async" {
		t.Errorf("Async memory kind incorrect: got %v, want %v", mem.Kind, "async")
	}
	
	addr := &Signal{Name: "addr", Width: 5, Kind: "wire"}
	data := &Signal{Name: "data", Width: 8, Kind: "wire"}
	enable := &Signal{Name: "we", Width: 1, Kind: "wire"}
	
	writeStmt := mem.Write(addr, data, enable)
	expected := "assign async_mem[addr] = (we) ? data : async_mem[addr];"
	if writeStmt != expected {
		t.Errorf("Async memory write incorrect: got %v, want %v", writeStmt, expected)
	}
}