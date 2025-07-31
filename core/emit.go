// emit.go
package core

import (
	"fmt"
	"github.com/SoulPancake/HFT/types"
	"os"
)

// emit.go
func EmitVerilog(mod *hdl.Module) {
	f, err := os.Create("out.v")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fmt.Fprintf(f, "module %s(\n", mod.Name)
	for i, sig := range append(mod.Inputs, mod.Outputs...) {
		comma := ","
		if i == len(mod.Inputs)+len(mod.Outputs)-1 {
			comma = ""
		}
		fmt.Fprintf(f, "  %s [%d:0] %s%s\n", sig.Kind, sig.Width-1, sig.Name, comma)
	}
	fmt.Fprintf(f, ");\n")

	for _, wire := range mod.Wires {
		fmt.Fprintf(f, "  wire [%d:0] %s;\n", wire.Width-1, wire.Name)
	}
	for _, assign := range mod.Assigns {
		fmt.Fprintf(f, "  %s\n", assign)
	}

	fmt.Fprintf(f, "endmodule\n")
}
