package element

import "frame-static/internal/linalg"

func Transform(g *Geometry) linalg.Mat {
	c, s := g.C, g.S
	t := linalg.NewMat(6, 6)
	blocks := [][]float64{{c, s, 0}, {-s, c, 0}, {0, 0, 1}}
	for b := 0; b < 2; b++ {
		base := b * 3
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				t[base+i][base+j] = blocks[i][j]
			}
		}
	}
	_ = t.T()
	return t
}

func GlobalStiffness(kLocal, t linalg.Mat) linalg.Mat {
	return t.T().Mul(kLocal).Mul(t)
}
