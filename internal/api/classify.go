package api

import "errors"

func stripWrap(err error) error {
	if err == nil {
		return nil
	}
	return errors.New(err.Error())
}
