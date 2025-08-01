package core

import (
	"fmt"
	"github.com/SoulPancake/HFT/types"
	"os"
)

func EmitVerilog(mod *hdl.Module) {
	f, err := os.Create("out.v")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Emit module header with parameters
	if len(mod.Parameters) > 0 {
		fmt.Fprintf(f, "module %s #(\n", mod.Name)
		paramKeys := make([]string, 0, len(mod.Parameters))
		for k := range mod.Parameters {
			paramKeys = append(paramKeys, k)
		}
		for i, key := range paramKeys {
			comma := ","
			if i == len(paramKeys)-1 {
				comma = ""
			}
			fmt.Fprintf(f, "  parameter %s = %v%s\n", key, mod.Parameters[key], comma)
		}
		fmt.Fprintf(f, ") (\n")
	} else {
		fmt.Fprintf(f, "module %s(\n", mod.Name)
	}
	
	// Collect all ports
	allPorts := append([]*hdl.Signal{}, mod.Inputs...)
	allPorts = append(allPorts, mod.Outputs...)
	
	for i, sig := range allPorts {
		comma := ","
		if i == len(allPorts)-1 {
			comma = ""
		}
		fmt.Fprintf(f, "  %s %s %s%s\n", sig.Kind, sig.Width.Bits(), sig.Name, comma)
	}
	fmt.Fprintf(f, ");\n\n")

	// Emit wire declarations
	for _, wire := range mod.Wires {
		fmt.Fprintf(f, "  wire %s %s;\n", wire.Width.Bits(), wire.Name)
	}
	
	// Emit reg declarations
	for _, reg := range mod.Regs {
		fmt.Fprintf(f, "  reg %s %s;\n", reg.Width.Bits(), reg.Name)
	}
	
	// Emit bundle field declarations
	for _, bundle := range mod.Bundles {
		for _, field := range bundle.Fields {
			fmt.Fprintf(f, "  %s %s %s;\n", field.Kind, field.Width.Bits(), field.Name)
		}
	}
	
	// Emit Vec element declarations
	for _, vec := range mod.Vecs {
		for _, element := range vec.Elements {
			fmt.Fprintf(f, "  %s %s %s;\n", element.Kind, element.Width.Bits(), element.Name)
		}
	}
	
	// Emit memory declarations
	for _, mem := range mod.Memories {
		fmt.Fprintf(f, "  reg %s %s [0:%d];\n", mem.Width.Bits(), mem.Name, mem.Depth-1)
	}
	
	if len(mod.Wires) > 0 || len(mod.Regs) > 0 || len(mod.Bundles) > 0 || len(mod.Vecs) > 0 || len(mod.Memories) > 0 {
		fmt.Fprintf(f, "\n")
	}

	// Emit assign statements
	for _, assign := range mod.Assigns {
		fmt.Fprintf(f, "  %s\n", assign)
	}
	
	if len(mod.Assigns) > 0 {
		fmt.Fprintf(f, "\n")
	}

	// Emit module instances
	for _, inst := range mod.Instances {
		if len(inst.Parameters) > 0 {
			fmt.Fprintf(f, "  %s #(\n", inst.ModuleName)
			paramKeys := make([]string, 0, len(inst.Parameters))
			for k := range inst.Parameters {
				paramKeys = append(paramKeys, k)
			}
			for i, key := range paramKeys {
				comma := ","
				if i == len(paramKeys)-1 {
					comma = ""
				}
				fmt.Fprintf(f, "    .%s(%v)%s\n", key, inst.Parameters[key], comma)
			}
			fmt.Fprintf(f, "  ) %s (\n", inst.InstanceName)
		} else {
			fmt.Fprintf(f, "  %s %s (\n", inst.ModuleName, inst.InstanceName)
		}
		
		connKeys := make([]string, 0, len(inst.Connections))
		for k := range inst.Connections {
			connKeys = append(connKeys, k)
		}
		for i, port := range connKeys {
			comma := ","
			if i == len(connKeys)-1 {
				comma = ""
			}
			signal := inst.Connections[port]
			fmt.Fprintf(f, "    .%s(%s)%s\n", port, signal.Name, comma)
		}
		fmt.Fprintf(f, "  );\n\n")
	}

	// Emit clock domain information as comments
	if len(mod.ClockDomains) > 0 {
		fmt.Fprintf(f, "  // Clock Domains:\n")
		for _, cd := range mod.ClockDomains {
			fmt.Fprintf(f, "  // - %s: clk=%s, rst=%s", cd.Name, cd.Clock.Name, cd.Reset.Name)
			if cd.Frequency > 0 {
				fmt.Fprintf(f, ", freq=%d Hz", cd.Frequency)
			}
			fmt.Fprintf(f, "\n")
		}
		fmt.Fprintf(f, "\n")
	}

	// Emit mutex arbiters
	for _, mutex := range mod.Mutexes {
		fmt.Fprintf(f, "  // Mutex: %s (%s arbitration)\n", mutex.Name, mutex.Arbitration)
		for i := 0; i < len(mutex.Requests); i++ {
			fmt.Fprintf(f, "  // Request[%d]: %s -> Grant[%d]: %s\n", 
				i, mutex.Requests[i].Name, i, mutex.Grants[i].Name)
		}
		fmt.Fprintf(f, "\n")
	}

	// Emit polymorphic template info as comments
	if len(mod.Templates) > 0 {
		fmt.Fprintf(f, "  // Polymorphic Templates:\n")
		for _, template := range mod.Templates {
			fmt.Fprintf(f, "  // - %s with type parameters: %v\n", template.Name, template.TypeParams)
		}
		fmt.Fprintf(f, "\n")
	}

	// Emit always blocks
	for _, always := range mod.Always {
		fmt.Fprintf(f, "  %s\n\n", always)
	}

	fmt.Fprintf(f, "endmodule\n")
}
