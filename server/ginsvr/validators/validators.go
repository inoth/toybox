package validators

import (
	"strings"

	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

var (
	trans ut.Translator
)

type Validator interface {
	Tag() string
	Validation() validator.Func
	RegisterTranslations() validator.RegisterTranslationsFunc
	Translation() validator.TranslationFunc
}

type RegisterTranslationsFunc func(string) validator.RegisterTranslationsFunc
type TranslationFunc func(string) validator.TranslationFunc

type validators struct {
	tag string

	validFunc   validator.Func
	regTranFunc validator.RegisterTranslationsFunc
	tranFunc    validator.TranslationFunc
}

func NewValidator(tag string, validFunc validator.Func, regTranFunc validator.RegisterTranslationsFunc, tranFunc validator.TranslationFunc) Validator {
	return &validators{
		tag:         tag,
		validFunc:   validFunc,
		regTranFunc: regTranFunc,
		tranFunc:    tranFunc,
	}
}

func (v validators) Tag() string {
	return v.tag
}

func (v *validators) Validation() validator.Func {
	return v.validFunc
}

func (v *validators) RegisterTranslations() validator.RegisterTranslationsFunc {
	return v.regTranFunc
}

func (v *validators) Translation() validator.TranslationFunc {
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
