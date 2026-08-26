package model

func Validate(m *Model) error {
	errs := ValidateAll(m)
	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}
