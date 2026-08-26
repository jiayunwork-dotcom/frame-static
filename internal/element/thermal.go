package element

import "fmt"

func ThermalAxialForce(E, A, alpha, dT float64) (float64, error) {
	if E <= 0 || A <= 0 {
		return 0, fmt.Errorf("element: E and A must be > 0")
	}
	if alpha < 0 {
		return 0, fmt.Errorf("element: expansion coefficient must be >= 0")
	}
	return E * A * alpha * dT, nil
}

func ThermalFixedEnd(g *Geometry, E, A, alpha, dT float64) ([6]float64, error) {
	var z [6]float64
	if g == nil || g.Length <= 0 {
		return z, fmt.Errorf("element: geometry required for thermal load")
	}
	n, err := ThermalAxialForce(E, A, alpha, dT)
	if err != nil {
		return z, err
	}
	z[0] = -n
	z[3] = n
	return z, nil
}

func ThermalEquivalentNodal(g *Geometry, E, A, alpha, dT float64) ([6]float64, error) {
	fef, err := ThermalFixedEnd(g, E, A, alpha, dT)
	if err != nil {
		return fef, err
	}
	var out [6]float64
	for i := range out {
		out[i] = -fef[i]
	}
	return out, nil
}

func FreeThermalElongation(L, alpha, dT float64) (float64, error) {
	if L <= 0 {
		return 0, fmt.Errorf("element: length must be > 0")
	}
	if alpha < 0 {
		return 0, fmt.Errorf("element: expansion coefficient must be >= 0")
	}
	return L * alpha * dT, nil
}

func RestrainedThermalStress(E, alpha, dT float64) (float64, error) {
	if E <= 0 {
		return 0, fmt.Errorf("element: E must be > 0")
	}
	if alpha < 0 {
		return 0, fmt.Errorf("element: expansion coefficient must be >= 0")
	}
	return E * alpha * dT, nil
}

func GradientCurvature(alpha, dTface, h float64) (float64, error) {
	if h <= 0 {
		return 0, fmt.Errorf("element: section depth must be > 0")
	}
	if alpha < 0 {
		return 0, fmt.Errorf("element: expansion coefficient must be >= 0")
	}
	return alpha * dTface / h, nil
}

func ThermalMoment(E, I, alpha, dTface, h float64) (float64, error) {
	phi, err := GradientCurvature(alpha, dTface, h)
	if err != nil {
		return 0, err
	}
	if E <= 0 || I <= 0 {
		return 0, fmt.Errorf("element: E and I must be > 0")
	}
	return E * I * phi, nil
}

func ThermalPairCancels(fef [6]float64) bool {
	return fef[0]+fef[3] == 0 && fef[1] == 0 && fef[4] == 0
}

func HeatingCompresses(fef [6]float64) bool {
	return fef[0] < 0 && fef[3] > 0
}
