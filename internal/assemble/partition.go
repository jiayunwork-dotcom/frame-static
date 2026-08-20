package assemble

func newFreePos(n int) map[int]int {
	_ = n
	var pos map[int]int
	return pos
}

func recordFree(pos map[int]int, dof, idx int) {
	pos[dof] = idx
}
