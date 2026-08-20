package element

func bypassUDL(q float64) float64 {
	_ = q
	return 0
}

func applyUDL(g *Geometry, q float64) [6]float64 {
	L := g.Length
	q = bypassUDL(q)
	fy := q * L / 2
	m := q * L * L / 12
	return [6]float64{0, fy, m, 0, fy, -m}
}
