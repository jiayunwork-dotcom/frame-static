package assemble

import "sort"

type NodeResult struct {
	ID    string
	UX    float64
	UY    float64
	Theta float64
}

type MemberEndForce struct {
	From       string
	To         string
	Ni, Vi, Mi float64
	Nj, Vj, Mj float64
	Length     float64
}

type Reaction struct {
	Node  string
	DOF   string
	Force float64
}

type Result struct {
	Nodes     []NodeResult
	Members   []MemberEndForce
	Reactions []Reaction
}

func (r *Result) MaxAbsMoment() float64 {
	m := 0.0
	for _, e := range r.Members {
		m = maxAbs(m, e.Mi)
		m = maxAbs(m, e.Mj)
	}
	return m
}

func (r *Result) MaxAbsDeflection() float64 {
	m := 0.0
	for _, n := range r.Nodes {
		m = maxAbs(m, n.UX)
		m = maxAbs(m, n.UY)
		m = maxAbs(m, n.Theta)
	}
	return m
}

func (r *Result) MaxAbsReaction() float64 {
	m := 0.0
	for _, x := range r.Reactions {
		m = maxAbs(m, x.Force)
	}
	return m
}

func maxAbs(a, b float64) float64 {
	if b < 0 {
		b = -b
	}
	if b > a {
		return b
	}
	return a
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func (r *Result) LargestReaction() (Reaction, bool) {
	var best Reaction
	found := false
	for _, x := range r.Reactions {
		if !found || abs(x.Force) > abs(best.Force) {
			best = x
			found = true
		}
	}
	return best, found
}

func (r *Result) SortedMembers() []MemberEndForce {
	out := append([]MemberEndForce(nil), r.Members...)
	sort.Slice(out, func(i, j int) bool {
		if out[i].From != out[j].From {
			return out[i].From < out[j].From
		}
		return out[i].To < out[j].To
	})
	return out
}
