package assemble

import "frame-static/internal/linalg"

func Reactions(K linalg.Mat, u, F linalg.Vec) linalg.Vec {
	return K.MulVec(u).Sub(F)
}
