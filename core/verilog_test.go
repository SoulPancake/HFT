package core

import (
	"os"
	"strings"
	"testing"
	
	hdl "github.com/SoulPancake/HFT/types"
)

// Test Verilog generation for multiple clock domains
func TestVerilogClockDomains(t *testing.T) {
	m := NewModule("ClockDomainTest")
	
	// Create multiple clock domains
	clk1 := m.Input("clk1", 1)
	rst1 := m.Input("rst1", 1)
	clk2 := m.Input("clk2", 1)
	rst2 := m.Input("rst2", 1)
	
	domain1 := m.NewClockDomain("fast_domain", clk1, rst1).SetFrequency(200000000)
	domain2 := m.NewClockDomain("slow_domain", clk2, rst2).SetFrequency(50000000)
	
	// Create signals in different domains
	sig1 := m.Wire("sig1", 8).WithClockDomain(domain1)
	sig2 := m.Wire("sig2", 8).WithClockDomain(domain2)
	
	// Create CDC synchronizer
	syncSig := m.CDCSynchronizer("cross_sync", sig1, domain2, 2)
	
	// Generate Verilog
	EmitVerilog(m)
	
	// Read and verify output
	content, err := os.ReadFile("out.v")
	if err != nil {
		t.Fatalf("Failed to read generated Verilog: %v", err)
	}
	
	verilog := string(content)
	
	// Check clock domain comments
	if !strings.Contains(verilog, "Clock Domains:") {
		t.Error("Clock domain comments not found")
	}
	if !strings.Contains(verilog, "fast_domain: clk=clk1, rst=rst1, freq=200000000 Hz") {
		t.Error("Fast domain info not found")
	}
	if !strings.Contains(verilog, "slow_domain: clk=clk2, rst=rst2, freq=50000000 Hz") {
		t.Error("Slow domain info not found")
	}
	
	// Check CDC synchronizer logic
	if !strings.Contains(verilog, "cross_sync_sync_stage0") {
		t.Error("CDC synchronizer stage 0 not found")
	}
	if !strings.Contains(verilog, "cross_sync_sync_stage1") {
		t.Error("CDC synchronizer stage 1 not found")
	}
	
	// Use variables to avoid warnings
	_ = sig2
	_ = syncSig
	
	// Clean up
	os.Rename("out.v", "test_clock_domains.v")
}

// Test Verilog generation for hardware mutexes
func TestVerilogMutexes(t *testing.T) {
	m := NewModule("MutexTest")
	
	clk := m.Input("clk", 1)
	rst := m.Input("rst", 1)
	
	// Create priority mutex
	priorityMutex := m.Mutex("priority_arb", 3, "priority")
	priorityMutex.GeneratePriority(m)
	
	// Create round-robin mutex
	rrMutex := m.Mutex("rr_arb", 2, "round_robin")
	rrMutex.GenerateRoundRobin(m, clk)
	
	// Generate Verilog
	EmitVerilog(m)
	
	// Read and verify output
	content, err := os.ReadFile("out.v")
	if err != nil {
		t.Fatalf("Failed to read generated Verilog: %v", err)
	}
	
	verilog := string(content)
	
	// Check mutex comments
	if !strings.Contains(verilog, "Mutex: priority_arb (priority arbitration)") {
		t.Error("Priority mutex info not found")
	}
	if !strings.Contains(verilog, "Mutex: rr_arb (round_robin arbitration)") {
		t.Error("Round-robin mutex info not found")
	}
	
	// Check priority arbitration logic
	if !strings.Contains(verilog, "priority_arb_grant_2 = priority_arb_req_2") {
		t.Error("Priority arbitration logic not found")
	}
	
	// Check round-robin arbitration logic
	if !strings.Contains(verilog, "rr_arb_counter") {
		t.Error("Round-robin counter not found")
	}
	if !strings.Contains(verilog, "always @(posedge clk)") {
		t.Error("Round-robin always block not found")
	}
	
	// Use variables to avoid warnings
	_ = rst
	
	// Clean up
	os.Rename("out.v", "test_mutexes.v")
}

// Test Verilog generation for polymorphic components
func TestVerilogPolymorphicComponents(t *testing.T) {
	m := NewModule("PolymorphicTest")
	
	// Create polymorphic FIFO
	fifo8 := m.InstantiateTemplate(hdl.FIFOTemplate(), "fifo_8x16", map[string]interface{}{
		"DATA_WIDTH": hdl.Width(8),
		"DEPTH":      16,
	})
	
	fifo32 := m.InstantiateTemplate(hdl.FIFOTemplate(), "fifo_32x64", map[string]interface{}{
		"DATA_WIDTH": hdl.Width(32),
		"DEPTH":      64,
	})
	
	// Create polymorphic register file
	regFile := m.InstantiateTemplate(hdl.RegisterFileTemplate(), "registers", map[string]interface{}{
		"DATA_WIDTH": hdl.Width(16),
		"ADDR_WIDTH": hdl.Width(4),
	})
	
	// Use the instances to avoid warnings
	_ = fifo8
	_ = fifo32
	_ = regFile
	
	// Generate Verilog (would need actual module instantiation for full test)
	EmitVerilog(m)
	
	// Read and verify output
	content, err := os.ReadFile("out.v")
	if err != nil {
		t.Fatalf("Failed to read generated Verilog: %v", err)
	}
	
	verilog := string(content)
	
	// Basic check that module was generated
	if !strings.Contains(verilog, "module PolymorphicTest") {
		t.Error("Module declaration not found")
	}
	
	// Clean up
	os.Rename("out.v", "test_polymorphic.v")
}

// Test complex system integration
func TestVerilogComplexSystem(t *testing.T) {
	m := NewModule("ComplexSystem")
	
	// Multiple clock domains
	cpuClk := m.Input("cpu_clk", 1)
	cpuRst := m.Input("cpu_rst", 1)
	ddrClk := m.Input("ddr_clk", 1)
	ddrRst := m.Input("ddr_rst", 1)
	
	cpuDomain := m.NewClockDomain("cpu", cpuClk, cpuRst).SetFrequency(100000000)
	ddrDomain := m.NewClockDomain("ddr", ddrClk, ddrRst).SetFrequency(200000000)
	
	// Hardware mutexes
	busMutex := m.Mutex("bus", 4, "priority")
	busMutex.GeneratePriority(m)
	
	memMutex := m.Mutex("memory", 2, "round_robin")
	memMutex.GenerateRoundRobin(m, cpuClk)
	
	// Clock domain crossing
	cpuData := m.Wire("cpu_data", 32).WithClockDomain(cpuDomain)
	ddrDataSync := m.CDCSynchronizer("cpu_to_ddr", cpuData, ddrDomain, 3)
	
	// Async FIFO
	wrData, rdData, wrFull, rdEmpty := m.AsyncFIFO("main_fifo", 32, 128, cpuDomain, ddrDomain)
	
	// System outputs
	dataOut := m.Output("data_out", 32)
	statusOut := m.Output("status_out", 8)
	
	// Connections
	m.Assign(dataOut, ddrDataSync)
	statusBits := hdl.Cat(
		busMutex.Grants[0],
		memMutex.Grants[0],
		wrFull,
		rdEmpty,
		hdl.Fill(4, "0"),
	)
	m.Assign(statusOut, statusBits)
	
	// Use variables to avoid warnings
	_ = wrData
	_ = rdData
	
	// Generate Verilog
	EmitVerilog(m)
	
	// Read and verify output
	content, err := os.ReadFile("out.v")
	if err != nil {
		t.Fatalf("Failed to read generated Verilog: %v", err)
	}
	
	verilog := string(content)
	
	// Comprehensive checks
	checks := []string{
		"module ComplexSystem",
		"Clock Domains:",
		"cpu: clk=cpu_clk, rst=cpu_rst, freq=100000000 Hz",
		"ddr: clk=ddr_clk, rst=ddr_rst, freq=200000000 Hz",
		"Mutex: bus (priority arbitration)",
		"Mutex: memory (round_robin arbitration)",
		"cpu_to_ddr_sync_stage0",
		"cpu_to_ddr_sync_stage1", 
		"cpu_to_ddr_sync_stage2",
		"AsyncFIFO main_fifo implementation",
		"always @(posedge cpu_clk)",
		"always @(posedge ddr_clk)",
		"assign data_out =",
		"assign status_out =",
	}
	
	for _, check := range checks {
		if !strings.Contains(verilog, check) {
			t.Errorf("Missing required content: %s", check)
		}
	}
	
	// Check that we have the expected number of features
	clockDomainCount := strings.Count(verilog, "clk=")
	if clockDomainCount < 2 {
		t.Errorf("Expected at least 2 clock domains, found %d", clockDomainCount)
	}
	
	mutexCount := strings.Count(verilog, "Mutex:")
	if mutexCount < 2 {
		t.Errorf("Expected at least 2 mutexes, found %d", mutexCount)
	}
	
	// Clean up
	os.Rename("out.v", "test_complex_system.v")
}

// Test Verilog generation for async FIFO
func TestVerilogAsyncFIFO(t *testing.T) {
	m := NewModule("AsyncFIFOTest")
	
	// Create clock domains
	wrClk := m.Input("wr_clk", 1)
	wrRst := m.Input("wr_rst", 1)
	rdClk := m.Input("rd_clk", 1)
	rdRst := m.Input("rd_rst", 1)
	
	wrDomain := m.NewClockDomain("write", wrClk, wrRst)
	rdDomain := m.NewClockDomain("read", rdClk, rdRst)
	
	// Create async FIFO
	wrData, rdData, wrFull, rdEmpty := m.AsyncFIFO("test_fifo", 16, 32, wrDomain, rdDomain)
	
	// Create outputs to connect the FIFO signals
	wrDataOut := m.Output("wr_data_out", 16)
	rdDataOut := m.Output("rd_data_out", 16)
	wrFullOut := m.Output("wr_full_out", 1)
	rdEmptyOut := m.Output("rd_empty_out", 1)
	
	m.Assign(wrDataOut, wrData)
	m.Assign(rdDataOut, rdData)
	m.Assign(wrFullOut, wrFull)
	m.Assign(rdEmptyOut, rdEmpty)
	
	// Generate Verilog
	EmitVerilog(m)
	
	// Read and verify output
	content, err := os.ReadFile("out.v")
	if err != nil {
		t.Fatalf("Failed to read generated Verilog: %v", err)
	}
	
	verilog := string(content)
	
	// Check async FIFO components
	checks := []string{
		"AsyncFIFO test_fifo implementation",
		"test_fifo_mem",
		"test_fifo_wr_ptr",
		"test_fifo_rd_ptr",
		"test_fifo_wr_ptr_sync",
		"test_fifo_rd_ptr_sync",
		"always @(posedge wr_clk)",
		"always @(posedge rd_clk)",
	}
	
	for _, check := range checks {
		if !strings.Contains(verilog, check) {
			t.Errorf("Missing async FIFO component: %s", check)
		}
	}
	
	// Clean up
	os.Rename("out.v", "test_async_fifo.v")
}

// Test width inference in complex expressions
func TestVerilogWidthInference(t *testing.T) {
	m := NewModule("WidthInferenceTest")
	
	// Create signals with different widths
	a := m.Input("a", 8)
	b := m.Input("b", 16)
	c := m.Input("c", 4)
	
	// Create operations with width inference
	addResult := a.Add(b)           // Should be 16 bits (max width)
	mulResult := a.Mul(c)           // Should be 12 bits (8+4)
	catResult := a.Cat(b, c)        // Should be 28 bits (8+16+4)
	
	// Create outputs
	addOut := m.Output("add_out", addResult.Width)
	mulOut := m.Output("mul_out", mulResult.Width)
	catOut := m.Output("cat_out", catResult.Width)
	
	m.Assign(addOut, addResult)
	m.Assign(mulOut, mulResult)
	m.Assign(catOut, catResult)
	
	// Generate Verilog
	EmitVerilog(m)
	
	// Read and verify output
	content, err := os.ReadFile("out.v")
	if err != nil {
		t.Fatalf("Failed to read generated Verilog: %v", err)
	}
	
	verilog := string(content)
	
	// Check width declarations
	if !strings.Contains(verilog, "output [15:0] add_out") {
		t.Error("Add result width not correct (should be 16 bits)")
	}
	if !strings.Contains(verilog, "output [11:0] mul_out") {
		t.Error("Mul result width not correct (should be 12 bits)")
	}
	if !strings.Contains(verilog, "output [27:0] cat_out") {
		t.Error("Cat result width not correct (should be 28 bits)")
	}
	
	// Check assign statements
	if !strings.Contains(verilog, "assign add_out = a + b") {
		t.Error("Add assignment not found")
	}
	if !strings.Contains(verilog, "assign mul_out = a * c") {
		t.Error("Mul assignment not found")
	}
	if !strings.Contains(verilog, "assign cat_out = {a, b, c}") {
		t.Error("Cat assignment not found")
	}
	
	// Clean up
	os.Rename("out.v", "test_width_inference.v")
}