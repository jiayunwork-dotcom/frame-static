package element

import (
	"math"
	"testing"
)

func TestFixedEndForcesSigns(t *testing.T) {
	g := &Geometry{Length: 6, C: 1, S: 0}
	fef := FixedEndForces(g, 10)
	if math.Abs(fef[1]-30) > 1e-9 || math.Abs(fef[4]-30) > 1e-9 {
		t.Fatalf("shear terms = %v, %v; want 30,30", fef[1], fef[4])
	}
	if math.Abs(fef[2]-30) > 1e-9 {
		t.Fatalf("moment i = %v, want 30", fef[2])
	}
	if math.Abs(fef[5]+30) > 1e-9 {
		t.Fatalf("moment j = %v, want -30", fef[5])
	}
}

func TestEquivalentNodalLoadIsNegatedFEF(t *testing.T) {
	g := &Geometry{Length: 4, C: 1, S: 0}
	fef := FixedEndForces(g, 5)
	eq := EquivalentNodalLoad(g, 5)
	for i := 0; i < 6; i++ {
		if math.Abs(eq[i]+fef[i]) > 1e-9 {
			t.Fatalf("eq[%d]=%v not -fef[%d]=%v", i, eq[i], i, -fef[i])
		}
	}
}

func TestThermalSelfStressBalances(t *testing.T) {
	g := &Geometry{Length: 5, C: 1, S: 0}
	fef, err := ThermalFixedEnd(g, 2.1e11, 0.01, 1.2e-5, 40)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(fef[0]+fef[3]) > 1e-9 {
		t.Fatalf("thermal axial pair %v %v must cancel", fef[0], fef[3])
	}
	if fef[0] >= 0 {
		t.Fatalf("heating a restrained bar should put compression at i, got %v", fef[0])
	}
	eq, err := ThermalEquivalentNodal(g, 2.1e11, 0.01, 1.2e-5, 40)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(eq[0]+fef[0]) > 1e-9 {
		t.Fatalf("equivalent nodal not -FEF")
	}
	free, err := FreeThermalElongation(5, 1.2e-5, 40)
	if err != nil {
		t.Fatal(err)
	}
	if free <= 0 {
		t.Fatalf("heating should elongate, got %v", free)
	}
}
