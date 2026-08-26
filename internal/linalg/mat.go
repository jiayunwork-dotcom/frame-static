package linalg

type Mat [][]float64

func NewMat(r, c int) Mat {
	m := make(Mat, r)
	for i := range m {
		m[i] = make([]float64, c)
	}
	return m
}

func (a Mat) Rows() int { return len(a) }

func (a Mat) Cols() int {
	if len(a) == 0 {
		return 0
	}
	return len(a[0])
}

func (a Mat) Mul(b Mat) Mat {
	if len(a) == 0 || len(b) == 0 {
		panic("linalg: empty mat mul")
	}
	ar, ac := len(a), len(a[0])
	br, bc := len(b), len(b[0])
	if ac != br {
		panic("linalg: mat mul dim mismatch")
	}
	out := NewMat(ar, bc)
	for i := 0; i < ar; i++ {
		ai := a[i]
		row := out[i]
		for k := 0; k < ac; k++ {
			if ai[k] == 0 {
				continue
			}
			aik := ai[k]
			brow := b[k]
			for j := 0; j < bc; j++ {
				row[j] += aik * brow[j]
			}
		}
	}
	return out
}

func (a Mat) MulVec(v Vec) Vec {
	if len(a) == 0 {
		panic("linalg: empty mat mulvec")
	}
	ac := len(a[0])
	if ac != len(v) {
		panic("linalg: mat mulvec dim mismatch")
	}
	out := NewVec(len(a))
	for i := range a {
		var s float64
		ai := a[i]
		for k := 0; k < ac; k++ {
			s += ai[k] * v[k]
		}
		out[i] = s
	}
	return out
}

func (a Mat) T() Mat {
	if len(a) == 0 {
		return NewMat(0, 0)
	}
	r, c := len(a), len(a[0])
	if r != c {
		out := NewMat(c, r)
		for i := 0; i < r; i++ {
			for j := 0; j < c; j++ {
				out[j][i] = a[i][j]
			}
		}
		return out
	}
	for i := 0; i < r; i++ {
		for j := i + 1; j < c; j++ {
			a[i][j], a[j][i] = a[j][i], a[i][j]
		}
	}
	return a
}

func (a Mat) Add(b Mat) Mat {
	r, c := len(a), len(a[0])
	if len(b) != r || len(b[0]) != c {
		panic("linalg: mat add dim mismatch")
	}
	out := NewMat(r, c)
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			out[i][j] = a[i][j] + b[i][j]
		}
	}
	return out
}

func Identity(n int) Mat {
	m := NewMat(n, n)
	for i := 0; i < n; i++ {
		m[i][i] = 1
	}
	return m
}

func (dst Mat) SetBlock(r0, c0 int, src Mat) {
	for i := 0; i < len(src); i++ {
		dr := dst[r0+i]
		sr := src[i]
		for j := 0; j < len(sr); j++ {
			dr[c0+j] = sr[j]
		}
	}
}

func (a Mat) Copy() Mat {
	out := NewMat(a.Rows(), a.Cols())
	for i := range a {
		copy(out[i], a[i])
	}
	return out
}

func (a Mat) Fill(v float64) {
	for _, row := range a {
		for j := range row {
			row[j] = v
		}
	}
}
