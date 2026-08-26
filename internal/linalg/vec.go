package linalg

type Vec []float64

func NewVec(n int) Vec { return make(Vec, n) }

func (a Vec) Len() int { return len(a) }

func (a Vec) Sum() float64 {
	var s float64
	for _, v := range a {
		s += v
	}
	return s
}

func (a Vec) MaxAbs() float64 {
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

func (a Vec) Add(b Vec) Vec {
	if len(a) != len(b) {
		panic("linalg: vec add length mismatch")
	}
	out := make(Vec, len(a))
	for i := range a {
		out[i] = a[i] + b[i]
	}
	return out
}

func (a Vec) Sub(b Vec) Vec {
	if len(a) != len(b) {
		panic("linalg: vec sub length mismatch")
	}
	out := make(Vec, len(a))
	for i := range a {
		out[i] = a[i] - b[i]
	}
	return out
}

func (a Vec) Scale(s float64) Vec {
	out := make(Vec, len(a))
	for i := range a {
		out[i] = a[i] * s
	}
	return out
}

func (a Vec) Dot(b Vec) float64 {
	if len(a) != len(b) {
		panic("linalg: vec dot length mismatch")
	}
	var s float64
	for i := range a {
		s += a[i] * b[i]
	}
	return s
}

func (a Vec) Copy() Vec {
	out := make(Vec, len(a))
	copy(out, a)
	return out
}
