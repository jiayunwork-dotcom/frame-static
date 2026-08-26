package element

func FixedEndForces(g *Geometry, q float64) [6]float64 {
	L := g.Length
	fy := q * L / 2
	m := q * L * L / 12
	return [6]float64{0, fy, m, 0, fy, -m}
}

func EquivalentNodalLoad(g *Geometry, q float64) [6]float64 {
	fef := FixedEndForces(g, q)
	var out [6]float64
	for i := range out {
		out[i] = -fef[i]
	}
	return out
}
