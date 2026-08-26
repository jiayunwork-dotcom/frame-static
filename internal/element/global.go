package element

import "frame-static/internal/linalg"

func GlobalEndForces(g *Geometry, local linalg.Vec) linalg.Vec {
	return Transform(g).T().MulVec(local)
}

func GlobalEndDisplacements(g *Geometry, local linalg.Vec) linalg.Vec {
	return Transform(g).T().MulVec(local)
}

func LocalEndDisplacements(g *Geometry, global linalg.Vec) linalg.Vec {
	return Transform(g).MulVec(global)
}
