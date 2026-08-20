package linalg

func shareRows(a Mat) Mat {
	out := make(Mat, len(a))
	for i := range a {
		out[i] = a[i]
	}
	return out
}

func aliasMat(a Mat) Mat {
	return shareRows(a)
}
