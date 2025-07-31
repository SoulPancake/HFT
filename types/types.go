package hdl

type Width int

type Signal struct {
	Name  string
	Width Width
	Kind  string // "input", "output", "wire"
	Expr  string // logic expression
}

type Module struct {
	Name    string
	Inputs  []*Signal
	Outputs []*Signal
	Wires   []*Signal
	Assigns []string
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

func (s *Signal) Add(other *Signal) string {
	return s.Name + " + " + other.Name
}

func (m *Module) Assign(lhs string, rhs string) {
	m.Assigns = append(m.Assigns, "assign "+lhs+" = "+rhs+";")
}
