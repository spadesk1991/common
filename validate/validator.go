package validate

import (
	"github.com/go-playground/validator"
	"github.com/pkg/errors"

	"reflect"
)

var v *validator.Validate

func init() {
	v = validator.New()
}

func Validate(in interface{}) (err error) {
	return handleError(in, v.Struct(in))
}

func handleError(obj interface{}, err error) (e error) {
	t := reflect.TypeOf(obj)
	if err != nil {
		if es, ok := (err).(validator.ValidationErrors); ok {
			for _, error := range es {
				if f, exists := t.Elem().FieldByName(error.Field()); exists {
					if msg, ok := f.Tag.Lookup("errorMsg"); ok {
						e = errors.New(msg)
						return
					} else {
						e = errors.Errorf(`%s`, error)
						return
					}
				} else {
					e = errors.Errorf(`%s`, error)
					return
				}
			}
		}
	}
	return
}
