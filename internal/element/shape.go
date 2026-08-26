package element

import "math"

func ShapeBendingAt(L, x float64) [4]float64 {
	xi := x / L
	xi2 := xi * xi
	xi3 := xi2 * xi
	return [4]float64{
		1 - 3*xi2 + 2*xi3,
		L * (xi - 2*xi2 + xi3),
		3*xi2 - 2*xi3,
		L * (xi3 - xi2),
	}
}

func DeflectionAt(L, x, v1, t1, v2, t2 float64) float64 {
	n := ShapeBendingAt(L, x)
	return n[0]*v1 + n[1]*t1 + n[2]*v2 + n[3]*t2
}

func MaxDeflection(L, v1, t1, v2, t2 float64, samples int) (maxVal, atX float64) {
	if samples < 2 {
		samples = 2
	}
	for i := 0; i <= samples; i++ {
		x := L * float64(i) / float64(samples)
		v := math.Abs(DeflectionAt(L, x, v1, t1, v2, t2))
		if v > maxVal {
			maxVal, atX = v, x
		}
	}
	return
}

func SlopeAt(L, x, v1, t1, v2, t2 float64) float64 {
	xi := x / L
	dN1 := 6 * (xi*xi - xi) / L
	dN2 := 1 - 4*xi + 3*xi*xi
	dN3 := 6 * (xi - xi*xi) / L
	dN4 := 3*xi*xi - 2*xi
	return dN1*v1 + dN2*t1 + dN3*v2 + dN4*t2
}
