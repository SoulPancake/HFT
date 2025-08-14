# HFT (Hardware Function Templates)

A modern, Go-based Hardware Description Language (HDL) that generates production-quality Verilog code with advanced features like multiple clock domains, hardware mutexes, polymorphic templates, and clock domain crossing logic.

## Features

✨ **Core Capabilities:**
- 🔧 **Type-safe signal operations** with automatic width inference
- 🕒 **Multiple clock domains** with frequency specification  
- 🔒 **Hardware mutexes** with priority and round-robin arbitration
- 📦 **Polymorphic templates** for reusable components (FIFO, register files, etc.)
- 🌉 **Clock domain crossing** with synchronizers and async FIFOs
- ⚡ **Production-quality Verilog generation**

## Quick Start

### Prerequisites

- Go 1.24 or later
- Optional: Verilog simulator (iverilog, ModelSim, etc.) for running simulations

### Installation

```bash
git clone https://github.com/SoulPancake/HFT.git
cd HFT
```

### Build and Run

```bash
# Build the project
make build

# Run a simple demo
make demo

# Run advanced features demo  
make demo-advanced

# Run simulation (requires Verilog simulator)
make sim

# Clean generated files
make clean
```

## Development Workflow

The HFT development workflow follows a clear path from high-level Go descriptions to simulation results:

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Go Code       │    │   HFT Framework │    │  Verilog Output │
│                 │    │                 │    │                 │
│ • Module()      │───▶│ • Type System   │───▶│ • module.v      │
│ • Input()       │    │ • Width Infer   │    │ • Standard      │
│ • Output()      │    │ • Clock Domains │    │   Verilog       │
│ • Assign()      │    │ • Templates     │    │ • Synthesizable │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                                              │
         │ go run                                       │
         ▼                                              ▼
┌─────────────────┐                          ┌─────────────────┐
│   Build & Test  │                          │   Simulation    │
│                 │                          │                 │
│ • make build    │                          │ • iverilog      │
│ • make test     │                          │ • vvp           │
│ • make demo     │                          │ • gtkwave       │
└─────────────────┘                          └─────────────────┘
         │                                              │
         └──────────────────┬───────────────────────────┘
                            ▼
                   ┌─────────────────┐
                   │    Results      │
                   │                 │
                   │ • VCD waveforms │
                   │ • Console output│
                   │ • Verification  │
                   └─────────────────┘
```

## Quick Flow Diagram

```
  HFT Go Code              Verilog Output           Simulation Results
┌─────────────────┐       ┌─────────────────┐       ┌─────────────────┐
│ module.go       │  ─▶   │ module.v        │  ─▶   │ waveforms.vcd   │
│ • NewModule()   │ make  │ • module decl   │ ivlog │ • timing        │
│ • Input()       │ demo  │ • assign stmts  │  +    │ • signal values │
│ • Output()      │  ─▶   │ • always blocks │ vvp   │ • test results  │
│ • Assign()      │       │ • wire decls    │  ─▶   │ • verification  │
│ • EmitVerilog() │       │ • standard RTL  │       │ • debugging     │
└─────────────────┘       └─────────────────┘       └─────────────────┘
       │                           │                           │
       ▼                           ▼                           ▼
  Go Compiler               Verilog Compiler             Waveform Viewer
  Type Safety               Syntax Check                 Visual Analysis
  Width Inference           Logic Synthesis              Timing Debug
```

### Workflow Steps

1. **Design**: Write hardware description in Go using HFT APIs
2. **Generate**: Run `make demo` to emit Verilog from Go code  
3. **Simulate**: Run `make sim` to compile and simulate with iverilog
4. **Analyze**: View results in console output and VCD waveforms

### 1. Design Hardware in Go

```go
package main

import (
    "github.com/SoulPancake/HFT/core"
    hdl "github.com/SoulPancake/HFT/types"
)

func main() {
    // Create a module
    m := core.NewModule("SimpleAdder")
    
    // Define inputs and outputs
    a := m.Input("a", 8)
    b := m.Input("b", 8) 
    sum := m.Output("sum", 8)
    
    // Describe logic with type-safe operations
    m.Assign(sum, a.Add(b))
    
    // Generate Verilog
    core.EmitVerilog(m)
}
```

### 2. Generate Verilog

```bash
go run examples/simple_demo.go
```

This produces `advanced_demo.v`:

```verilog
module AdvancedHFTDemo(
  input [0:0] cpu_clk,
  input [0:0] cpu_rst,
  input [0:0] mem_clk,
  input [0:0] mem_rst,
  // ... more ports
);
  // Generated Verilog implementation
endmodule
```

### 3. Simulate (Optional)

```bash
# Compile and run simulation
make sim

# View waveforms
make waves
```

## Examples

### Basic Adder
```bash
go run driver/main.go
```

### Advanced Features Demo
```bash
go run examples/simple_demo.go
```

Features demonstrated:
- Multiple clock domains (CPU @ 100MHz, Memory @ 200MHz)
- Hardware mutexes (priority and round-robin arbitration)
- Polymorphic FIFO and register file templates
- Clock domain crossing with synchronizers
- Async FIFO for high-throughput CDC

### More Examples

- `examples/data_structures_demo.go` - Data structure templates
- `examples/advanced_features_demo.go` - Complex system integration
- `examples/phase3_demo.go` - Phase 3 features

## Architecture

### Core Components

- **`core/`** - Module builder, Verilog emission, core HDL functionality
- **`types/`** - Type system, signal operations, width inference
- **`driver/`** - Simple driver program
- **`examples/`** - Example designs and demos

### Key APIs

#### Module Creation
```go
m := core.NewModule("ModuleName")
input := m.Input("signal_name", width)
output := m.Output("signal_name", width)  
wire := m.Wire("signal_name", width)
```

#### Clock Domains
```go
clk := m.Input("clk", 1)
rst := m.Input("rst", 1)
domain := m.NewClockDomain("cpu", clk, rst).SetFrequency(100000000)
signal := m.Wire("data", 32).WithClockDomain(domain)
```

#### Hardware Mutexes
```go
mutex := m.Mutex("arbiter", 4, "round_robin")
mutex.GenerateRoundRobin(m, clk)
// Access mutex.Requests[i] and mutex.Grants[i]
```

#### Polymorphic Templates
```go
fifo := m.InstantiateTemplate(hdl.FIFOTemplate(), "my_fifo", map[string]interface{}{
    "DATA_WIDTH": hdl.Width(32),
    "DEPTH":      64,
})
```

#### Clock Domain Crossing
```go
dataSync := m.CDCSynchronizer("sync", inputSignal, targetDomain, stages)
wrData, rdData, wrFull, rdEmpty := m.AsyncFIFO("fifo", 32, 64, domain1, domain2)
```

## Testing

Run the comprehensive test suite:

```bash
# Run all tests
make test

# Run specific test categories
go test ./core -v
go test ./types -v
```

## Makefile Targets

| Target | Description |
|--------|-------------|
| `build` | Build the entire project |
| `test` | Run all tests |
| `demo` | Run simple demo |
| `demo-advanced` | Run advanced features demo |
| `sim` | Run Verilog simulation |
| `waves` | View simulation waveforms |
| `clean` | Clean generated files |
| `help` | Show available targets |

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run `make test` to ensure all tests pass
6. Submit a pull request

## Project Structure

```
HFT/
├── README.md           # This file
├── Makefile           # Build automation
├── go.mod             # Go module definition
├── core/              # Core HDL functionality
│   ├── builder.go     # Module builder
│   ├── emit.go        # Verilog generation
│   └── *_test.go      # Tests
├── types/             # Type system and operations
│   ├── types.go       # Core types and operations
│   └── *_test.go      # Tests
├── examples/          # Example designs
│   ├── simple_demo.go
│   └── advanced_*.go
├── driver/            # Simple driver program
└── testbench.v        # Basic testbench
```

## License

[Add your license information here]

## Acknowledgments

Built with modern Go practices and inspired by Chisel/FIRRTL, Bluespec, and other advanced HDLs.