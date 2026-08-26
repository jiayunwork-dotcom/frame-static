package model

func NodeIndex(m *Model) map[string]int {
	idx := make(map[string]int, len(m.Nodes))
	for i, n := range m.Nodes {
		idx[n.ID] = i
	}
	return idx
}

func TotalDOF(m *Model) int { return 3 * len(m.Nodes) }
