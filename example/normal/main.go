package main

import (
	"demo1/router"
	"regexp"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/inoth/toybox"
	"github.com/inoth/toybox/component/logger"
	"github.com/inoth/toybox/server/ginsvr"
	"github.com/inoth/toybox/server/ginsvr/httpgin"
	"github.com/inoth/toybox/server/ginsvr/validators"
)

func main() {
	tb := toybox.New(
		toybox.WithLoadConf(),
		logger.WithLogger(),
		httpgin.NewHttpGin(func(hgs *httpgin.HttpGinServer) {
			hgs.Use(ginsvr.Recovery())
		}, router.WithUserRouter(),
			httpgin.WithValidator(validators.NewValidator("phoneValidate", func(fl validator.FieldLevel) bool {
				ok, _ := regexp.MatchString(`^1[3-9][0-9]{9}$`, fl.Field().String())
				return ok
			}, func(ut ut.Translator) error {
				return ut.Add("phoneValidate", "{0}不是一个合法的手机号", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("phoneValidate", fe.Field())
				return t
			})),
		),
	)
	if err := tb.Run(); err != nil {
		panic(err)
	}
}
