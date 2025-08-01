package main

import (
	"fmt"
	"os"
	"github.com/SoulPancake/HFT/core"
	hdl "github.com/SoulPancake/HFT/types"
)

func main() {
	fmt.Println("HFT Advanced Features Demo")
	fmt.Println("==========================")
	
	// Create a comprehensive demo showing multiple clock domains, mutexes, and polymorphism
	m := core.NewModule("AdvancedHFTDemo")
	
	// === Multiple Clock Domains ===
	fmt.Println("Creating multiple clock domains...")
	cpuClk := m.Input("cpu_clk", 1)
	cpuRst := m.Input("cpu_rst", 1)
	memClk := m.Input("mem_clk", 1)
	memRst := m.Input("mem_rst", 1)
	
	cpuDomain := m.NewClockDomain("cpu_domain", cpuClk, cpuRst).SetFrequency(100000000)
	memDomain := m.NewClockDomain("mem_domain", memClk, memRst).SetFrequency(200000000)
	
	// === Hardware Mutexes ===
	fmt.Println("Creating hardware mutexes...")
	busMutex := m.Mutex("bus_arbiter", 3, "priority")
	busMutex.GeneratePriority(m)
	
	memMutex := m.Mutex("memory_arbiter", 2, "round_robin") 
	memMutex.GenerateRoundRobin(m, cpuClk)
	
	// === Polymorphic Components ===
	fmt.Println("Creating polymorphic components...")
	
	// Create FIFO instances with different sizes
	smallFIFO := m.InstantiateTemplate(hdl.FIFOTemplate(), "small_fifo", map[string]interface{}{
		"DATA_WIDTH": hdl.Width(8),
		"DEPTH":      16,
	})
	
	largeFIFO := m.InstantiateTemplate(hdl.FIFOTemplate(), "large_fifo", map[string]interface{}{
		"DATA_WIDTH": hdl.Width(32),
		"DEPTH":      128,
	})
	
	// Create register file
	regFile := m.InstantiateTemplate(hdl.RegisterFileTemplate(), "main_registers", map[string]interface{}{
		"DATA_WIDTH": hdl.Width(32),
		"ADDR_WIDTH": hdl.Width(5),
	})
	
	// === Clock Domain Crossing ===
	fmt.Println("Creating clock domain crossing logic...")
	cpuData := m.Wire("cpu_data", 32).WithClockDomain(cpuDomain)
	memDataSync := m.CDCSynchronizer("cpu_to_mem", cpuData, memDomain, 2)
	
	// Create async FIFO for high throughput CDC
	wrData, rdData, wrFull, rdEmpty := m.AsyncFIFO("cpu_mem_fifo", 32, 64, cpuDomain, memDomain)
	
	// === System Connections ===
	fmt.Println("Creating system connections...")
	
	// Inputs and outputs
	dataIn := m.Input("data_in", 32)
	addrIn := m.Input("addr_in", 5)
	enableIn := m.Input("enable_in", 1)
	
	dataOut := m.Output("data_out", 32)
	statusOut := m.Output("status_out", 8)
	
	// Connect data flow
	m.Assign(cpuData, dataIn)
	m.Assign(dataOut, memDataSync)
	
	// Create status word from various signals
	statusWord := hdl.Cat(
		busMutex.Grants[0],    // Bus arbiter grant 0
		memMutex.Grants[0],    // Memory arbiter grant 0  
		wrFull,                // FIFO write full
		rdEmpty,               // FIFO read empty
		enableIn,              // Enable input
		hdl.Fill(3, "0"),      // Padding
	)
	m.Assign(statusOut, statusWord)
	
	// Use all variables to avoid warnings
	_ = smallFIFO
	_ = largeFIFO
	_ = regFile
	_ = wrData
	_ = rdData
	_ = addrIn
	
	// === Generate Verilog ===
	fmt.Println("Generating Verilog output...")
	core.EmitVerilog(m)
	
	// Rename the output file to show what it contains
	if err := os.Rename("out.v", "advanced_demo.v"); err == nil {
		fmt.Println("Generated advanced_demo.v")
	}
	
	fmt.Println("\n=== Demo Complete ===")
	fmt.Println("Features demonstrated:")
	fmt.Printf("✓ %d Clock domains with frequency specification\n", len(m.ClockDomains))
	fmt.Printf("✓ %d Hardware mutexes with different arbitration schemes\n", len(m.Mutexes))
	fmt.Println("✓ Polymorphic FIFO and register file templates")
	fmt.Println("✓ Clock domain crossing with synchronizers and async FIFO")
	fmt.Println("✓ Type-safe signal operations with automatic width inference")
	fmt.Println("✓ Production-quality Verilog generation")
}