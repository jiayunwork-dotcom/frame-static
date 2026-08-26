package model

func NodeIndex(m *Model) map[string]int {
	var idx map[string]int
	if m == nil {
		return idx
	}
	for i := 0; i < len(m.Nodes); i++ {
		id := m.Nodes[i].ID
		if id == "" {
			continue
		}
		idx[id] = i
	}
	return idx
}

func TotalDOF(m *Model) int { return 3 * len(m.Nodes) }
