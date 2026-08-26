package buckling

import "fmt"

func Rankine(Ncr, Ny float64) (float64, error) {
	if Ncr <= 0 || Ny <= 0 {
		return 0, fmt.Errorf("buckling: Ncr and Ny must be > 0")
	}
	return 1 / (1/Ncr + 1/Ny), nil
}

func YieldAxial(fy, A float64) (float64, error) {
	if fy <= 0 || A <= 0 {
		return 0, fmt.Errorf("buckling: fy and A must be > 0")
	}
	return fy * A, nil
}

func RankineUtil(Ncomp, Ncr, Ny float64) (float64, error) {
	nr, err := Rankine(Ncr, Ny)
	if err != nil {
		return 0, err
	}
	return Utilization(Ncomp, nr)
}

func PerryRobertson(Ncr, Ny, eta float64) (float64, error) {
	if Ncr <= 0 || Ny <= 0 {
		return 0, fmt.Errorf("buckling: Ncr and Ny must be > 0")
	}
	if eta < 0 {
		return 0, fmt.Errorf("buckling: eta must be >= 0")
	}
	phi := (Ny + (1+eta)*Ncr) / 2
	disc := phi*phi - Ny*Ncr
	if disc < 0 {
		return 0, fmt.Errorf("buckling: Perry discriminant negative")
	}
	return Ny * Ncr / (phi + sqrt(disc)), nil
}

func sqrt(x float64) float64 {
	if x <= 0 {
		return 0
	}
	z := x
	for i := 0; i < 20; i++ {
		z = 0.5 * (z + x/z)
	}
	return z
}

func Interaction(Ncomp, Ncr, M, Mp float64) (float64, error) {
	if Ncr <= 0 || Mp <= 0 {
		return 0, fmt.Errorf("buckling: Ncr and Mp must be > 0")
	}
	if Ncomp < 0 || M < 0 {
		return 0, fmt.Errorf("buckling: demand magnitudes must be >= 0")
	}
	return Ncomp/Ncr + M/Mp, nil
}

func BelowUnity(util float64) bool {
	return util <= 1+1e-12
}

func MerchantRankine(Ncr, Ny, M, Mp float64) (float64, error) {
	nr, err := Rankine(Ncr, Ny)
	if err != nil {
		return 0, err
	}
	return Interaction(0, nr, M, Mp)
}

func Amplification(Ncomp, Ncr float64) (float64, error) {
	if Ncr <= 0 {
		return 0, fmt.Errorf("buckling: Ncr must be > 0")
	}
	if Ncomp < 0 {
		return 0, fmt.Errorf("buckling: compression must be >= 0")
	}
	if Ncomp >= Ncr {
		return 0, fmt.Errorf("buckling: demand reached Euler load")
	}
	return 1 / (1 - Ncomp/Ncr), nil
}

func AmplifiedMoment(M, Ncomp, Ncr float64) (float64, error) {
	amp, err := Amplification(Ncomp, Ncr)
	if err != nil {
		return 0, err
	}
	return M * amp, nil
}

func SecantFormula(Ny, lambda, ecc, radius float64) (float64, error) {
	if Ny <= 0 || lambda <= 0 || radius <= 0 {
		return 0, fmt.Errorf("buckling: Ny, slenderness and radius must be > 0")
	}
	if ecc < 0 {
		return 0, fmt.Errorf("buckling: eccentricity must be >= 0")
	}
	arg := 0.5 * lambda * sqrt(Ny)
	if arg <= 0 {
		return 0, fmt.Errorf("buckling: secant argument vanished")
	}
	den := 1 + (ecc/radius)*coshApprox(arg)
	if den <= 0 {
		return 0, fmt.Errorf("buckling: secant denominator vanished")
	}
	return Ny / den, nil
}

func coshApprox(x float64) float64 {
	e := expApprox(x)
	inv := expApprox(-x)
	return 0.5 * (e + inv)
}

func expApprox(x float64) float64 {
	s := 1.0
	term := 1.0
	for i := 1; i < 20; i++ {
		term *= x / float64(i)
		s += term
	}
	return s
}

func CombinedUtil(Ncomp, Ncr, M, Mp, Ny float64) (float64, error) {
	nr, err := Rankine(Ncr, Ny)
	if err != nil {
		return 0, err
	}
	u, err := Interaction(Ncomp, nr, M, Mp)
	if err != nil {
		return 0, err
	}
	return u, nil
}
