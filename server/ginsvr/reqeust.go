package ginsvr

import (
	"github.com/gin-gonic/gin"
	"github.com/inoth/toybox/server/ginsvr/res"
)

func ParseJsonParam[T interface{}](c *gin.Context) (T, bool) {
	var req T
	if err := c.ShouldBindJSON(&req); err != nil {
		res.ErrParams(c, err.Error())
		return req, false
	}
	return req, true
}

func ParseQueryParam[T interface{}](c *gin.Context) (T, bool) {
	var req T
	if err := c.ShouldBindQuery(&req); err != nil {
		res.ErrParams(c, err.Error())
		return req, false
	}
	return req, true
}

func ParseXMLParam[T interface{}](c *gin.Context) (T, bool) {
	var req T
	if err := c.ShouldBindXML(&req); err != nil {
		res.ErrParams(c, err.Error())
		return req, false
	}
	return req, true
}
