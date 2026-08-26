package model

import "math"

func (m *Model) NodeByID(id string) (Node, int, bool) {
	for i, n := range m.Nodes {
		if n.ID == id {
			return n, i, true
		}
	}
	return Node{}, 0, false
}

func (m *Model) ElementsIncident(id string) []Element {
	var out []Element
	for _, e := range m.Elements {
		if e.From == id || e.To == id {
			out = append(out, e)
		}
	}
	return out
}

func (m *Model) MaxSpan() float64 {
	idx := NodeIndex(m)
	max := 0.0
	for _, e := range m.Elements {
		fi, ok1 := idx[e.From]
		ti, ok2 := idx[e.To]
		if !ok1 || !ok2 {
			continue
		}
		nf, nt := m.Nodes[fi], m.Nodes[ti]
		if L := math.Hypot(nt.X-nf.X, nt.Y-nf.Y); L > max {
			max = L
		}
	}
	return max
}

func (m *Model) RestrainedCount() int {
	n := 0
	for _, nd := range m.Nodes {
		r := nd.EffectiveRestraint()
		for _, b := range r {
			if b {
				n++
			}
		}
	}
	return n
}

func (m *Model) FreeCount() int {
	return TotalDOF(m) - m.RestrainedCount()
}

func (m *Model) Clone() *Model {
	out := &Model{
		Nodes:    make([]Node, len(m.Nodes)),
		Elements: make([]Element, len(m.Elements)),
		Loads:    make([]NodeLoad, len(m.Loads)),
	}
	copy(out.Nodes, m.Nodes)
	copy(out.Elements, m.Elements)
	copy(out.Loads, m.Loads)
	return out
}
