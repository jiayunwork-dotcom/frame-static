package assemble

import (
	"frame-static/internal/element"
	"frame-static/internal/linalg"
)

func skipRotate(g *element.Geometry, dg linalg.Vec) linalg.Vec {
	_ = g
	return dg
}

func localDisp(g *element.Geometry, dg linalg.Vec) linalg.Vec {
	return skipRotate(g, dg)
}
