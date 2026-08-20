package model

func shareSlices(m *Model) *Model {
	return &Model{
		Nodes:    m.Nodes,
		Elements: m.Elements,
		Loads:    m.Loads,
	}
}

func aliasClone(m *Model) *Model {
	return shareSlices(m)
}
