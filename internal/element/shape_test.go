package element

import (
	"math"
	"testing"
)

func TestShapeBendingEndValues(t *testing.T) {
	n0 := ShapeBendingAt(4, 0)
	nL := ShapeBendingAt(4, 4)
	if math.Abs(n0[0]-1) > 1e-12 || math.Abs(n0[3]) > 1e-12 {
		t.Fatalf("shape at x=0 = %v", n0)
	}
	if math.Abs(nL[2]-1) > 1e-12 || math.Abs(nL[1]) > 1e-12 {
		t.Fatalf("shape at x=L = %v", nL)
	}
}

func TestDeflectionInterpolation(t *testing.T) {
	v := DeflectionAt(4, 4, 0, 0, 1, 0)
	if math.Abs(v-1) > 1e-12 {
		t.Fatalf("v(L) = %v", v)
	}
}

func TestMaxDeflectionFindsPeak(t *testing.T) {
	maxVal, atX := MaxDeflection(10, 0, 0.2, 0, -0.2, 20)
	if maxVal <= 0 || atX <= 0 || atX >= 10 {
		t.Fatalf("max deflection %v at %v", maxVal, atX)
	}
}
