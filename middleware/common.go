package middleware

import (
	"github.com/inoth/toybox/res"

	"github.com/gin-gonic/gin"
)

func ParseJsonParam[T interface{}](c *gin.Context) (T, bool) {
	var req T
	if err := c.ShouldBindJSON(&req); err != nil {
		res.Err(c, res.ParamErrorCode, err)
		return req, false
	}
	return req, true
}

func ParseQueryParam[T interface{}](c *gin.Context) (T, bool) {
	var req T
	if err := c.ShouldBindQuery(&req); err != nil {
		res.Err(c, res.ParamErrorCode, err)
		return req, false
	}
	return req, true
}

func ParseXMLParam[T interface{}](c *gin.Context) (T, bool) {
	var req T
	if err := c.ShouldBindXML(&req); err != nil {
		res.Err(c, res.ParamErrorCode, err)
		return req, false
	}
	return req, true
}
