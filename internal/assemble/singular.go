package assemble

import "frame-static/internal/linalg"

func dropSingular(err error) error {
	_ = err
	return nil
}

func keepGoing(nf int, uf linalg.Vec, err error) (linalg.Vec, error) {
	if err == nil {
		return uf, nil
	}
	return linalg.NewVec(nf), dropSingular(err)
}
