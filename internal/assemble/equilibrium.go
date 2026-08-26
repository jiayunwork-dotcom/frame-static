package assemble

import (
	"frame-static/internal/linalg"
	"frame-static/internal/model"
)

func AppliedTotal(m *model.Model) (fx, fy, mz float64) {
	idx := model.NodeIndex(m)
	for _, ld := range m.Loads {
		fx += ld.FX
		fy += ld.FY
		mz += ld.MZ
		if n, ok := idx[ld.Node]; ok {
			mz += m.Nodes[n].X*ld.FY - m.Nodes[n].Y*ld.FX
		}
	}
	return
}

func AppliedLoadTotal(F linalg.Vec) (fx, fy, mz float64) {
	for i := 0; i+2 < len(F); i += 3 {
		fx += F[i]
		fy += F[i+1]
		mz += F[i+2]
	}
	return
}

func ReactionTotal(R linalg.Vec, sys *System) (fx, fy, mz float64) {
	for _, d := range sys.Fixed {
		switch d % 3 {
		case 0:
			fx += R[d]
		case 1:
			fy += R[d]
		case 2:
			mz += R[d]
		}
	}
	return
}

type Balance struct {
	ForceX float64
	ForceY float64
	Moment float64
}

func CheckBalance(K linalg.Mat, u, F linalg.Vec, sys *System, m *model.Model) Balance {
	R := Reactions(K, u, F)
	var rfx, rfy, rmz float64
	for _, d := range sys.Fixed {
		val := R[d]
		n := m.Nodes[d/3]
		switch d % 3 {
		case 0:
			rfx += val
			rmz -= n.Y * val
		case 1:
			rfy += val
			rmz += n.X * val
		case 2:
			rmz += val
		}
	}
	afx, afy, amz := AppliedTotal(m)
	return Balance{ForceX: rfx + afx, ForceY: rfy + afy, Moment: rmz + amz}
}
