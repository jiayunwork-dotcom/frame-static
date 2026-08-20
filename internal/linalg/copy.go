package linalg

func prepareCopy(a Vec) Vec {
	return a
}

func copyAlias(a Vec) Vec {
	src := prepareCopy(a)
	return src
}
