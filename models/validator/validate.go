package validator

import (
	"fmt"
	"net/http"

	"github.com/go-playground/form/v4"
	"github.com/go-playground/validator/v10"
)

var (
	formDecoder = form.NewDecoder()
	validate    = validator.New()
)

func bindAndValidate(r *http.Request, dst any) (errs map[string]string, err error) {
	if err = r.ParseForm(); err != nil {
		return nil, err
	}
	if err = formDecoder.Decode(dst, r.PostForm); err != nil {
		return nil, err
	}
	if verr := validate.Struct(dst); verr != nil {
		errs = make(map[string]string)
		for _, fe := range verr.(validator.ValidationErrors) {
			errs[fe.Field()] = fmt.Sprintf("must satisfy %s", fe.Tag())
		}
	}
	return errs, nil
}

func BindAndValidateForm[T any](r *http.Request) (T, map[string]string, error) {
	var form T
	errs, err := bindAndValidate(r, &form)
	return form, errs, err
}
