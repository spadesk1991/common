package validate

import (
	"regexp"
	"time"

	"github.com/go-playground/validator"
	"github.com/pkg/errors"

	"reflect"
)

var v *validator.Validate

func init() {
	v = validator.New()
	v.RegisterValidation("YYYMMDD", CheckYYYYMMDD)
}

func Validate(in interface{}) (err error) {
	return handleError(in, v.Struct(in))
}

func CheckYYYYMMDD(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	f, err := regexp.MatchString("\\w{4}-\\w{2}-\\w{2}", value)
	if err != nil {
		return false
	}
	if !f {
		return false
	}
	_, err = time.Parse("2006-0102", value)
	if err != nil {
		return false
	}
	return true
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
