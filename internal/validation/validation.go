// Package validation turns invalid input into a field -> problem map
// that the web app can show under each form field.
package validation

import (
	"errors"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Error struct {
	Fields map[string]string
}

func (e *Error) Error() string { return "validation failed" }

func Setup() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(f reflect.StructField) string {
			name, _, _ := strings.Cut(f.Tag.Get("json"), ",")
			return name
		})
	}
}

func Details(err error) (fields map[string]string, ok bool) {
	var ve validator.ValidationErrors
	if !errors.As(err, &ve) {
		return nil, false
	}

	fields = make(map[string]string, len(ve))
	for _, fe := range ve {
		fields[fe.Field()] = message(fe)
	}

	return fields, true
}

func message(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "required"
	case "max":
		return "must be at most " + fe.Param() + " characters"
	default:
		return "invalid"
	}
}
