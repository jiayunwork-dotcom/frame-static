package element

import (
	"math"

	"frame-static/internal/model"
)

func (g *Geometry) Midpoint(n1, n2 model.Node) (mx, my float64) {
	return (n1.X + n2.X) / 2, (n1.Y + n2.Y) / 2
}

func (g *Geometry) AngleDeg() float64 {
	return math.Atan2(g.S, g.C) * 180 / math.Pi
}

func (g *Geometry) PerpDirection() (cx, cy float64) {
	return -g.S, g.C
}

func (g *Geometry) LocalProjection(dx, dy float64) (axial, transverse float64) {
	axial = g.C*dx + g.S*dy
	transverse = -g.S*dx + g.C*dy
	return
}

func (g *Geometry) GlobalPoint(x0, y0, x float64) (gx, gy float64) {
	return x0 + g.C*x, y0 + g.S*x
}
