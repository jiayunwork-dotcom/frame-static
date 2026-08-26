package buckling

import (
	"fmt"
	"math"

	"frame-static/internal/assemble"
	"frame-static/internal/element"
	"frame-static/internal/model"
)

func Radius(A, I float64) (float64, error) {
	if !(A > 0) || math.IsNaN(A) || math.IsInf(A, 0) {
		return 0, fmt.Errorf("buckling: A must be positive")
	}
	if !(I > 0) || math.IsNaN(I) || math.IsInf(I, 0) {
		return 0, fmt.Errorf("buckling: I must be positive")
	}
	return math.Sqrt(I / A), nil
}

func Slenderness(L, i float64) (float64, error) {
	if !(L > 0) || math.IsNaN(L) || math.IsInf(L, 0) {
		return 0, fmt.Errorf("buckling: L must be positive")
	}
	if !(i > 0) || math.IsNaN(i) || math.IsInf(i, 0) {
		return 0, fmt.Errorf("buckling: radius of gyration must be positive")
	}
	return L / i, nil
}

func Euler(E, I, Le float64) (float64, error) {
	if !(E > 0) || math.IsNaN(E) || math.IsInf(E, 0) {
		return 0, fmt.Errorf("buckling: E must be positive")
	}
	if !(I > 0) {
		return 0, fmt.Errorf("buckling: I must be positive")
	}
	if !(Le > 0) {
		return 0, fmt.Errorf("buckling: effective length must be positive")
	}
	return math.Pi * math.Pi * E * I / (Le * Le), nil
}

func EffectiveLength(L, k float64) (float64, error) {
	if !(L > 0) {
		return 0, fmt.Errorf("buckling: L must be positive")
	}
	if !(k > 0) {
		return 0, fmt.Errorf("buckling: K factor must be positive")
	}
	return k * L, nil
}

func Utilization(Ncomp, Ncr float64) (float64, error) {
	if !(Ncr > 0) {
		return 0, fmt.Errorf("buckling: Ncr must be positive")
	}
	if Ncomp < 0 {
		return 0, fmt.Errorf("buckling: compression magnitude must be >= 0")
	}
	return Ncomp / Ncr, nil
}

func Compression(Ni, Nj float64) float64 {
	n := 0.5 * (Ni + Nj)
	if n < 0 {
		return -n
	}
	return 0
}

type Check struct {
	From   string
	To     string
	Length float64
	I      float64
	A      float64
	Ncr    float64
	Ncomp  float64
	Util   float64
}

func Member(el model.Element, n1, n2 model.Node, Ni, Nj, kFactor float64) (Check, error) {
	g, err := element.GeometryOf(n1, n2)
	if err != nil {
		return Check{}, err
	}
	le, err := EffectiveLength(g.Length, kFactor)
	if err != nil {
		return Check{}, err
	}
	ncr, err := Euler(el.E, el.I, le)
	if err != nil {
		return Check{}, err
	}
	ncomp := Compression(Ni, Nj)
	u, err := Utilization(ncomp, ncr)
	if err != nil {
		return Check{}, err
	}
	return Check{
		From:   el.From,
		To:     el.To,
		Length: g.Length,
		I:      el.I,
		A:      el.A,
		Ncr:    ncr,
		Ncomp:  ncomp,
		Util:   u,
	}, nil
}

func Frame(m *model.Model, res *assemble.Result, kFactor float64) ([]Check, error) {
	if m == nil || res == nil {
		return nil, fmt.Errorf("buckling: nil model or result")
	}
	index := map[string]model.Node{}
	for _, n := range m.Nodes {
		index[n.ID] = n
	}
	force := map[string]assemble.MemberEndForce{}
	for _, mem := range res.Members {
		force[mem.From+"|"+mem.To] = mem
	}
	out := make([]Check, 0, len(m.Elements))
	for _, el := range m.Elements {
		n1, ok1 := index[el.From]
		n2, ok2 := index[el.To]
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("buckling: missing node for %s-%s", el.From, el.To)
		}
		mem, ok := force[el.From+"|"+el.To]
		if !ok {
			return nil, fmt.Errorf("buckling: missing member force for %s-%s", el.From, el.To)
		}
		c, err := Member(el, n1, n2, mem.Ni, mem.Nj, kFactor)
		if err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, nil
}

func Worst(checks []Check) (Check, bool) {
	var best Check
	found := false
	for _, c := range checks {
		if !found || c.Util > best.Util {
			best = c
			found = true
		}
	}
	return best, found
}

func PinnedK() float64      { return 1.0 }
func FixedK() float64       { return 0.5 }
func CantileverK() float64  { return 2.0 }
func FixedPinnedK() float64 { return 0.7 }
