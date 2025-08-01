package hdl

import (
	"strings"
	"testing"
)

// === CLOCK DOMAIN TESTS ===

func TestClockDomainCreation(t *testing.T) {
	m := &Module{
		Name:         "TestModule",
		ClockDomains: []*ClockDomain{},
	}
	
	clk := &Signal{Name: "clk_100m", Width: 1, Kind: "input"}
	rst := &Signal{Name: "rst_n", Width: 1, Kind: "input"}
	
	cd := m.NewClockDomain("main_domain", clk, rst)
	
	if cd.Name != "main_domain" {
		t.Errorf("Expected clock domain name 'main_domain', got '%s'", cd.Name)
	}
	if cd.Clock != clk {
		t.Errorf("Clock signal not properly assigned")
	}
	if cd.Reset != rst {
		t.Errorf("Reset signal not properly assigned")
	}
	if len(m.ClockDomains) != 1 {
		t.Errorf("Expected 1 clock domain, got %d", len(m.ClockDomains))
	}
}

func TestClockDomainFrequency(t *testing.T) {
	clk := &Signal{Name: "clk", Width: 1, Kind: "input"}
	rst := &Signal{Name: "rst", Width: 1, Kind: "input"}
	
	cd := &ClockDomain{Name: "test", Clock: clk, Reset: rst}
	cd.SetFrequency(100000000) // 100 MHz
	
	if cd.Frequency != 100000000 {
		t.Errorf("Expected frequency 100000000, got %d", cd.Frequency)
	}
}

func TestSignalClockDomainAssociation(t *testing.T) {
	clk := &Signal{Name: "clk", Width: 1, Kind: "input"}
	rst := &Signal{Name: "rst", Width: 1, Kind: "input"}
	cd := &ClockDomain{Name: "test", Clock: clk, Reset: rst}
	
	sig := &Signal{Name: "data", Width: 8, Kind: "wire"}
	sig.WithClockDomain(cd)
	
	if sig.ClockDomain != cd {
		t.Errorf("Signal not properly associated with clock domain")
	}
}

func TestCDCSynchronizer(t *testing.T) {
	m := &Module{
		Name:         "CDCTest",
		Regs:         []*Signal{},
		Always:       []string{},
		ClockDomains: []*ClockDomain{},
	}
	
	srcClk := &Signal{Name: "src_clk", Width: 1, Kind: "input"}
	srcRst := &Signal{Name: "src_rst", Width: 1, Kind: "input"}
	destClk := &Signal{Name: "dest_clk", Width: 1, Kind: "input"}
	destRst := &Signal{Name: "dest_rst", Width: 1, Kind: "input"}
	
	srcDomain := m.NewClockDomain("src", srcClk, srcRst)
	destDomain := m.NewClockDomain("dest", destClk, destRst)
	
	srcSignal := &Signal{Name: "src_data", Width: 8, Kind: "wire"}
	syncSignal := m.CDCSynchronizer("sync", srcSignal, destDomain, 2)
	
	// Use srcDomain in a comment to avoid unused variable
	_ = srcDomain
	
	if syncSignal.Width != srcSignal.Width {
		t.Errorf("Synchronized signal width mismatch: expected %d, got %d", srcSignal.Width, syncSignal.Width)
	}
	
	if !strings.Contains(syncSignal.Name, "sync_stage1") {
		t.Errorf("Expected synchronized signal to reference sync stages")
	}
	
	// Check that synchronizer logic was added
	if len(m.Always) < 2 {
		t.Errorf("Expected at least 2 always blocks for 2-stage synchronizer, got %d", len(m.Always))
	}
}

func TestAsyncFIFO(t *testing.T) {
	m := &Module{
		Name:         "AsyncFIFOTest",
		Inputs:       []*Signal{},
		Outputs:      []*Signal{},
		Regs:         []*Signal{},
		Memories:     []*Memory{},
		Always:       []string{},
		ClockDomains: []*ClockDomain{},
	}
	
	wrClk := &Signal{Name: "wr_clk", Width: 1, Kind: "input"}
	rdClk := &Signal{Name: "rd_clk", Width: 1, Kind: "input"}
	wrRst := &Signal{Name: "wr_rst", Width: 1, Kind: "input"}
	rdRst := &Signal{Name: "rd_rst", Width: 1, Kind: "input"}
	
	wrDomain := m.NewClockDomain("wr", wrClk, wrRst)
	rdDomain := m.NewClockDomain("rd", rdClk, rdRst)
	
	wrData, rdData, wrFull, rdEmpty := m.AsyncFIFO("test_fifo", 32, 16, wrDomain, rdDomain)
	
	// Check signal widths
	if wrData.Width != 32 {
		t.Errorf("Expected write data width 32, got %d", wrData.Width)
	}
	if rdData.Width != 32 {
		t.Errorf("Expected read data width 32, got %d", rdData.Width)
	}
	if wrFull.Width != 1 {
		t.Errorf("Expected write full width 1, got %d", wrFull.Width)
	}
	if rdEmpty.Width != 1 {
		t.Errorf("Expected read empty width 1, got %d", rdEmpty.Width)
	}
	
	// Check clock domain assignments
	if wrData.ClockDomain != wrDomain {
		t.Errorf("Write data not in write domain")
	}
	if rdData.ClockDomain != rdDomain {
		t.Errorf("Read data not in read domain")
	}
}

// === MUTEX TESTS ===

func TestMutexCreation(t *testing.T) {
	m := &Module{
		Name:    "MutexTest",
		Inputs:  []*Signal{},
		Outputs: []*Signal{},
		Mutexes: []*Mutex{},
	}
	
	mutex := m.Mutex("bus_arbiter", 4, "round_robin")
	
	if mutex.Name != "bus_arbiter" {
		t.Errorf("Expected mutex name 'bus_arbiter', got '%s'", mutex.Name)
	}
	if mutex.Width != 4 {
		t.Errorf("Expected mutex width 4, got %d", mutex.Width)
	}
	if mutex.Arbitration != "round_robin" {
		t.Errorf("Expected arbitration 'round_robin', got '%s'", mutex.Arbitration)
	}
	if len(mutex.Requests) != 4 {
		t.Errorf("Expected 4 request signals, got %d", len(mutex.Requests))
	}
	if len(mutex.Grants) != 4 {
		t.Errorf("Expected 4 grant signals, got %d", len(mutex.Grants))
	}
	if len(m.Mutexes) != 1 {
		t.Errorf("Expected 1 mutex in module, got %d", len(m.Mutexes))
	}
}

func TestMutexRoundRobin(t *testing.T) {
	m := &Module{
		Name:    "RoundRobinTest",
		Inputs:  []*Signal{},
		Outputs: []*Signal{},
		Regs:    []*Signal{},
		Always:  []string{},
		Mutexes: []*Mutex{},
	}
	
	clk := &Signal{Name: "clk", Width: 1, Kind: "input"}
	mutex := m.Mutex("arbiter", 3, "round_robin")
	mutex.GenerateRoundRobin(m, clk)
	
	// Check that counter register was created
	found := false
	for _, reg := range m.Regs {
		if strings.Contains(reg.Name, "counter") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected counter register for round-robin arbitration")
	}
	
	// Check that always block was generated
	if len(m.Always) == 0 {
		t.Errorf("Expected always block for round-robin arbitration")
	}
}

func TestMutexPriority(t *testing.T) {
	m := &Module{
		Name:    "PriorityTest",
		Assigns: []string{},
		Mutexes: []*Mutex{},
	}
	
	mutex := m.Mutex("priority_arbiter", 3, "priority")
	mutex.GeneratePriority(m)
	
	// Check that assign statements were generated
	if len(m.Assigns) != 3 {
		t.Errorf("Expected 3 assign statements for priority arbitration, got %d", len(m.Assigns))
	}
}

// === POLYMORPHISM TESTS ===

func TestModuleTemplateCreation(t *testing.T) {
	template := NewModuleTemplate("GenericAdder", []string{"WIDTH"})
	
	if template.Name != "GenericAdder" {
		t.Errorf("Expected template name 'GenericAdder', got '%s'", template.Name)
	}
	if len(template.TypeParams) != 1 || template.TypeParams[0] != "WIDTH" {
		t.Errorf("Expected type parameter 'WIDTH', got %v", template.TypeParams)
	}
}

func TestTemplateConstraints(t *testing.T) {
	template := NewModuleTemplate("Test", []string{"WIDTH", "DEPTH"})
	template.AddConstraint("WIDTH", Width(0))
	template.AddConstraint("DEPTH", int(0))
	
	if len(template.Params) != 2 {
		t.Errorf("Expected 2 constraints, got %d", len(template.Params))
	}
	
	if _, exists := template.Params["WIDTH"]; !exists {
		t.Errorf("WIDTH constraint not found")
	}
	if _, exists := template.Params["DEPTH"]; !exists {
		t.Errorf("DEPTH constraint not found")
	}
}

func TestTemplateInstantiation(t *testing.T) {
	template := NewModuleTemplate("GenericAdder", []string{"WIDTH"})
	template.AddConstraint("WIDTH", Width(0))
	
	// Set generator function
	template.SetGenerator(func(params map[string]interface{}) *Module {
		width := params["WIDTH"].(Width)
		m := &Module{
			Name:       "GenericAdder",
			Inputs:     []*Signal{},
			Outputs:    []*Signal{},
			Parameters: make(map[string]interface{}),
		}
		m.Input("a", width)
		m.Input("b", width)
		m.Output("sum", width+1)
		m.SetParameter("WIDTH", width)
		return m
	})
	
	// Instantiate template
	hostModule := &Module{Name: "Host"}
	instance := hostModule.InstantiateTemplate(template, "adder_8bit", map[string]interface{}{
		"WIDTH": Width(8),
	})
	
	if instance.Name != "adder_8bit" {
		t.Errorf("Expected instance name 'adder_8bit', got '%s'", instance.Name)
	}
	
	if len(instance.Inputs) != 2 {
		t.Errorf("Expected 2 inputs, got %d", len(instance.Inputs))
	}
	if len(instance.Outputs) != 1 {
		t.Errorf("Expected 1 output, got %d", len(instance.Outputs))
	}
	
	// Check that WIDTH parameter was set
	if width, exists := instance.Parameters["WIDTH"]; !exists || width != Width(8) {
		t.Errorf("WIDTH parameter not properly set")
	}
}

func TestTemplateValidation(t *testing.T) {
	template := NewModuleTemplate("Test", []string{"WIDTH"})
	template.AddConstraint("WIDTH", Width(0))
	
	template.SetGenerator(func(params map[string]interface{}) *Module {
		return &Module{Name: "Test", Parameters: make(map[string]interface{})}
	})
	
	hostModule := &Module{Name: "Host"}
	
	// Test missing type argument
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for missing type argument")
		}
	}()
	
	hostModule.InstantiateTemplate(template, "test", map[string]interface{}{})
}

func TestFIFOTemplate(t *testing.T) {
	template := FIFOTemplate()
	
	if template.Name != "GenericFIFO" {
		t.Errorf("Expected template name 'GenericFIFO', got '%s'", template.Name)
	}
	
	if len(template.TypeParams) != 2 {
		t.Errorf("Expected 2 type parameters, got %d", len(template.TypeParams))
	}
	
	// Test instantiation
	hostModule := &Module{Name: "Host"}
	instance := hostModule.InstantiateTemplate(template, "fifo_32x16", map[string]interface{}{
		"DATA_WIDTH": Width(32),
		"DEPTH":      16,
	})
	
	if instance.Name != "fifo_32x16" {
		t.Errorf("Expected instance name 'fifo_32x16', got '%s'", instance.Name)
	}
	
	// Check FIFO interface
	expectedInputs := 5  // clk, rst, wr_data, wr_en, rd_en
	expectedOutputs := 3 // wr_full, rd_data, rd_empty
	
	if len(instance.Inputs) != expectedInputs {
		t.Errorf("Expected %d inputs, got %d", expectedInputs, len(instance.Inputs))
	}
	if len(instance.Outputs) != expectedOutputs {
		t.Errorf("Expected %d outputs, got %d", expectedOutputs, len(instance.Outputs))
	}
	
	// Check memory was created
	if len(instance.Memories) != 1 {
		t.Errorf("Expected 1 memory, got %d", len(instance.Memories))
	}
}

func TestRegisterFileTemplate(t *testing.T) {
	template := RegisterFileTemplate()
	
	if template.Name != "GenericRegisterFile" {
		t.Errorf("Expected template name 'GenericRegisterFile', got '%s'", template.Name)
	}
	
	// Test instantiation
	hostModule := &Module{Name: "Host"}
	instance := hostModule.InstantiateTemplate(template, "regfile_32x5", map[string]interface{}{
		"DATA_WIDTH": Width(32),
		"ADDR_WIDTH": Width(5),
	})
	
	if instance.Name != "regfile_32x5" {
		t.Errorf("Expected instance name 'regfile_32x5', got '%s'", instance.Name)
	}
	
	// Check register file interface
	expectedInputs := 7  // clk, rst, wr_addr, wr_data, wr_en, rd1_addr, rd2_addr
	expectedOutputs := 2 // rd1_data, rd2_data
	
	if len(instance.Inputs) != expectedInputs {
		t.Errorf("Expected %d inputs, got %d", expectedInputs, len(instance.Inputs))
	}
	if len(instance.Outputs) != expectedOutputs {
		t.Errorf("Expected %d outputs, got %d", expectedOutputs, len(instance.Outputs))
	}
	
	// Check memory was created
	if len(instance.Memories) != 1 {
		t.Errorf("Expected 1 memory, got %d", len(instance.Memories))
	}
}

// === INTEGRATION TESTS ===

func TestMultiClockDomainSystem(t *testing.T) {
	// Create a system with multiple clock domains
	m := &Module{
		Name:         "MultiClockSystem",
		Inputs:       []*Signal{},
		Outputs:      []*Signal{},
		Wires:        []*Signal{},
		Regs:         []*Signal{},
		Always:       []string{},
		ClockDomains: []*ClockDomain{},
	}
	
	// Create clock domains
	cpuClk := m.Input("cpu_clk", 1)
	cpuRst := m.Input("cpu_rst", 1)
	ddrClk := m.Input("ddr_clk", 1)
	ddrRst := m.Input("ddr_rst", 1)
	
	cpuDomain := m.NewClockDomain("cpu_domain", cpuClk, cpuRst).SetFrequency(100000000)
	ddrDomain := m.NewClockDomain("ddr_domain", ddrClk, ddrRst).SetFrequency(200000000)
	
	// Create signals in different domains
	cpuData := m.Wire("cpu_data", 32).WithClockDomain(cpuDomain)
	ddrData := m.Wire("ddr_data", 32).WithClockDomain(ddrDomain)
	
	// Create CDC synchronizer
	syncData := m.CDCSynchronizer("cpu_to_ddr", cpuData, ddrDomain, 3)
	
	// Use ddrData to avoid unused variable
	_ = ddrData
	
	// Verify the system
	if len(m.ClockDomains) != 2 {
		t.Errorf("Expected 2 clock domains, got %d", len(m.ClockDomains))
	}
	
	if cpuDomain.Frequency != 100000000 {
		t.Errorf("CPU domain frequency mismatch")
	}
	if ddrDomain.Frequency != 200000000 {
		t.Errorf("DDR domain frequency mismatch")
	}
	
	if syncData.Width != cpuData.Width {
		t.Errorf("Synchronized data width mismatch")
	}
}

func TestPolymorphicSystemWithMutex(t *testing.T) {
	// Create a system using polymorphic templates and hardware mutexes
	m := &Module{
		Name:    "PolymorphicSystem",
		Inputs:  []*Signal{},
		Outputs: []*Signal{},
		Mutexes: []*Mutex{},
	}
	
	// Create clock and reset
	clk := m.Input("clk", 1)
	rst := m.Input("rst", 1)
	
	// Use rst to avoid unused variable
	_ = rst
	
	// Create a mutex for resource arbitration
	busMutex := m.Mutex("bus_arbiter", 3, "round_robin")
	busMutex.GenerateRoundRobin(m, clk)
	
	// Create polymorphic FIFO instances
	fifoTemplate := FIFOTemplate()
	
	fifo1 := m.InstantiateTemplate(fifoTemplate, "fifo_8x32", map[string]interface{}{
		"DATA_WIDTH": Width(8),
		"DEPTH":      32,
	})
	
	fifo2 := m.InstantiateTemplate(fifoTemplate, "fifo_16x64", map[string]interface{}{
		"DATA_WIDTH": Width(16),
		"DEPTH":      64,
	})
	
	// Verify system
	if len(m.Mutexes) != 1 {
		t.Errorf("Expected 1 mutex, got %d", len(m.Mutexes))
	}
	
	if fifo1.Name != "fifo_8x32" {
		t.Errorf("FIFO1 name mismatch: %s", fifo1.Name)
	}
	if fifo2.Name != "fifo_16x64" {
		t.Errorf("FIFO2 name mismatch: %s", fifo2.Name)
	}
	
	// Check different data widths were instantiated correctly
	fifo1DataWidth := fifo1.Parameters["DATA_WIDTH"].(Width)
	fifo2DataWidth := fifo2.Parameters["DATA_WIDTH"].(Width)
	
	if fifo1DataWidth != 8 {
		t.Errorf("FIFO1 data width mismatch: %d", fifo1DataWidth)
	}
	if fifo2DataWidth != 16 {
		t.Errorf("FIFO2 data width mismatch: %d", fifo2DataWidth)
	}
}