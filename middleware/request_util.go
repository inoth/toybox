package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/inoth/ino-toybox/res"
)

func RequestJsonParamHandler[T interface{}](c *gin.Context) (T, bool) {
	var req T
	if err := c.ShouldBindJSON(&req); err != nil {
		res.ResultErr(c, res.ParamErrorCode, err)
		return req, false
	}
	return req, true
}

func RequestQueryParamHandler[T interface{}](c *gin.Context) (T, bool) {
	var req T
	if err := c.ShouldBindQuery(&req); err != nil {
		res.ResultErr(c, res.ParamErrorCode, err)
		return req, false
	}
	return req, true
}

func RequestXMLParamHandler[T interface{}](c *gin.Context) (T, bool) {
	var req T
	if err := c.ShouldBindXML(&req); err != nil {
		res.ResultErr(c, res.ParamErrorCode, err)
		return req, false
	}
	return req, true
}
