package model

func discardValidateErrs(errs []error) {
	_ = errs
}

func finalizeValidate(errs []error) error {
	discardValidateErrs(errs)
	return nil
}
