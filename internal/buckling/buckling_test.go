package buckling

import (
	"math"
	"testing"

	"frame-static/internal/assemble"
	"frame-static/internal/model"
)

func TestRadiusAndEulerScale(t *testing.T) {
	i, err := Radius(0.01, 8e-5)
	if err != nil {
		t.Fatal(err)
	}
	want := math.Sqrt(8e-5 / 0.01)
	if math.Abs(i-want) > 1e-15 {
		t.Fatalf("radius %g want %g", i, want)
	}
	n1, err := Euler(2.1e11, 8e-5, 4)
	if err != nil {
		t.Fatal(err)
	}
	n2, err := Euler(2.1e11, 8e-5, 2)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(n2/n1-4) > 1e-9 {
		t.Fatalf("halving L should *4 Ncr, got %g / %g", n2, n1)
	}
}

func TestFixedColumnHasFourTimesPinnedCapacity(t *testing.T) {
	L := 4.0
	E, I := 2.1e11, 8e-5
	leP, err := EffectiveLength(L, PinnedK())
	if err != nil {
		t.Fatal(err)
	}
	leF, err := EffectiveLength(L, FixedK())
	if err != nil {
		t.Fatal(err)
	}
	np, err := Euler(E, I, leP)
	if err != nil {
		t.Fatal(err)
	}
	nf, err := Euler(E, I, leF)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(nf/np-4) > 1e-9 {
		t.Fatalf("fixed/pinned %g", nf/np)
	}
}

func TestPortalCompressionBelowEuler(t *testing.T) {
	f, err := osOpenPortal(t)
	if err != nil {
		t.Fatal(err)
	}
	res, err := assemble.Solve(f)
	if err != nil {
		t.Fatal(err)
	}
	checks, err := Frame(f, res, PinnedK())
	if err != nil {
		t.Fatal(err)
	}
	if len(checks) != 3 {
		t.Fatalf("members %d", len(checks))
	}
	w, ok := Worst(checks)
	if !ok {
		t.Fatal("no checks")
	}
	if w.Util >= 1 {
		t.Fatalf("portal should stay below Euler, util=%g on %s-%s", w.Util, w.From, w.To)
	}
}

func TestRankineBelowBothLimits(t *testing.T) {
	nr, err := Rankine(1e6, 2e6)
	if err != nil {
		t.Fatal(err)
	}
	if nr >= 1e6 || nr >= 2e6 {
		t.Fatalf("Rankine %v should sit below both Ncr and Ny", nr)
	}
	u, err := RankineUtil(5e5, 1e6, 2e6)
	if err != nil {
		t.Fatal(err)
	}
	if !BelowUnity(u) {
		t.Fatalf("half Rankine demand should be below unity, util=%v", u)
	}
	ia, err := Interaction(0, 1e6, 0, 1e5)
	if err != nil {
		t.Fatal(err)
	}
	if ia != 0 {
		t.Fatalf("zero demand interaction %v", ia)
	}
}

func osOpenPortal(t *testing.T) (*model.Model, error) {
	t.Helper()
	return model.ParseModelBytes([]byte(`{
  "nodes": [
    {"id": "A", "x": 0, "y": 0, "support": "fixed"},
    {"id": "B", "x": 0, "y": 4},
    {"id": "C", "x": 6, "y": 4},
    {"id": "D", "x": 6, "y": 0, "support": "fixed"}
  ],
  "elements": [
    {"from": "A", "to": "B", "E": 2.1e11, "A": 0.01, "I": 8e-5},
    {"from": "B", "to": "C", "E": 2.1e11, "A": 0.01, "I": 1.2e-4},
    {"from": "C", "to": "D", "E": 2.1e11, "A": 0.01, "I": 8e-5}
  ],
  "loads": [{"node": "B", "fx": -50}]
}`))
}
