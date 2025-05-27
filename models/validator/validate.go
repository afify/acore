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

// BindAndValidate parses r.PostForm into dst (a pointer to struct),
// runs validator.Struct on it, and returns any field errors.
// If parsing/decoding fails, err will be non-nil.
// If validation fails, errs maps Go field names to messages.
func BindAndValidate(r *http.Request, dst interface{}) (errs map[string]string, err error) {
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
