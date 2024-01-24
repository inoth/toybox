package validaton

import (
	"strings"

	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

var (
	trans ut.Translator
)

type Validaton interface {
	Tag() string
	Validation() validator.Func
	RegisterTranslations() validator.RegisterTranslationsFunc
	Translation() validator.TranslationFunc
}

type RegisterTranslationsFunc func(Validaton) validator.RegisterTranslationsFunc
type TranslationFunc func(Validaton) validator.TranslationFunc

type validaton struct {
	tag string

	validFunc   validator.Func
	regTranFunc validator.RegisterTranslationsFunc
	tranFunc    validator.TranslationFunc
}

func NewValidator(tag string, validFunc validator.Func, regTranFunc RegisterTranslationsFunc, tranFunc TranslationFunc) Validaton {
	valid := validaton{
		tag:       tag,
		validFunc: validFunc,
	}
	if regTranFunc != nil {
		valid.regTranFunc = regTranFunc(&valid)
	}
	if tranFunc != nil {
		valid.tranFunc = tranFunc(&valid)
	}
	return &valid
}

func (v validaton) Tag() string {
	return v.tag
}

func (v validaton) Validation() validator.Func {
	return v.validFunc
}

func (v validaton) RegisterTranslations() validator.RegisterTranslationsFunc {
	return v.regTranFunc
}

func (v validaton) Translation() validator.TranslationFunc {
	return v.tranFunc
}

func GetTranslator() ut.Translator {
	if trans == nil {
		translator := zh.New()
		uni := ut.New(translator)
		trans, _ = uni.GetTranslator("zh")
	}
	return trans
}

func ValidatorTranslate(err error) interface{} {
	if errs, ok := err.(validator.ValidationErrors); ok {
		return removeTopStruct(errs.Translate(trans))
	}
	return err.Error()
}

func removeTopStruct(fields map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fields {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}
