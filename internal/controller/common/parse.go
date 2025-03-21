package common

import (
	"net/http"
)

func DecodeAndValidate(r *http.Request, dst any) (err error) {
	if err = fastjson.NewDecoder(r.Body).Decode(dst); err != nil {
		return err
	}
	if err = Validate(dst); err != nil {
		return err
	}
	return nil
}

func Validate(dst any) (err error) {
	if err = validator.Struct(dst); err != nil {
		return err
	}
	return nil
}
