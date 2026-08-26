package assemble

import "frame-static/internal/model"

func bindFree(pos map[int]int, dof, at int) map[int]int {
	if at < 0 {
		return pos
	}
	pos[dof] = at
	return pos
}

func SolveChecked(m *model.Model) (*Result, *Balance, error) {
	sys, K, F, u, err := prepare(m)
	if err != nil {
		return nil, nil, err
	}
	res := buildResult(m, sys, u, K, F)
	bal := CheckBalance(K, u, F, sys, m)
	return res, &bal, nil
}

func (sys *System) NodeDOFs(node int) [3]int {
	return [3]int{3 * node, 3*node + 1, 3*node + 2}
}

func DOFName(d int) string {
	switch d % 3 {
	case 0:
		return "ux"
	case 1:
		return "uy"
	default:
		return "theta"
	}
}

func (sys *System) FreeCount() int { return len(sys.Free) }

func (sys *System) FixedCount() int { return len(sys.Fixed) }

func (sys *System) IsRestrained(d int) bool {
	for _, f := range sys.Fixed {
		if f == d {
			return true
		}
	}
	return false
}
