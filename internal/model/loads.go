package model

func (m *Model) TotalAppliedForce() (fx, fy, mz float64) {
	for _, ld := range m.Loads {
		fx += ld.FX
		fy += ld.FY
		mz += ld.MZ
	}
	return
}

func (m *Model) HasDistLoads() bool {
	for _, e := range m.Elements {
		if e.Dist != nil {
			return true
		}
	}
	return false
}

func (m *Model) TotalDistLoad() float64 {
	var s float64
	for _, e := range m.Elements {
		if e.Dist != nil {
			s += e.Dist.Q
		}
	}
	return s
}

func (m *Model) LoadedNodes() []string {
	var out []string
	for _, ld := range m.Loads {
		out = append(out, ld.Node)
	}
	return out
}

func (m *Model) LoadCount() int { return len(m.Loads) }
