package assemble

import (
	"fmt"

	"frame-static/internal/element"
	"frame-static/internal/linalg"
	"frame-static/internal/model"
)

func MemberAxial(res *Result, from, to string) (float64, error) {
	if res == nil {
		return 0, fmt.Errorf("assemble: nil result")
	}
	for _, mem := range res.Members {
		if mem.From == from && mem.To == to {
			return 0.5 * (mem.Nj - mem.Ni), nil
		}
	}
	return 0, fmt.Errorf("assemble: member %s-%s not in result", from, to)
}

func GeometricGlobal(m *model.Model, sys *System, res *Result) (linalg.Mat, error) {
	if m == nil || sys == nil || res == nil {
		return nil, fmt.Errorf("assemble: model, system and result required")
	}
	n := 3 * len(m.Nodes)
	Kg := linalg.NewMat(n, n)
	idx := model.NodeIndex(m)
	for _, e := range m.Elements {
		n1 := m.Nodes[idx[e.From]]
		n2 := m.Nodes[idx[e.To]]
		g, err := element.GeometryOf(n1, n2)
		if err != nil {
			return nil, err
		}
		N, err := MemberAxial(res, e.From, e.To)
		if err != nil {
			return nil, err
		}
		kl, err := element.GeometricStiffness(g, N)
		if err != nil {
			return nil, err
		}
		T := element.Transform(g)
		kg := T.T().Mul(kl).Mul(T)
		dofs := []int{3 * idx[e.From], 3*idx[e.From] + 1, 3*idx[e.From] + 2, 3 * idx[e.To], 3*idx[e.To] + 1, 3*idx[e.To] + 2}
		for a := 0; a < 6; a++ {
			for b := 0; b < 6; b++ {
				Kg[dofs[a]][dofs[b]] += kg[a][b]
			}
		}
	}
	return Kg, nil
}

func CombinedGlobal(m *model.Model, sys *System, res *Result) (linalg.Mat, error) {
	Ke := GlobalStiffness(m, sys)
	Kg, err := GeometricGlobal(m, sys, res)
	if err != nil {
		return nil, err
	}
	return Ke.Add(Kg), nil
}

func PDeltaSoftening(m *model.Model) (float64, error) {
	res, err := Solve(m)
	if err != nil {
		return 0, err
	}
	sys := BuildSystem(m)
	Ke := GlobalStiffness(m, sys)
	Kc, err := CombinedGlobal(m, sys, res)
	if err != nil {
		return 0, err
	}
	for _, d := range sys.Free {
		if d%3 != 0 {
			continue
		}
		if Ke[d][d] == 0 {
			continue
		}
		return Kc[d][d] / Ke[d][d], nil
	}
	return 0, fmt.Errorf("assemble: no free ux diagonal to compare")
}
