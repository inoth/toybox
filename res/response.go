package res

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

// 自定义错误码
const (
	SuccessCode       int = 100 * iota // 成功返回
	UndefErrorCode                     // 未定义错误
	ValidErrorCode                     // 自定义错误
	InternalErrorCode                  // 内部错误
	ParamErrorCode                     // 参数错误

	InvalidRequestErrorCode = 401 // 无效请求
)

type Result struct {
	ErrorCode int         `json:"code"`
	ErrorMsg  string      `json:"msg"`
	Data      interface{} `json:"data"`
	TraceId   interface{} `json:"trace_id"`
}

func (r Result) String() []byte {
	buf, _ := json.Marshal(r)
	return buf
}

func ResultErr(c *gin.Context, code int, err error) {
	traceId, _ := c.Get("trace_id")
	resp := &Result{ErrorCode: code, ErrorMsg: err.Error(), Data: "", TraceId: traceId}
	c.Set("result", string(resp.String()))

	c.JSON(200, resp)
	c.AbortWithError(200, err)
}

func ResultOk(c *gin.Context, code int, data ...interface{}) {
	traceId, _ := c.Get("trace_id")
	var res interface{}
	if len(data) > 0 {
		res = data[0]
	}
	resp := &Result{ErrorCode: code, ErrorMsg: "ok", Data: res, TraceId: traceId}
	c.Set("result", string(resp.String()))

	c.JSON(200, resp)
}
