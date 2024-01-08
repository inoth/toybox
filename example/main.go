package main

import (
	"github/inoth/toybox"
	"github/inoth/toybox/server/httpgin"

	"github.com/gin-gonic/gin"
)

func main() {
	tb := toybox.New(toybox.WithLoadConf(), httpgin.NewHttpGin(func(hgs *httpgin.HttpGinServer) {
		hgs.GET("test", func(c *gin.Context) {
			data, _, _ := hgs.Do("test", func() (interface{}, error) {
				return "123", nil
			})
			c.String(200, data.(string))
		})
	}))

	// aaa := toybox.GetString("aaaa")
	// fmt.Println(aaa)
	// list := toybox.GetStringSlice("list")
	// for _, v := range list {
	// 	fmt.Println(v)
	// }

	// mysql.GetDB(mysql.SetConfig(toybox.GetConfMate()))

	// mysql.GetDB(mysql.SetName("mysql2"), mysql.SetConfig(toybox.GetConfMate()))

	tb.Run()
}
