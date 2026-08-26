package element

import (
	"errors"
	"math"

	"frame-static/internal/model"
)

type Geometry struct {
	Length float64
	C      float64
	S      float64
}

func GeometryOf(n1, n2 model.Node) (*Geometry, error) {
	dx, dy := n2.X-n1.X, n2.Y-n1.Y
	L := math.Hypot(dx, dy)
	if L == 0 {
		return nil, errors.New("zero-length element")
	}
	return &Geometry{Length: L, C: dx / L, S: dy / L}, nil
}
