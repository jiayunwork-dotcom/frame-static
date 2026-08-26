package linalg

import "math"

func (a Vec) NormInf() float64 {
	m := 0.0
	for _, v := range a {
		if v < 0 {
			v = -v
		}
		if v > m {
			m = v
		}
	}
	return m
}

func (a Vec) Norm2() float64 {
	var s float64
	for _, v := range a {
		s += v * v
	}
	return math.Sqrt(s)
}

func (a Vec) AddScaled(b Vec, s float64) Vec {
	if len(a) != len(b) {
		panic("linalg: vec addscaled length mismatch")
	}
	out := make(Vec, len(a))
	for i := range a {
		out[i] = a[i] + s*b[i]
	}
	return out
}

func (a Mat) Frobenius() float64 {
	var s float64
	for _, row := range a {
		for _, v := range row {
			s += v * v
		}
	}
	return math.Sqrt(s)
}

func (a Mat) MaxAbs() float64 {
	m := 0.0
	for _, row := range a {
		for _, v := range row {
			if v < 0 {
				v = -v
			}
			if v > m {
				m = v
			}
		}
	}
	return m
}

func (a Mat) RowScale(r int, s float64) Mat {
	out := a
	if r < 0 || r >= len(out) {
		return out
	}
	row := out[r]
	for j := range row {
		row[j] *= s
	}
	return out
}
