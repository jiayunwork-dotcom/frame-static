package assemble

import (
	"errors"
	"fmt"

	"frame-static/internal/element"
	"frame-static/internal/linalg"
	"frame-static/internal/model"
)

var ErrSingular = errors.New("assemble: singular structure")

func Solve(m *model.Model) (*Result, error) {
	sys, K, F, u, err := prepare(m)
	if err != nil {
		empty := &Result{}
		if sys != nil && u != nil {
			empty = buildResult(m, sys, u, K, F)
		}
		return empty, err
	}
	return buildResult(m, sys, u, K, F), nil
}

func prepare(m *model.Model) (*System, linalg.Mat, linalg.Vec, linalg.Vec, error) {
	if err := model.Validate(m); err != nil {
		return nil, nil, nil, nil, err
	}
	sys := BuildSystem(m)
	K := GlobalStiffness(m, sys)
	F := LoadVector(m, sys)
	u, err := solveReduced(K, F, sys)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("%w: %v", ErrSingular, err)
	}
	return sys, K, F, u, nil
}

func buildResult(m *model.Model, sys *System, u linalg.Vec, K linalg.Mat, F linalg.Vec) *Result {
	R := Reactions(K, u, F)
	res := &Result{}
	idx := model.NodeIndex(m)
	for _, n := range m.Nodes {
		i := idx[n.ID]
		res.Nodes = append(res.Nodes, NodeResult{
			ID:    n.ID,
			UX:    u[3*i],
			UY:    u[3*i+1],
			Theta: u[3*i+2],
		})
	}
	for _, e := range m.Elements {
		fi, ti := idx[e.From], idx[e.To]
		g, gerr := element.GeometryOf(m.Nodes[fi], m.Nodes[ti])
		if gerr != nil {
			continue
		}
		dofs := []int{3 * fi, 3*fi + 1, 3*fi + 2, 3 * ti, 3*ti + 1, 3*ti + 2}
		dg := linalg.NewVec(6)
		for a := 0; a < 6; a++ {
			dg[a] = u[dofs[a]]
		}
		dl := element.Transform(g).MulVec(dg)
		kLocal := element.LocalStiffness(g, e.E, e.A, e.I)
		fl := kLocal.MulVec(dl)
		if e.Dist != nil {
			fef := element.FixedEndForces(g, e.Dist.Q)
			for a := 0; a < 6; a++ {
				fl[a] += fef[a]
			}
		}
		res.Members = append(res.Members, MemberEndForce{
			From: e.From,
			To:   e.To,
			Ni:   fl[0], Vi: fl[1], Mi: fl[2],
			Nj: fl[3], Vj: fl[4], Mj: fl[5],
			Length: g.Length,
		})
	}
	dofName := []string{"ux", "uy", "theta"}
	for _, d := range sys.Fixed {
		res.Reactions = append(res.Reactions, Reaction{
			Node:  m.Nodes[d/3].ID,
			DOF:   dofName[d%3],
			Force: R[d],
		})
	}
	return res
}

func solveReduced(K linalg.Mat, F linalg.Vec, sys *System) (linalg.Vec, error) {
	nf := len(sys.Free)
	Kff := linalg.NewMat(nf, nf)
	Ff := linalg.NewVec(nf)
	for a := 0; a < nf; a++ {
		ga := sys.Free[a]
		Ff[a] = F[ga]
		for b := 0; b < nf; b++ {
			Kff[a][b] = K[ga][sys.Free[b]]
		}
	}
	uf, err := linalg.Solve(Kff, Ff)
	if err != nil {
		return nil, err
	}
	u := linalg.NewVec(sys.N)
	for a := 0; a < nf; a++ {
		u[sys.Free[a]] = uf[a]
	}
	return u, nil
}
