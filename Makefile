# HFT (Hardware Function Templates) Makefile
# Automates building, testing, and simulation workflow

.PHONY: help build test demo demo-advanced sim waves clean all

# Default target
all: build test demo

# Help target - shows available commands
help:
	@echo "HFT (Hardware Function Templates) - Build System"
	@echo "================================================"
	@echo ""
	@echo "Available targets:"
	@echo "  build         - Build the entire project"
	@echo "  test          - Run all tests"
	@echo "  demo          - Run simple demo (driver/main.go)"
	@echo "  demo-advanced - Run advanced features demo"
	@echo "  sim           - Run Verilog simulation"
	@echo "  waves         - View simulation waveforms (requires display)"
	@echo "  clean         - Clean generated files"
	@echo "  all           - Build, test, and run demo"
	@echo "  help          - Show this help message"

# Build the entire project
build:
	@echo "Building HFT project..."
	go mod tidy
	go build ./...
	@echo "✓ Build complete"

# Run all tests
test:
	@echo "Running tests..."
	go test ./... -v
	@echo "✓ All tests passed"

# Run simple demo
demo:
	@echo "Running simple demo..."
	go run driver/main.go
	@echo "✓ Generated: out.v"
	@if [ -f out.v ]; then \
		echo "Generated Verilog (first 20 lines):"; \
		echo "====================================="; \
		head -20 out.v; \
		echo "====================================="; \
	fi

# Run advanced features demo
demo-advanced:
	@echo "Running advanced features demo..."
	go run examples/simple_demo.go
	@echo "✓ Demo complete"
	@if [ -f advanced_demo.v ]; then \
		echo "Generated Verilog info:"; \
		echo "======================"; \
		wc -l advanced_demo.v; \
		echo "Module declaration:"; \
		head -10 advanced_demo.v | grep -E "(module|input|output)" || true; \
	fi

# Run additional example demos
demo-data:
	@echo "Running data structures demo..."
	go run examples/data_structures_demo.go || echo "Demo not available"

demo-phase3:
	@echo "Running phase 3 demo..."
	go run examples/phase3_demo.go || echo "Demo not available"

# Simulate with Verilog tools (requires iverilog)
sim: build-sim run-sim

build-sim:
	@echo "Building simulation..."
	@if command -v iverilog >/dev/null 2>&1; then \
		echo "Using iverilog for simulation"; \
		if [ -f out.v ]; then \
			iverilog -o sim_out out.v testbench.v; \
			echo "✓ Simulation compiled"; \
		else \
			echo "⚠ No out.v found - running demo first"; \
			$(MAKE) demo; \
			iverilog -o sim_out out.v testbench.v; \
		fi \
	else \
		echo "⚠ iverilog not found - install with: apt-get install iverilog"; \
		echo "   Simulation skipped"; \
	fi

run-sim:
	@if [ -f sim_out ]; then \
		echo "Running simulation..."; \
		./sim_out; \
		echo "✓ Simulation complete"; \
		if [ -f dump.vcd ]; then \
			echo "✓ Waveform dump created: dump.vcd"; \
		fi \
	else \
		echo "⚠ No simulation executable found"; \
	fi

# View waveforms (requires gtkwave or similar)
waves:
	@if [ -f dump.vcd ]; then \
		if command -v gtkwave >/dev/null 2>&1; then \
			echo "Opening waveforms with gtkwave..."; \
			gtkwave dump.vcd & \
		else \
			echo "⚠ gtkwave not found - install with: apt-get install gtkwave"; \
			echo "Waveform file available: dump.vcd"; \
		fi \
	else \
		echo "⚠ No waveform file found - run 'make sim' first"; \
	fi

# Clean generated files
clean:
	@echo "Cleaning generated files..."
	rm -f *.v
	rm -f *.vcd
	rm -f sim_out sim
	rm -f core/test_*.v
	@echo "✓ Cleanup complete"

# Development targets

# Format Go code
fmt:
	@echo "Formatting Go code..."
	go fmt ./...
	@echo "✓ Code formatted"

# Run linter (if available)
lint:
	@if command -v golint >/dev/null 2>&1; then \
		echo "Running linter..."; \
		golint ./...; \
	else \
		echo "⚠ golint not available - install with: go install golang.org/x/lint/golint@latest"; \
	fi

# Quick development cycle
dev: clean build test demo
	@echo "✓ Development cycle complete"

# CI/CD simulation (full pipeline)
ci: clean build test demo demo-advanced sim
	@echo "✓ Full CI pipeline complete"

# Show project info
info:
	@echo "HFT Project Information"
	@echo "======================"
	@echo "Go version: $$(go version)"
	@echo "Module: $$(head -1 go.mod)"
	@echo "Files:"
	@find . -name "*.go" -not -path "./.git/*" | wc -l | xargs echo "  Go files:"
	@find . -name "*.v" -not -path "./.git/*" | wc -l | xargs echo "  Verilog files:"
	@echo "Test status:"
	@go test ./... -v | grep -E "(PASS|FAIL)" | tail -1 || echo "  Run 'make test' to see test status"