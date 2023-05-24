package res

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 自定义错误码
const (
	SuccessCode       int = 100 * iota // 成功返回
	UndefErrorCode                     // 100 未定义错误
	ValidErrorCode                     // 200 自定义错误
	InternalErrorCode                  // 300 内部错误
	ParamErrorCode                     // 400 参数错误

	Unauthorized            = 401 //  权限错误
	InvalidRequestErrorCode = 403 // 无效请求

	StatusInternalServerError = 500
)

type Result struct {
	ErrorCode int         `json:"err_code"`
	ErrorMsg  string      `json:"err_msg"`
	Data      interface{} `json:"data"`
}

func (r Result) String() string {
	buf, _ := json.Marshal(r)
	return string(buf)
}

func OK(c *gin.Context, msg string, data ...interface{}) {
	r := Result{
		ErrorCode: SuccessCode,
		ErrorMsg:  msg,
	}
	if len(data) > 0 {
		r.Data = data[0]
	}
	c.JSON(http.StatusOK, r)
}

func Err(c *gin.Context, code int, err error) {
	r := Result{
		ErrorCode: code,
		ErrorMsg:  err.Error(),
	}
	switch code {
	case StatusInternalServerError:
		c.JSON(StatusInternalServerError, r)
		c.AbortWithError(StatusInternalServerError, err)
	default:
		c.JSON(http.StatusOK, r)
		c.AbortWithError(http.StatusOK, err)
	}
}
