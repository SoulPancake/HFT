package hdl

import (
	"fmt"
	"strings"
)

type Width int

// Helper functions for width calculations
func WidthOf(value int) Width {
	if value == 0 {
		return 1
	}
	width := Width(0)
	for value > 0 {
		width++
		value >>= 1
	}
	return width
}

func MaxWidth(a, b Width) Width {
	if a > b {
		return a
	}
	return b
}

func (w Width) Bits() string {
	if w == 1 {
		return "[0:0]"
	}
	return fmt.Sprintf("[%d:0]", w-1)
}

type Signal struct {
	Name  string
	Width Width
	Kind  string // "input", "output", "wire", "reg"
	Expr  string // logic expression
	ClockDomain *ClockDomain // Associated clock domain
}

// Clock Domain for managing multiple clock domains
type ClockDomain struct {
	Name     string
	Clock    *Signal
	Reset    *Signal
	Frequency int // in Hz, optional for documentation
}

// Hardware Mutex for resource arbitration
type Mutex struct {
	Name     string
	Width    Width // Number of requesters
	Requests []*Signal
	Grants   []*Signal
	Arbitration string // "round_robin", "priority", "fifo"
}

// Generic/Polymorphic module template
type ModuleTemplate struct {
	Name       string
	TypeParams []string // Type parameter names
	Params     map[string]interface{} // Parameter constraints
	Generator  func(params map[string]interface{}) *Module
}

// Width-aware signal operations
func (s *Signal) WithWidth(width Width) *Signal {
	return &Signal{
		Name:  s.Name,
		Width: width,
		Kind:  s.Kind,
		Expr:  s.Expr,
	}
}

func (s *Signal) CheckWidth(other *Signal) error {
	if s.Width != other.Width {
		return fmt.Errorf("width mismatch: %s has width %d, %s has width %d", 
			s.Name, s.Width, other.Name, other.Width)
	}
	return nil
}

// Module instantiation support
type ModuleInstance struct {
	ModuleName   string
	InstanceName string
	Parameters   map[string]interface{}
	Connections  map[string]*Signal
}

// Bundle represents a collection of related signals
type Bundle struct {
	Name   string
	Fields map[string]*Signal
}

// Vec represents an array of signals
type Vec struct {
	Name     string
	Elements []*Signal
	Width    Width
	Size     int
}

// Memory represents a memory primitive
type Memory struct {
	Name      string
	Width     Width
	Depth     int
	AddrWidth Width
	Kind      string // "sync", "async"
}

type Module struct {
	Name       string
	Inputs     []*Signal
	Outputs    []*Signal
	Wires      []*Signal
	Regs       []*Signal
	Bundles    []*Bundle
	Vecs       []*Vec
	Memories   []*Memory
	Instances  []*ModuleInstance
	Assigns    []string
	Always     []string
	Clock      *Signal
	Reset      *Signal
	Parameters map[string]interface{}
	ClockDomains []*ClockDomain // Multiple clock domains
	Mutexes    []*Mutex        // Hardware mutexes
	Templates  []*ModuleTemplate // Polymorphic templates
}

func (m *Module) Input(name string, width Width) *Signal {
	s := &Signal{Name: name, Width: width, Kind: "input"}
	m.Inputs = append(m.Inputs, s)
	return s
}

func (m *Module) Output(name string, width Width) *Signal {
	s := &Signal{Name: name, Width: width, Kind: "output"}
	m.Outputs = append(m.Outputs, s)
	return s
}

func (m *Module) Wire(name string, width Width) *Signal {
	s := &Signal{Name: name, Width: width, Kind: "wire"}
	m.Wires = append(m.Wires, s)
	return s
}

func (m *Module) Reg(name string, width Width) *Signal {
	s := &Signal{Name: name, Width: width, Kind: "reg"}
	m.Regs = append(m.Regs, s)
	return s
}

func (m *Module) Bundle(name string) *Bundle {
	b := &Bundle{
		Name:   name,
		Fields: make(map[string]*Signal),
	}
	m.Bundles = append(m.Bundles, b)
	return b
}

func (m *Module) Vec(name string, size int, width Width) *Vec {
	v := &Vec{
		Name:     name,
		Elements: make([]*Signal, size),
		Width:    width,
		Size:     size,
	}
	
	// Create individual signals for each element
	for i := 0; i < size; i++ {
		elemName := fmt.Sprintf("%s[%d]", name, i)
		v.Elements[i] = &Signal{Name: elemName, Width: width, Kind: "wire"}
	}
	
	m.Vecs = append(m.Vecs, v)
	return v
}

// Bundle methods
func (b *Bundle) AddField(name string, signal *Signal) {
	fieldName := b.Name + "_" + name
	newSignal := &Signal{
		Name:  fieldName,
		Width: signal.Width,
		Kind:  signal.Kind,
		Expr:  signal.Expr,
	}
	b.Fields[name] = newSignal
}

func (b *Bundle) GetField(name string) *Signal {
	return b.Fields[name]
}

func (b *Bundle) Connect(other *Bundle) []string {
	var connections []string
	for name, signal := range b.Fields {
		if otherSignal, exists := other.Fields[name]; exists {
			connections = append(connections, fmt.Sprintf("assign %s = %s;", signal.Name, otherSignal.Name))
		}
	}
	return connections
}

// Vec methods
func (v *Vec) Get(index int) *Signal {
	if index >= 0 && index < v.Size {
		return v.Elements[index]
	}
	return nil
}

func (v *Vec) Set(index int, signal *Signal) {
	if index >= 0 && index < v.Size {
		v.Elements[index] = signal
	}
}

func (v *Vec) Connect(other *Vec) []string {
	var connections []string
	size := v.Size
	if other.Size < size {
		size = other.Size
	}
	
	for i := 0; i < size; i++ {
		connections = append(connections, fmt.Sprintf("assign %s = %s;", v.Elements[i].Name, other.Elements[i].Name))
	}
	return connections
}

// Memory methods
func (m *Module) SyncMem(name string, width Width, depth int) *Memory {
	addrWidth := Width(0)
	for d := depth - 1; d > 0; d >>= 1 {
		addrWidth++
	}
	if addrWidth == 0 {
		addrWidth = 1
	}
	
	mem := &Memory{
		Name:      name,
		Width:     width,
		Depth:     depth,
		AddrWidth: addrWidth,
		Kind:      "sync",
	}
	m.Memories = append(m.Memories, mem)
	return mem
}

func (m *Module) AsyncMem(name string, width Width, depth int) *Memory {
	addrWidth := Width(0)
	for d := depth - 1; d > 0; d >>= 1 {
		addrWidth++
	}
	if addrWidth == 0 {
		addrWidth = 1
	}
	
	mem := &Memory{
		Name:      name,
		Width:     width,
		Depth:     depth,
		AddrWidth: addrWidth,
		Kind:      "async",
	}
	m.Memories = append(m.Memories, mem)
	return mem
}

func (mem *Memory) Read(addr *Signal) *Signal {
	return &Signal{
		Name:  mem.Name + "[" + addr.Name + "]",
		Width: mem.Width,
		Kind:  "wire",
	}
}

func (mem *Memory) Write(addr *Signal, data *Signal, enable *Signal) string {
	if mem.Kind == "sync" {
		return fmt.Sprintf("if (%s) %s[%s] <= %s;", enable.Name, mem.Name, addr.Name, data.Name)
	} else {
		return fmt.Sprintf("assign %s[%s] = (%s) ? %s : %s[%s];", 
			mem.Name, addr.Name, enable.Name, data.Name, mem.Name, addr.Name)
	}
}

func (m *Module) SetClock(clk *Signal) {
	m.Clock = clk
}

func (m *Module) SetReset(rst *Signal) {
	m.Reset = rst
}

func (m *Module) SetParameter(name string, value interface{}) {
	if m.Parameters == nil {
		m.Parameters = make(map[string]interface{})
	}
	m.Parameters[name] = value
}

// Module instantiation
func (m *Module) Instance(moduleName, instanceName string) *ModuleInstance {
	inst := &ModuleInstance{
		ModuleName:   moduleName,
		InstanceName: instanceName,
		Parameters:   make(map[string]interface{}),
		Connections:  make(map[string]*Signal),
	}
	m.Instances = append(m.Instances, inst)
	return inst
}

func (inst *ModuleInstance) SetParameter(name string, value interface{}) *ModuleInstance {
	inst.Parameters[name] = value
	return inst
}

func (inst *ModuleInstance) Connect(portName string, signal *Signal) *ModuleInstance {
	inst.Connections[portName] = signal
	return inst
}

func (inst *ModuleInstance) IO(portName string, width Width) *Signal {
	// Create a signal for the connection
	signalName := inst.InstanceName + "_" + portName
	signal := &Signal{
		Name:  signalName,
		Width: width,
		Kind:  "wire",
	}
	inst.Connections[portName] = signal
	return signal
}

// Multiplexer function with width inference
func Mux(sel *Signal, inputs ...*Signal) *Signal {
	if len(inputs) == 0 {
		return &Signal{Name: "", Width: 0, Kind: "wire"}
	}
	if len(inputs) == 1 {
		return inputs[0]
	}
	
	// Determine result width from inputs
	resultWidth := Width(0)
	for _, input := range inputs {
		if input.Width > resultWidth {
			resultWidth = input.Width
		}
	}
	
	// Generate case statement for multiplexer
	result := "("
	for i := len(inputs) - 1; i >= 0; i-- {
		if i == len(inputs) - 1 {
			result += inputs[i].Name
		} else {
			result += fmt.Sprintf(" : (%s == %d) ? %s", sel.Name, i, inputs[i].Name)
		}
	}
	result += ")"
	
	return &Signal{
		Name:  result,
		Width: resultWidth,
		Kind:  "wire",
	}
}

func (m *Module) Assign(lhs *Signal, rhs *Signal) {
	m.Assigns = append(m.Assigns, fmt.Sprintf("assign %s = %s;", lhs.Name, rhs.Name))
}

func (m *Module) AssignExpr(lhs *Signal, expr string) {
	m.Assigns = append(m.Assigns, fmt.Sprintf("assign %s = %s;", lhs.Name, expr))
}

// Legacy support - assign using signal name
func (m *Module) AssignName(lhsName string, rhsExpr string) {
	m.Assigns = append(m.Assigns, fmt.Sprintf("assign %s = %s;", lhsName, rhsExpr))
}

// Literal signal creation helpers
func Lit(value int, width Width) *Signal {
	return &Signal{
		Name:  fmt.Sprintf("%d'h%x", width, value),
		Width: width,
		Kind:  "wire",
	}
}

func Bool(value bool) *Signal {
	if value {
		return &Signal{Name: "1'b1", Width: 1, Kind: "wire"}
	}
	return &Signal{Name: "1'b0", Width: 1, Kind: "wire"}
}



// Arithmetic operations with width inference
func (s *Signal) Add(other *Signal) *Signal {
	resultWidth := MaxWidth(s.Width, other.Width)
	return &Signal{
		Name:  s.Name + " + " + other.Name,
		Width: resultWidth,
		Kind:  "wire",
	}
}

func (s *Signal) Sub(other *Signal) *Signal {
	resultWidth := MaxWidth(s.Width, other.Width)
	return &Signal{
		Name:  s.Name + " - " + other.Name,
		Width: resultWidth,
		Kind:  "wire",
	}
}

func (s *Signal) Mul(other *Signal) *Signal {
	resultWidth := s.Width + other.Width // Multiplication doubles width
	return &Signal{
		Name:  s.Name + " * " + other.Name,
		Width: resultWidth,
		Kind:  "wire",
	}
}

func (s *Signal) Div(other *Signal) *Signal {
	resultWidth := s.Width // Division keeps numerator width
	return &Signal{
		Name:  s.Name + " / " + other.Name,
		Width: resultWidth,
		Kind:  "wire",
	}
}

func (s *Signal) Mod(other *Signal) *Signal {
	resultWidth := other.Width // Modulo has denominator width
	return &Signal{
		Name:  s.Name + " % " + other.Name,
		Width: resultWidth,
		Kind:  "wire",
	}
}

// Bitwise operations preserve width
func (s *Signal) And(other *Signal) *Signal {
	resultWidth := MaxWidth(s.Width, other.Width)
	return &Signal{
		Name:  s.Name + " & " + other.Name,
		Width: resultWidth,
		Kind:  "wire",
	}
}

func (s *Signal) Or(other *Signal) *Signal {
	resultWidth := MaxWidth(s.Width, other.Width)
	return &Signal{
		Name:  s.Name + " | " + other.Name,
		Width: resultWidth,
		Kind:  "wire",
	}
}

func (s *Signal) Xor(other *Signal) *Signal {
	resultWidth := MaxWidth(s.Width, other.Width)
	return &Signal{
		Name:  s.Name + " ^ " + other.Name,
		Width: resultWidth,
		Kind:  "wire",
	}
}

func (s *Signal) Not() *Signal {
	return &Signal{
		Name:  "~" + s.Name,
		Width: s.Width,
		Kind:  "wire",
	}
}

// Logical operations return 1-bit results
func (s *Signal) LogicAnd(other *Signal) *Signal {
	return &Signal{
		Name:  s.Name + " && " + other.Name,
		Width: 1,
		Kind:  "wire",
	}
}

func (s *Signal) LogicOr(other *Signal) *Signal {
	return &Signal{
		Name:  s.Name + " || " + other.Name,
		Width: 1,
		Kind:  "wire",
	}
}

func (s *Signal) LogicNot() *Signal {
	return &Signal{
		Name:  "!" + s.Name,
		Width: 1,
		Kind:  "wire",
	}
}

// Comparison operations return 1-bit results
func (s *Signal) Eq(other *Signal) *Signal {
	return &Signal{
		Name:  s.Name + " == " + other.Name,
		Width: 1,
		Kind:  "wire",
	}
}

func (s *Signal) Neq(other *Signal) *Signal {
	return &Signal{
		Name:  s.Name + " != " + other.Name,
		Width: 1,
		Kind:  "wire",
	}
}

func (s *Signal) Lt(other *Signal) *Signal {
	return &Signal{
		Name:  s.Name + " < " + other.Name,
		Width: 1,
		Kind:  "wire",
	}
}

func (s *Signal) Lte(other *Signal) *Signal {
	return &Signal{
		Name:  s.Name + " <= " + other.Name,
		Width: 1,
		Kind:  "wire",
	}
}

func (s *Signal) Gt(other *Signal) *Signal {
	return &Signal{
		Name:  s.Name + " > " + other.Name,
		Width: 1,
		Kind:  "wire",
	}
}

func (s *Signal) Gte(other *Signal) *Signal {
	return &Signal{
		Name:  s.Name + " >= " + other.Name,
		Width: 1,
		Kind:  "wire",
	}
}

// Shift operations
func (s *Signal) Shl(amount int) *Signal {
	return &Signal{
		Name:  s.Name + " << " + fmt.Sprintf("%d", amount),
		Width: s.Width + Width(amount), // Left shift increases width
		Kind:  "wire",
	}
}

func (s *Signal) Shr(amount int) *Signal {
	newWidth := s.Width - Width(amount)
	if newWidth < 1 {
		newWidth = 1
	}
	return &Signal{
		Name:  s.Name + " >> " + fmt.Sprintf("%d", amount),
		Width: newWidth, // Right shift decreases width
		Kind:  "wire",
	}
}

// Bit selection and concatenation
func (s *Signal) Bits(high, low int) *Signal {
	width := Width(high - low + 1)
	if high == low {
		return &Signal{
			Name:  s.Name + "[" + fmt.Sprintf("%d", high) + "]",
			Width: 1,
			Kind:  "wire",
		}
	}
	return &Signal{
		Name:  s.Name + "[" + fmt.Sprintf("%d", high) + ":" + fmt.Sprintf("%d", low) + "]",
		Width: width,
		Kind:  "wire",
	}
}

// Cat method for signals - instance method version
func (s *Signal) Cat(others ...*Signal) *Signal {
	allSignals := append([]*Signal{s}, others...)
	return Cat(allSignals...)
}

func Cat(signals ...*Signal) *Signal {
	if len(signals) == 0 {
		return &Signal{Name: "", Width: 0, Kind: "wire"}
	}
	
	totalWidth := Width(0)
	names := make([]string, len(signals))
	for i, sig := range signals {
		totalWidth += sig.Width
		names[i] = sig.Name
	}
	
	result := "{"
	for i, name := range names {
		if i > 0 {
			result += ", "
		}
		result += name
	}
	result += "}"
	
	return &Signal{
		Name:  result,
		Width: totalWidth,
		Kind:  "wire",
	}
}

func Fill(n int, value string) *Signal {
	return &Signal{
		Name:  "{" + fmt.Sprintf("%d", n) + "{" + value + "}}",
		Width: Width(n), // Assumes value is 1 bit
		Kind:  "wire",
	}
}

// === CLOCK DOMAIN METHODS ===

// Create a new clock domain
func (m *Module) NewClockDomain(name string, clk *Signal, rst *Signal) *ClockDomain {
	cd := &ClockDomain{
		Name:  name,
		Clock: clk,
		Reset: rst,
	}
	m.ClockDomains = append(m.ClockDomains, cd)
	return cd
}

// Set frequency for documentation purposes
func (cd *ClockDomain) SetFrequency(freq int) *ClockDomain {
	cd.Frequency = freq
	return cd
}

// Associate a signal with a clock domain
func (s *Signal) WithClockDomain(cd *ClockDomain) *Signal {
	s.ClockDomain = cd
	return s
}

// Clock Domain Crossing (CDC) synchronizer
func (m *Module) CDCSynchronizer(name string, srcSignal *Signal, destDomain *ClockDomain, stages int) *Signal {
	if stages < 2 {
		stages = 2 // Minimum 2 stages for proper synchronization
	}
	
	var syncStages []*Signal
	for i := 0; i < stages; i++ {
		stageName := fmt.Sprintf("%s_sync_stage%d", name, i)
		stage := m.Reg(stageName, srcSignal.Width).WithClockDomain(destDomain)
		syncStages = append(syncStages, stage)
	}
	
	// Connect the synchronizer chain
	for i := 0; i < stages; i++ {
		if i == 0 {
			m.Always = append(m.Always, fmt.Sprintf("always @(posedge %s) %s <= %s;", 
				destDomain.Clock.Name, syncStages[i].Name, srcSignal.Name))
		} else {
			m.Always = append(m.Always, fmt.Sprintf("always @(posedge %s) %s <= %s;", 
				destDomain.Clock.Name, syncStages[i].Name, syncStages[i-1].Name))
		}
	}
	
	return syncStages[stages-1]
}

// Asynchronous FIFO for CDC
func (m *Module) AsyncFIFO(name string, width Width, depth int, wrDomain, rdDomain *ClockDomain) (*Signal, *Signal, *Signal, *Signal) {
	// Calculate address width
	addrWidth := Width(0)
	for d := depth - 1; d > 0; d >>= 1 {
		addrWidth++
	}
	if addrWidth == 0 {
		addrWidth = 1
	}
	
	// Create FIFO memory
	mem := m.SyncMem(name+"_mem", width, depth)
	
	// Write domain signals
	wrData := m.Input(name+"_wr_data", width).WithClockDomain(wrDomain)
	wrEn := m.Input(name+"_wr_en", 1).WithClockDomain(wrDomain)
	wrFull := m.Output(name+"_wr_full", 1).WithClockDomain(wrDomain)
	
	// Read domain signals
	rdData := m.Output(name+"_rd_data", width).WithClockDomain(rdDomain)
	rdEn := m.Input(name+"_rd_en", 1).WithClockDomain(rdDomain)
	rdEmpty := m.Output(name+"_rd_empty", 1).WithClockDomain(rdDomain)
	
	// Gray code pointers for CDC
	wrPtr := m.Reg(name+"_wr_ptr", addrWidth+1).WithClockDomain(wrDomain)
	rdPtr := m.Reg(name+"_rd_ptr", addrWidth+1).WithClockDomain(rdDomain)
	
	// Cross-domain synchronizers
	wrPtrSync := m.CDCSynchronizer(name+"_wr_ptr_sync", wrPtr, rdDomain, 2)
	rdPtrSync := m.CDCSynchronizer(name+"_rd_ptr_sync", rdPtr, wrDomain, 2)
	
	// FIFO logic would be added here (simplified for brevity)
	m.Always = append(m.Always, fmt.Sprintf("// AsyncFIFO %s implementation", name))
	m.Always = append(m.Always, fmt.Sprintf("// Memory: %s", mem.Name))
	m.Always = append(m.Always, fmt.Sprintf("// Write enable: %s", wrEn.Name))
	m.Always = append(m.Always, fmt.Sprintf("// Read enable: %s", rdEn.Name))
	m.Always = append(m.Always, fmt.Sprintf("// Write pointer sync: %s", wrPtrSync.Name))
	m.Always = append(m.Always, fmt.Sprintf("// Read pointer sync: %s", rdPtrSync.Name))
	
	return wrData, rdData, wrFull, rdEmpty
}

// === MUTEX METHODS ===

// Create a hardware mutex with multiple requesters
func (m *Module) Mutex(name string, numRequesters int, arbitration string) *Mutex {
	if arbitration == "" {
		arbitration = "round_robin"
	}
	
	mutex := &Mutex{
		Name:        name,
		Width:       Width(numRequesters),
		Arbitration: arbitration,
	}
	
	// Create request and grant signals
	for i := 0; i < numRequesters; i++ {
		req := m.Input(fmt.Sprintf("%s_req_%d", name, i), 1)
		grant := m.Output(fmt.Sprintf("%s_grant_%d", name, i), 1)
		mutex.Requests = append(mutex.Requests, req)
		mutex.Grants = append(mutex.Grants, grant)
	}
	
	m.Mutexes = append(m.Mutexes, mutex)
	return mutex
}

// Round-robin arbitration
func (mutex *Mutex) GenerateRoundRobin(m *Module, clk *Signal) {
	if len(mutex.Requests) == 0 {
		return
	}
	
	// Create round-robin counter
	counterWidth := WidthOf(len(mutex.Requests) - 1)
	counter := m.Reg(mutex.Name+"_counter", counterWidth)
	
	// Generate arbitration logic
	m.Always = append(m.Always, fmt.Sprintf("always @(posedge %s) begin", clk.Name))
	
	for i, grant := range mutex.Grants {
		condition := fmt.Sprintf("(%s == %d) && %s", counter.Name, i, mutex.Requests[i].Name)
		m.Always = append(m.Always, fmt.Sprintf("  %s <= %s;", grant.Name, condition))
	}
	
	// Update counter
	anyGrant := strings.Join(func() []string {
		var grants []string
		for _, g := range mutex.Grants {
			grants = append(grants, g.Name)
		}
		return grants
	}(), " || ")
	
	m.Always = append(m.Always, fmt.Sprintf("  if (%s) %s <= (%s + 1) %% %d;", 
		anyGrant, counter.Name, counter.Name, len(mutex.Requests)))
	m.Always = append(m.Always, "end")
}

// Priority arbitration (highest index wins)
func (mutex *Mutex) GeneratePriority(m *Module) {
	if len(mutex.Requests) == 0 {
		return
	}
	
	// Generate priority encoder
	for i := len(mutex.Grants) - 1; i >= 0; i-- {
		if i == len(mutex.Grants) - 1 {
			// Highest priority
			m.Assign(mutex.Grants[i], mutex.Requests[i])
		} else {
			// Lower priority - only grant if higher priorities not requesting
			higherReqs := []string{}
			for j := i + 1; j < len(mutex.Requests); j++ {
				higherReqs = append(higherReqs, mutex.Requests[j].Name)
			}
			condition := fmt.Sprintf("%s && !(%s)", mutex.Requests[i].Name, strings.Join(higherReqs, " || "))
			m.AssignExpr(mutex.Grants[i], condition)
		}
	}
}

// === POLYMORPHISM METHODS ===

// Create a polymorphic module template
func NewModuleTemplate(name string, typeParams []string) *ModuleTemplate {
	return &ModuleTemplate{
		Name:       name,
		TypeParams: typeParams,
		Params:     make(map[string]interface{}),
	}
}

// Set generator function for the template
func (mt *ModuleTemplate) SetGenerator(gen func(params map[string]interface{}) *Module) *ModuleTemplate {
	mt.Generator = gen
	return mt
}

// Add type constraint
func (mt *ModuleTemplate) AddConstraint(param string, constraint interface{}) *ModuleTemplate {
	mt.Params[param] = constraint
	return mt
}

// Instantiate a polymorphic module with specific type parameters
func (m *Module) InstantiateTemplate(template *ModuleTemplate, instanceName string, typeArgs map[string]interface{}) *Module {
	// Validate type arguments against constraints
	for param, constraint := range template.Params {
		if arg, exists := typeArgs[param]; exists {
			// Simple type checking (can be extended for more complex constraints)
			switch constraint.(type) {
			case Width:
				if _, ok := arg.(Width); !ok {
					panic(fmt.Sprintf("Type parameter %s must be Width, got %T", param, arg))
				}
			case int:
				if _, ok := arg.(int); !ok {
					panic(fmt.Sprintf("Type parameter %s must be int, got %T", param, arg))
				}
			}
		} else {
			panic(fmt.Sprintf("Missing type argument for parameter %s", param))
		}
	}
	
	// Generate the module instance
	if template.Generator == nil {
		panic(fmt.Sprintf("No generator function set for template %s", template.Name))
	}
	
	instance := template.Generator(typeArgs)
	instance.Name = instanceName
	
	return instance
}

// Generic FIFO template
func FIFOTemplate() *ModuleTemplate {
	template := NewModuleTemplate("GenericFIFO", []string{"DATA_WIDTH", "DEPTH"})
	template.AddConstraint("DATA_WIDTH", Width(0))
	template.AddConstraint("DEPTH", int(0))
	
	template.SetGenerator(func(params map[string]interface{}) *Module {
		dataWidth := params["DATA_WIDTH"].(Width)
		depth := params["DEPTH"].(int)
		
		m := &Module{
			Name:       "GenericFIFO",
			Inputs:     []*Signal{},
			Outputs:    []*Signal{},
			Wires:      []*Signal{},
			Regs:       []*Signal{},
			Assigns:    []string{},
			Always:     []string{},
			Parameters: make(map[string]interface{}),
		}
		
		// Create FIFO interface
		clk := m.Input("clk", 1)
		rst := m.Input("rst", 1)
		wrData := m.Input("wr_data", dataWidth)
		wrEn := m.Input("wr_en", 1)
		wrFull := m.Output("wr_full", 1)
		rdData := m.Output("rd_data", dataWidth)
		rdEn := m.Input("rd_en", 1)
		rdEmpty := m.Output("rd_empty", 1)
		
		// Add FIFO implementation
		mem := m.SyncMem("fifo_mem", dataWidth, depth)
		
		// Use the signals in comments to avoid unused variable errors
		m.Always = append(m.Always, fmt.Sprintf("// FIFO clock: %s, reset: %s", clk.Name, rst.Name))
		m.Always = append(m.Always, fmt.Sprintf("// Write data: %s, enable: %s, full: %s", wrData.Name, wrEn.Name, wrFull.Name))
		m.Always = append(m.Always, fmt.Sprintf("// Read data: %s, enable: %s, empty: %s", rdData.Name, rdEn.Name, rdEmpty.Name))
		m.Always = append(m.Always, fmt.Sprintf("// Memory: %s", mem.Name))
		
		// Calculate address width
		addrWidth := WidthOf(depth - 1)
		wrPtr := m.Reg("wr_ptr", addrWidth)
		rdPtr := m.Reg("rd_ptr", addrWidth)
		
		// FIFO logic (simplified)
		m.Always = append(m.Always, fmt.Sprintf("always @(posedge %s) begin", clk.Name))
		m.Always = append(m.Always, fmt.Sprintf("  if (%s) begin", rst.Name))
		m.Always = append(m.Always, fmt.Sprintf("    %s <= 0;", wrPtr.Name))
		m.Always = append(m.Always, fmt.Sprintf("    %s <= 0;", rdPtr.Name))
		m.Always = append(m.Always, "  end else begin")
		m.Always = append(m.Always, fmt.Sprintf("    if (%s && !%s) %s <= %s + 1;", wrEn.Name, wrFull.Name, wrPtr.Name, wrPtr.Name))
		m.Always = append(m.Always, fmt.Sprintf("    if (%s && !%s) %s <= %s + 1;", rdEn.Name, rdEmpty.Name, rdPtr.Name, rdPtr.Name))
		m.Always = append(m.Always, "  end")
		m.Always = append(m.Always, "end")
		
		// Set parameters for Verilog generation
		m.SetParameter("DATA_WIDTH", dataWidth)
		m.SetParameter("DEPTH", depth)
		
		return m
	})
	
	return template
}

// Generic register file template
func RegisterFileTemplate() *ModuleTemplate {
	template := NewModuleTemplate("GenericRegisterFile", []string{"DATA_WIDTH", "ADDR_WIDTH"})
	template.AddConstraint("DATA_WIDTH", Width(0))
	template.AddConstraint("ADDR_WIDTH", Width(0))
	
	template.SetGenerator(func(params map[string]interface{}) *Module {
		dataWidth := params["DATA_WIDTH"].(Width)
		addrWidth := params["ADDR_WIDTH"].(Width)
		depth := 1 << addrWidth
		
		m := &Module{
			Name:       "GenericRegisterFile",
			Inputs:     []*Signal{},
			Outputs:    []*Signal{},
			Wires:      []*Signal{},
			Regs:       []*Signal{},
			Assigns:    []string{},
			Always:     []string{},
			Parameters: make(map[string]interface{}),
		}
		
		// Create register file interface
		clk := m.Input("clk", 1)
		rst := m.Input("rst", 1)
		
		// Write port
		wrAddr := m.Input("wr_addr", addrWidth)
		wrData := m.Input("wr_data", dataWidth)
		wrEn := m.Input("wr_en", 1)
		
		// Read ports (dual-port)
		rd1Addr := m.Input("rd1_addr", addrWidth)
		rd1Data := m.Output("rd1_data", dataWidth)
		rd2Addr := m.Input("rd2_addr", addrWidth)
		rd2Data := m.Output("rd2_data", dataWidth)
		
		// Create register array
		regFile := m.SyncMem("reg_file", dataWidth, depth)
		
		// Implement register file logic
		m.Always = append(m.Always, fmt.Sprintf("always @(posedge %s) begin", clk.Name))
		m.Always = append(m.Always, fmt.Sprintf("  if (%s) begin", rst.Name))
		m.Always = append(m.Always, "    // Reset logic if needed")
		m.Always = append(m.Always, "  end else begin")
		m.Always = append(m.Always, fmt.Sprintf("    if (%s) %s[%s] <= %s;", wrEn.Name, regFile.Name, wrAddr.Name, wrData.Name))
		m.Always = append(m.Always, "  end")
		m.Always = append(m.Always, "end")
		
		// Continuous read
		m.Assign(rd1Data, regFile.Read(rd1Addr))
		m.Assign(rd2Data, regFile.Read(rd2Addr))
		
		// Set parameters
		m.SetParameter("DATA_WIDTH", dataWidth)
		m.SetParameter("ADDR_WIDTH", addrWidth)
		
		return m
	})
	
	return template
}

