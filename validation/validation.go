package validation

import (
	"reflect"
	"strings"
	"sync"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh2 "github.com/go-playground/validator/v10/translations/zh"
)

var (
	once sync.Once
	lock sync.RWMutex

	trans    ut.Translator
	validate *validator.Validate
)

type Validation interface {
	Tag() string
	Validator() validator.Func
	RegisterTranslation() validator.RegisterTranslationsFunc
	Translation() validator.TranslationFunc
}

type validation struct {
	tag         string
	validFunc   validator.Func
	regTranFunc validator.RegisterTranslationsFunc
	tranFunc    validator.TranslationFunc
}

func (v *validation) Tag() string {
	return v.tag
}
func (v *validation) Validator() validator.Func {
	return v.validFunc
}
func (v *validation) RegisterTranslation() validator.RegisterTranslationsFunc {
	return v.regTranFunc
}
func (v *validation) Translation() validator.TranslationFunc {
	return v.tranFunc
}

func NewValidator(tag string,
	validFunc validator.Func,
	regTranFunc func(Validation) validator.RegisterTranslationsFunc,
	tranFunc func(Validation) validator.TranslationFunc,
) Validation {
	valid := validation{
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

func GetTranslator() ut.Translator {
	if trans == nil {
		translator := zh.New()
		uni := ut.New(translator)
		trans, _ = uni.GetTranslator("zh")
	}
	return trans
}

func ValidatorTranslate(err error) any {
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

func GetDefaultValidator() *validator.Validate {
	once.Do(func() {
		if validate == nil {
			validate, _ = binding.Validator.Engine().(*validator.Validate)
		}
	})
	return validate
}

func LoadValidation(valids []Validation) {
	lock.Lock()
	defer lock.Unlock()

	trans := GetTranslator()
	validate := GetDefaultValidator()
	_ = zh2.RegisterDefaultTranslations(validate, trans)
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	for _, valid := range valids {
		if valid.Validator() != nil {
			validate.RegisterValidation(valid.Tag(), valid.Validator())
		}
		if valid.RegisterTranslation() != nil && valid.Translation() != nil {
			validate.RegisterTranslation(valid.Tag(), trans, valid.RegisterTranslation(), valid.Translation())
		}
	}
}
