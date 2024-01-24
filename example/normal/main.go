package main

import (
	"demo1/router"

	"github.com/inoth/toybox"
	"github.com/inoth/toybox/component/logger"
	"github.com/inoth/toybox/server/ginsvr"
	"github.com/inoth/toybox/server/ginsvr/httpgin"
	"github.com/inoth/toybox/server/ginsvr/validaton"
)

func main() {
	tb := toybox.New(
		toybox.WithLoadConf(),
		logger.WithLogger(),
		httpgin.NewHttpGin(func(hgs *httpgin.HttpGinServer) {
			hgs.Use(ginsvr.Recovery())
		}, router.WithUserRouter(),
			httpgin.WithValidator(validaton.PhoneValidate),
		),
	)
	if err := tb.Run(); err != nil {
		panic(err)
	}
}

// func(ut ut.Translator) error {
// 	return ut.Add("phoneValidate", "{0}不是一个合法的手机号", true)
// }, func(ut ut.Translator, fe validator.FieldError) string {
// 	t, _ := ut.T("phoneValidate", fe.Field())
// 	return t
// }
