package main

import (
	"fmt"
	"github.com/SoulPancake/HFT/core"
	hdl "github.com/SoulPancake/HFT/types"
)

// Demonstrates multiple clock domains, CDC, and async FIFO
func createMultiClockDomainSystem() *hdl.Module {
	fmt.Println("=== Multi-Clock Domain System ===")
	
	m := core.NewModule("MultiClockSystem")
	
	// Create different clock domains
	cpuClk := m.Input("cpu_clk", 1)
	cpuRst := m.Input("cpu_rst_n", 1)
	ddrClk := m.Input("ddr_clk", 1) 
	ddrRst := m.Input("ddr_rst_n", 1)
	ethClk := m.Input("eth_clk", 1)
	ethRst := m.Input("eth_rst_n", 1)
	
	// Create clock domains with different frequencies
	cpuDomain := m.NewClockDomain("cpu_domain", cpuClk, cpuRst).SetFrequency(100000000) // 100 MHz
	ddrDomain := m.NewClockDomain("ddr_domain", ddrClk, ddrRst).SetFrequency(200000000) // 200 MHz  
	ethDomain := m.NewClockDomain("eth_domain", ethClk, ethRst).SetFrequency(125000000) // 125 MHz
	
	fmt.Printf("Created %d clock domains\n", len(m.ClockDomains))
	
	// Create signals in different domains
	cpuData := m.Wire("cpu_data", 32).WithClockDomain(cpuDomain)
	cpuValid := m.Wire("cpu_valid", 1).WithClockDomain(cpuDomain)
	
	ddrAddr := m.Wire("ddr_addr", 24).WithClockDomain(ddrDomain)
	ddrData := m.Wire("ddr_data", 64).WithClockDomain(ddrDomain)
	
	ethPacket := m.Wire("eth_packet", 8).WithClockDomain(ethDomain)
	ethValid := m.Wire("eth_valid", 1).WithClockDomain(ethDomain)
	
	// Clock Domain Crossing - CPU to DDR
	fmt.Println("Creating CDC synchronizers...")
	cpuToDdrData := m.CDCSynchronizer("cpu_to_ddr_data", cpuData, ddrDomain, 2)
	cpuToDdrValid := m.CDCSynchronizer("cpu_to_ddr_valid", cpuValid, ddrDomain, 3)
	
	// Clock Domain Crossing - ETH to CPU
	ethToCpuPacket := m.CDCSynchronizer("eth_to_cpu_packet", ethPacket, cpuDomain, 2)
	ethToCpuValid := m.CDCSynchronizer("eth_to_cpu_valid", ethValid, cpuDomain, 2)
	
	// Asynchronous FIFO for high-throughput CDC
	fmt.Println("Creating Async FIFO for CPU-DDR interface...")
	fifoWrData, fifoRdData, fifoWrFull, fifoRdEmpty := m.AsyncFIFO("cpu_ddr_fifo", 32, 64, cpuDomain, ddrDomain)
	
	// Connect FIFO to synchronized signals
	m.Assign(fifoWrData, cpuToDdrData)
	
	// Create output assignments to demonstrate connectivity
	cpuDataOut := m.Output("cpu_data_out", 32)
	ddrDataOut := m.Output("ddr_data_out", 64)
	ethDataOut := m.Output("eth_data_out", 8)
	
	// Use unused variables in signal concatenation
	paddedData := ethToCpuPacket.Cat(hdl.Fill(24, "0"))
	combinedData := fifoRdData.Cat(cpuToDdrData)
	
	m.Assign(cpuDataOut, paddedData)
	m.Assign(ddrDataOut, combinedData)
	m.Assign(ethDataOut, ethPacket)
	
	// Use unused variables to avoid warnings
	_ = ddrAddr
	_ = ddrData
	_ = cpuToDdrValid
	_ = ethToCpuValid
	_ = fifoWrFull
	_ = fifoRdEmpty
	
	return m
}

// Demonstrates hardware mutex with different arbitration schemes
func createMutexArbitrationSystem() *hdl.Module {
	fmt.Println("\n=== Hardware Mutex Arbitration System ===")
	
	m := core.NewModule("MutexArbitrationSystem")
	
	// System clock and reset
	clk := m.Input("clk", 1)
	rst := m.Input("rst", 1)
	
	// Use rst to avoid unused variable warning
	_ = rst
	
	// Create multiple requesters for different resources
	
	// Memory arbiter with round-robin arbitration
	fmt.Println("Creating memory arbiter with round-robin...")
	memoryMutex := m.Mutex("memory_arbiter", 4, "round_robin")
	memoryMutex.GenerateRoundRobin(m, clk)
	
	// Bus arbiter with priority arbitration  
	fmt.Println("Creating bus arbiter with priority...")
	busMutex := m.Mutex("bus_arbiter", 3, "priority")
	busMutex.GeneratePriority(m)
	
	// Cache arbiter with round-robin
	fmt.Println("Creating cache arbiter with round-robin...")
	cacheMutex := m.Mutex("cache_arbiter", 2, "round_robin")
	cacheMutex.GenerateRoundRobin(m, clk)
	
	// Create some logic to demonstrate usage
	memAddr := m.Wire("mem_addr", 32)
	busAddr := m.Wire("bus_addr", 16)
	cacheAddr := m.Wire("cache_addr", 20)
	
	// Memory access logic
	memAccess := m.Wire("mem_access", 1)
	m.AssignExpr(memAccess, memoryMutex.Grants[0].Name + " || " + memoryMutex.Grants[1].Name)
	
	// Bus access logic
	busAccess := m.Wire("bus_access", 1) 
	m.AssignExpr(busAccess, busMutex.Grants[0].Name + " || " + busMutex.Grants[2].Name)
	
	// Cache access logic
	cacheAccess := m.Wire("cache_access", 1)
	m.AssignExpr(cacheAccess, cacheMutex.Grants[0].Name + " && " + cacheMutex.Grants[1].Name)
	
	// Use unused variables to avoid warnings
	_ = memAddr
	_ = busAddr  
	_ = cacheAddr
	_ = memAccess
	_ = busAccess
	_ = cacheAccess
	
	// Create outputs to show arbitration results
	memGrantsOut := m.Output("memory_grants", 4)
	busGrantsOut := m.Output("bus_grants", 3)
	cacheGrantsOut := m.Output("cache_grants", 2)
	
	// Concatenate grants for output
	m.Assign(memGrantsOut, hdl.Cat(memoryMutex.Grants...))
	m.Assign(busGrantsOut, hdl.Cat(busMutex.Grants...))
	m.Assign(cacheGrantsOut, hdl.Cat(cacheMutex.Grants...))
	
	fmt.Printf("Created %d mutex arbiters\n", len(m.Mutexes))
	
	return m
}

// Demonstrates polymorphic module templates and generic components
func createPolymorphicSystem() *hdl.Module {
	fmt.Println("\n=== Polymorphic Component System ===")
	
	m := core.NewModule("PolymorphicSystem")
	
	// System signals
	clk := m.Input("clk", 1)
	rst := m.Input("rst", 1)
	
	// Create polymorphic FIFO instances with different parameters
	fmt.Println("Creating polymorphic FIFO instances...")
	fifoTemplate := hdl.FIFOTemplate()
	
	// Small 8-bit FIFO for control data
	controlFIFO := m.InstantiateTemplate(fifoTemplate, "control_fifo", map[string]interface{}{
		"DATA_WIDTH": hdl.Width(8),
		"DEPTH":      16,
	})
	
	// Large 32-bit FIFO for data path
	dataFIFO := m.InstantiateTemplate(fifoTemplate, "data_fifo", map[string]interface{}{
		"DATA_WIDTH": hdl.Width(32),
		"DEPTH":      256,
	})
	
	// High-speed 64-bit FIFO for DMA
	dmaFIFO := m.InstantiateTemplate(fifoTemplate, "dma_fifo", map[string]interface{}{
		"DATA_WIDTH": hdl.Width(64),
		"DEPTH":      128,
	})
	
	fmt.Printf("Created %d FIFO instances\n", 3)
	
	// Create polymorphic register file instances
	fmt.Println("Creating polymorphic register file instances...")
	regFileTemplate := hdl.RegisterFileTemplate()
	
	// Small register file for control registers
	controlRegs := m.InstantiateTemplate(regFileTemplate, "control_registers", map[string]interface{}{
		"DATA_WIDTH": hdl.Width(16),
		"ADDR_WIDTH": hdl.Width(4), // 16 registers
	})
	
	// Large register file for general purpose registers
	gpRegs := m.InstantiateTemplate(regFileTemplate, "gp_registers", map[string]interface{}{
		"DATA_WIDTH": hdl.Width(32),
		"ADDR_WIDTH": hdl.Width(5), // 32 registers
	})
	
	fmt.Printf("Created %d register file instances\n", 2)
	
	// Create a custom polymorphic module template
	fmt.Println("Creating custom ALU template...")
	aluTemplate := createALUTemplate()
	
	// Instantiate different ALU sizes
	alu8 := m.InstantiateTemplate(aluTemplate, "alu_8bit", map[string]interface{}{
		"WIDTH": hdl.Width(8),
	})
	
	alu16 := m.InstantiateTemplate(aluTemplate, "alu_16bit", map[string]interface{}{
		"WIDTH": hdl.Width(16),
	})
	
	alu32 := m.InstantiateTemplate(aluTemplate, "alu_32bit", map[string]interface{}{
		"WIDTH": hdl.Width(32),
	})
	
	fmt.Printf("Created %d ALU instances\n", 3)
	
	// Connect some signals to demonstrate the system
	dataInput := m.Input("data_in", 32)
	controlInput := m.Input("control_in", 8)
	
	// Simple data flow demonstration
	dataOut := m.Output("data_out", 32)
	controlOut := m.Output("control_out", 8)
	
	// Note: In a real implementation, you would connect these properly
	// This is simplified for demonstration
	m.AssignExpr(dataOut, dataInput.Name)
	m.AssignExpr(controlOut, controlInput.Name)
	
	// Use unused variables to avoid warnings
	_ = clk
	_ = rst
	_ = controlFIFO
	_ = dataFIFO
	_ = dmaFIFO
	_ = controlRegs
	_ = gpRegs
	_ = alu8
	_ = alu16
	_ = alu32
	
	return m
}

// Helper function to create a custom ALU template
func createALUTemplate() *hdl.ModuleTemplate {
	template := hdl.NewModuleTemplate("GenericALU", []string{"WIDTH"})
	template.AddConstraint("WIDTH", hdl.Width(0))
	
	template.SetGenerator(func(params map[string]interface{}) *hdl.Module {
		width := params["WIDTH"].(hdl.Width)
		
		m := &hdl.Module{
			Name:       "GenericALU",
			Inputs:     []*hdl.Signal{},
			Outputs:    []*hdl.Signal{},
			Wires:      []*hdl.Signal{},
			Assigns:    []string{},
			Parameters: make(map[string]interface{}),
		}
		
		// ALU interface
		a := m.Input("a", width)
		b := m.Input("b", width)
		op := m.Input("op", 4) // 4-bit operation code
		result := m.Output("result", width)
		zero := m.Output("zero", 1)
		overflow := m.Output("overflow", 1)
		
		// ALU operations (simplified)
		addResult := a.Add(b)
		subResult := a.Sub(b)
		andResult := a.And(b)
		orResult := a.Or(b)
		xorResult := a.Xor(b)
		
		// Multiplexer for operation selection
		aluOut := hdl.Mux(op, addResult, subResult, andResult, orResult, xorResult,
			hdl.Fill(int(width), "0"), hdl.Fill(int(width), "0"), hdl.Fill(int(width), "0"))
		
		m.Assign(result, aluOut)
		
		// Zero flag
		zeroFlag := aluOut.Eq(hdl.Lit(0, width))
		m.Assign(zero, zeroFlag)
		
		// Overflow detection (simplified)
		overflowWire := m.Wire("overflow_detect", 1)
		m.AssignExpr(overflowWire, "1'b0") // Simplified
		m.Assign(overflow, overflowWire)
		
		// Set parameters
		m.SetParameter("WIDTH", width)
		
		return m
	})
	
	return template
}

// Demonstrates a complex system integrating all advanced features
func createComplexIntegratedSystem() *hdl.Module {
	fmt.Println("\n=== Complex Integrated System ===")
	
	m := core.NewModule("ComplexSystem")
	
	// === CLOCK DOMAINS ===
	cpuClk := m.Input("cpu_clk", 1)
	cpuRst := m.Input("cpu_rst", 1)
	ddrClk := m.Input("ddr_clk", 1)
	ddrRst := m.Input("ddr_rst", 1)
	netClk := m.Input("net_clk", 1)
	netRst := m.Input("net_rst", 1)
	
	cpuDomain := m.NewClockDomain("cpu", cpuClk, cpuRst).SetFrequency(100000000)
	ddrDomain := m.NewClockDomain("ddr", ddrClk, ddrRst).SetFrequency(200000000)
	netDomain := m.NewClockDomain("net", netClk, netRst).SetFrequency(125000000)
	
	// === HARDWARE MUTEXES ===
	busMutex := m.Mutex("system_bus", 4, "priority")
	busMutex.GeneratePriority(m)
	
	memMutex := m.Mutex("memory_controller", 3, "round_robin")
	memMutex.GenerateRoundRobin(m, cpuClk)
	
	// === POLYMORPHIC COMPONENTS ===
	
	// CPU-side components
	cpuFIFO := m.InstantiateTemplate(hdl.FIFOTemplate(), "cpu_cmd_fifo", map[string]interface{}{
		"DATA_WIDTH": hdl.Width(32),
		"DEPTH":      64,
	})
	
	cpuRegs := m.InstantiateTemplate(hdl.RegisterFileTemplate(), "cpu_registers", map[string]interface{}{
		"DATA_WIDTH": hdl.Width(32),
		"ADDR_WIDTH": hdl.Width(5),
	})
	
	// DDR-side components  
	ddrFIFO := m.InstantiateTemplate(hdl.FIFOTemplate(), "ddr_buffer", map[string]interface{}{
		"DATA_WIDTH": hdl.Width(64),
		"DEPTH":      128,
	})
	
	// Network-side components
	netFIFO := m.InstantiateTemplate(hdl.FIFOTemplate(), "net_packet_fifo", map[string]interface{}{
		"DATA_WIDTH": hdl.Width(8),
		"DEPTH":      512,
	})
	
	// === CLOCK DOMAIN CROSSINGS ===
	
	// CPU to DDR data path
	cpuData := m.Wire("cpu_data", 32).WithClockDomain(cpuDomain)
	cpuValid := m.Wire("cpu_valid", 1).WithClockDomain(cpuDomain)
	
	ddrDataSync := m.CDCSynchronizer("cpu_to_ddr_data", cpuData, ddrDomain, 2)
	ddrValidSync := m.CDCSynchronizer("cpu_to_ddr_valid", cpuValid, ddrDomain, 2)
	
	// Network to CPU data path
	netData := m.Wire("net_data", 8).WithClockDomain(netDomain)
	netValid := m.Wire("net_valid", 1).WithClockDomain(netDomain)
	
	cpuNetDataSync := m.CDCSynchronizer("net_to_cpu_data", netData, cpuDomain, 3)
	cpuNetValidSync := m.CDCSynchronizer("net_to_cpu_valid", netValid, cpuDomain, 3)
	
	// Async FIFO for high-bandwidth DDR-Network path
	netToDdrWrData, netToDdrRdData, netToDdrWrFull, netToDdrRdEmpty := 
		m.AsyncFIFO("net_to_ddr", 16, 256, netDomain, ddrDomain)
	
	// === SYSTEM INTERCONNECT ===
	
	// Create system outputs to demonstrate connectivity
	systemDataOut := m.Output("system_data_out", 64)
	systemStatusOut := m.Output("system_status_out", 8)
	
	// Complex data path combining multiple domains and components
	combinedData := ddrDataSync.Cat(cpuNetDataSync.Cat(hdl.Fill(24, "0")))
	m.Assign(systemDataOut, combinedData)
	
	// Status combining mutex grants and FIFO status
	statusBits := hdl.Cat(
		busMutex.Grants[0],
		memMutex.Grants[0], 
		netToDdrWrFull,
		netToDdrRdEmpty,
		ddrValidSync,
		cpuNetValidSync,
		hdl.Fill(2, "0"),
	)
	m.Assign(systemStatusOut, statusBits)
	
	// Use all variables to avoid unused warnings
	_ = cpuFIFO
	_ = cpuRegs
	_ = ddrFIFO
	_ = netFIFO
	_ = netToDdrWrData
	_ = netToDdrRdData
	
	fmt.Println("Complex system created with:")
	fmt.Printf("  - %d clock domains\n", len(m.ClockDomains))
	fmt.Printf("  - %d hardware mutexes\n", len(m.Mutexes))
	fmt.Printf("  - Multiple polymorphic component instances\n")
	fmt.Printf("  - Cross-domain synchronizers and async FIFOs\n")
	
	return m
}

// Main demonstration function
func runAdvancedFeaturesDemo() {
	fmt.Println("HFT Advanced Features Demonstration")
	fmt.Println("===================================")
	
	// Create and emit each demonstration module
	
	// Multi-clock domain system
	multiClockSys := createMultiClockDomainSystem()
	fmt.Println("Emitting MultiClockSystem...")
	core.EmitVerilog(multiClockSys)
	// Move output to specific file
	// os.Rename("out.v", "test_multiclock.v")
	
	// Mutex arbitration system
	mutexSys := createMutexArbitrationSystem()
	fmt.Println("Emitting MutexArbitrationSystem...")
	core.EmitVerilog(mutexSys)
	
	// Polymorphic system
	polySys := createPolymorphicSystem()
	fmt.Println("Emitting PolymorphicSystem...")
	core.EmitVerilog(polySys)
	
	// Complex integrated system
	complexSys := createComplexIntegratedSystem()
	fmt.Println("Emitting ComplexSystem...")
	core.EmitVerilog(complexSys)
	
	fmt.Println("\n=== Demo Complete ===")
	fmt.Println("Advanced features demonstrated:")
	fmt.Println("✓ Multiple clock domains with frequency specification")
	fmt.Println("✓ Clock domain crossing (CDC) synchronizers")
	fmt.Println("✓ Asynchronous FIFOs for high-throughput CDC")
	fmt.Println("✓ Hardware mutexes with round-robin and priority arbitration")
	fmt.Println("✓ Polymorphic module templates with type constraints")
	fmt.Println("✓ Generic FIFO and register file templates")
	fmt.Println("✓ Custom polymorphic module creation")
	fmt.Println("✓ Complex system integration combining all features")
}