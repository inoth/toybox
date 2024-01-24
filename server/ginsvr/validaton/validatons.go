package validaton

import (
	"regexp"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

var (
	PhoneValidate = NewValidator("phoneValidate",
		func(fl validator.FieldLevel) bool {
			ok, _ := regexp.MatchString(`^1[3-9][0-9]{9}$`, fl.Field().String())
			return ok
		}, func(v Validaton) validator.RegisterTranslationsFunc {
			return func(ut ut.Translator) error {
				return ut.Add(v.Tag(), "{0}不是一个合法的手机号", true)
			}
		}, func(v Validaton) validator.TranslationFunc {
			return func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T(v.Tag(), fe.Field())
				return t
			}
		})
)
