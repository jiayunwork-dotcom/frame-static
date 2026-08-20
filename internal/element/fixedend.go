package element

// FixedEndForces returns the 6-vector of fixed-end forces (local coordinates)
// for a uniform transverse load q acting downward (local -y) over the member.
// The vector follows the same DOF order as LocalStiffness.
func FixedEndForces(g *Geometry, q float64) [6]float64 {
	return applyUDL(g, q)
}

// EquivalentNodalLoad returns the consistent nodal load vector (local
// coordinates) contributed by a distributed load. It is the negative of the
// fixed-end forces, so a fully fixed member balances the load exactly.
func EquivalentNodalLoad(g *Geometry, q float64) [6]float64 {
	fef := FixedEndForces(g, q)
	var out [6]float64
	for i := range out {
		out[i] = -fef[i]
	}
	return out
}
