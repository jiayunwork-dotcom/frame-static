package element

import (
	"fmt"

	"frame-static/internal/linalg"
)

func GeometricStiffness(g *Geometry, N float64) (linalg.Mat, error) {
	if g == nil || g.Length <= 0 {
		return nil, fmt.Errorf("element: geometry required for geometric stiffness")
	}
	L := g.Length
	k := linalg.NewMat(6, 6)
	a := N / L
	k[1][1] = 6.0 / 5.0 * a
	k[1][2] = L / 10.0 * a
	k[1][4] = -6.0 / 5.0 * a
	k[1][5] = L / 10.0 * a
	k[2][1] = L / 10.0 * a
	k[2][2] = 2.0 * L * L / 15.0 * a
	k[2][4] = -L / 10.0 * a
	k[2][5] = -L * L / 30.0 * a
	k[4][1] = -6.0 / 5.0 * a
	k[4][2] = -L / 10.0 * a
	k[4][4] = 6.0 / 5.0 * a
	k[4][5] = -L / 10.0 * a
	k[5][1] = L / 10.0 * a
	k[5][2] = -L * L / 30.0 * a
	k[5][4] = -L / 10.0 * a
	k[5][5] = 2.0 * L * L / 15.0 * a
	return k, nil
}

func CombinedLocal(g *Geometry, E, A, I, N float64) (linalg.Mat, error) {
	kg, err := GeometricStiffness(g, N)
	if err != nil {
		return nil, err
	}
	return LocalStiffness(g, E, A, I).Add(kg), nil
}

func LateralStiffness(g *Geometry, E, I, N float64) (float64, error) {
	if g == nil || g.Length <= 0 {
		return 0, fmt.Errorf("element: geometry required")
	}
	elastic := 12 * E * I / (g.Length * g.Length * g.Length)
	geo := 6.0 / 5.0 * N / g.Length
	return elastic + geo, nil
}

func SoftensUnderCompression(g *Geometry, E, I, Ncomp float64) (float64, error) {
	if Ncomp < 0 {
		return 0, fmt.Errorf("element: compression magnitude must be >= 0")
	}
	zero, err := LateralStiffness(g, E, I, 0)
	if err != nil {
		return 0, err
	}
	comp, err := LateralStiffness(g, E, I, -Ncomp)
	if err != nil {
		return 0, err
	}
	if zero <= 0 {
		return 0, fmt.Errorf("element: elastic lateral stiffness vanished")
	}
	return comp / zero, nil
}

func TensionStiffens(g *Geometry, E, I, Ntens float64) (float64, error) {
	if Ntens < 0 {
		return 0, fmt.Errorf("element: tension must be >= 0")
	}
	zero, err := LateralStiffness(g, E, I, 0)
	if err != nil {
		return 0, err
	}
	ten, err := LateralStiffness(g, E, I, Ntens)
	if err != nil {
		return 0, err
	}
	if zero <= 0 {
		return 0, fmt.Errorf("element: elastic lateral stiffness vanished")
	}
	return ten / zero, nil
}
