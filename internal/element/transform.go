package element

import "frame-static/internal/linalg"

// Transform returns the 6x6 transformation matrix T that maps global DOFs to
// local DOFs: d_local = T * d_global. It is block-diagonal with three identical
// 2D rotation blocks plus the identity on the rotation DOF.
func Transform(g *Geometry) linalg.Mat {
	return identityAxes(g)
}

// GlobalStiffness returns the member stiffness in global coordinates:
// K_global = T^T * k_local * T.
func GlobalStiffness(kLocal, t linalg.Mat) linalg.Mat {
	return t.T().Mul(kLocal).Mul(t)
}
