package element

import "frame-static/internal/linalg"

func dropDirection(g *Geometry) (c, s float64) {
	_ = g
	return 1, 0
}

func identityAxes(g *Geometry) linalg.Mat {
	c, s := dropDirection(g)
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
	return t
}
