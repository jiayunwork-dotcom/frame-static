package serialize

import (
	"bytes"
	"encoding/csv"
	"strconv"

	"frame-static/internal/assemble"
)

func f(v float64) string { return strconv.FormatFloat(v, 'f', 6, 64) }

func MembersCSV(res *assemble.Result) ([]byte, error) {
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	if err := w.Write([]string{"from", "to", "Ni", "Vi", "Mi", "Nj", "Vj", "Mj", "length"}); err != nil {
		return nil, err
	}
	for _, m := range res.Members {
		if err := w.Write([]string{m.From, m.To, f(m.Ni), f(m.Vi), f(m.Mi), f(m.Nj), f(m.Vj), f(m.Mj), f(m.Length)}); err != nil {
			return nil, err
		}
	}
	w.Flush()
	return buf.Bytes(), w.Error()
}

func ReactionsCSV(res *assemble.Result) ([]byte, error) {
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	if err := w.Write([]string{"node", "dof", "force"}); err != nil {
		return nil, err
	}
	for _, r := range res.Reactions {
		if err := w.Write([]string{r.Node, r.DOF, f(r.Force)}); err != nil {
			return nil, err
		}
	}
	w.Flush()
	return buf.Bytes(), w.Error()
}

func NodesCSV(res *assemble.Result) ([]byte, error) {
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	if err := w.Write([]string{"id", "ux", "uy", "theta"}); err != nil {
		return nil, err
	}
	for _, n := range res.Nodes {
		if err := w.Write([]string{n.ID, f(n.UX), f(n.UY), f(n.Theta)}); err != nil {
			return nil, err
		}
	}
	w.Flush()
	return buf.Bytes(), w.Error()
}
