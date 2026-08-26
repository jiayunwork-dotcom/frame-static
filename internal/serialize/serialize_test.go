package serialize

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"frame-static/internal/assemble"
	"frame-static/internal/model"
)

func sampleResult(t *testing.T) *assemble.Result {
	t.Helper()
	m := &model.Model{
		Nodes: []model.Node{
			{ID: "A", X: 0, Y: 0, Support: "pin"},
			{ID: "B", X: 4, Y: 0, Restraint: [3]bool{false, true, false}},
		},
		Elements: []model.Element{{From: "A", To: "B", E: 2.1e11, A: 0.01, I: 4e-5}},
		Loads:    []model.NodeLoad{{Node: "B", FY: -500}},
	}
	res, err := assemble.Solve(m)
	if err != nil {
		t.Fatalf("solve: %v", err)
	}
	return res
}

func TestMembersCSVHeader(t *testing.T) {
	b, err := MembersCSV(sampleResult(t))
	if err != nil {
		t.Fatalf("csv: %v", err)
	}
	if !strings.HasPrefix(string(b), "from,to,Ni") {
		t.Fatalf("header missing: %s", b)
	}
}

func TestReactionsCSV(t *testing.T) {
	b, err := ReactionsCSV(sampleResult(t))
	if err != nil {
		t.Fatalf("csv: %v", err)
	}
	if !strings.Contains(string(b), "A,uy") {
		t.Fatalf("reaction row missing: %s", b)
	}
}

func TestJSONRoundTrip(t *testing.T) {
	res := sampleResult(t)
	b, err := ToJSON(res)
	if err != nil {
		t.Fatalf("json: %v", err)
	}
	back, err := FromJSON(b)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(back.Nodes) != len(res.Nodes) || len(back.Reactions) != len(res.Reactions) {
		t.Fatalf("round trip lost data")
	}
}

func TestSnapshotTruncatedKeepsCommittedNodes(t *testing.T) {
	res := sampleResult(t)
	dir := t.TempDir()
	full := filepath.Join(dir, "full.snap")
	if err := WriteSnapshot(full, res); err != nil {
		t.Fatal(err)
	}
	got, err := ReadSnapshot(full)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Nodes) != len(res.Nodes) {
		t.Fatalf("round trip nodes %d", len(got.Nodes))
	}
	raw, err := os.ReadFile(full)
	if err != nil {
		t.Fatal(err)
	}
	if len(raw) < 40 {
		t.Fatalf("snapshot too small")
	}
	cut := raw[:len(raw)-20]
	bad := filepath.Join(dir, "trunc.snap")
	if err := os.WriteFile(bad, cut, 0o644); err != nil {
		t.Fatal(err)
	}
	part, err := ReadSnapshot(bad)
	if err != nil {
		t.Fatalf("truncated snapshot should keep prefix: %v", err)
	}
	if len(part.Nodes) == 0 {
		t.Fatal("prefix nodes vanished")
	}
	if len(part.Members) >= len(res.Members) && len(part.Nodes) == len(res.Nodes) {
		t.Fatal("truncated snapshot still looks complete")
	}
}
