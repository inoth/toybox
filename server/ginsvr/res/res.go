package res

import (
	"net/http"

	"github.com/inoth/toybox/util"

	"github.com/gin-gonic/gin"
)

const (
	Success             = iota * 100
	InvalidParameter          // 参数错误
	MissingParameter          // 参数缺损
	InvalidToken              // token 失效
	BadRequest                // 无效的访问
	Unauthorized        = 401 // 未授权
	NotFound            = 404 // 空数据
	InternalServerError = 500 // 内部服务器错误
)

type ResultBody struct {
	TraceId string      `json:"trace_id"`
	Ret     int         `json:"ret"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data,omitempty"`
}

func (rb *ResultBody) result(c *gin.Context) {
	c.JSON(http.StatusOK, rb)
}

func Ok(c *gin.Context, msg string, data ...any) {
	rb := ResultBody{
		TraceId: c.GetHeader("TraceId"),
		Ret:     Success,
		Msg:     msg,
		Data:    util.First(nil, data),
	}
	rb.result(c)
}

func Result(c *gin.Context, ret int, msg string, data ...any) {
	rb := ResultBody{
		TraceId: c.GetHeader("TraceId"),
		Ret:     ret,
		Msg:     msg,
		Data:    util.First(nil, data),
	}
	rb.result(c)
}
