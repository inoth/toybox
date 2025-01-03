package res

import (
	"github.com/inoth/toybox/util"

	"github.com/gin-gonic/gin"
)

func Failed(c *gin.Context, msg ...string) {
	rb := ResultBody{
		TraceId: c.GetHeader("trace_id"),
		Ret:     InternalServerError,
		Msg:     util.First("InternalServerError", msg),
	}
	rb.result(c)
}

func ErrParams(c *gin.Context, msg ...string) {
	rb := ResultBody{
		TraceId: c.GetHeader("trace_id"),
		Ret:     InvalidParameter,
		Msg:     util.First("InvalidParameter", msg),
	}
	rb.result(c)
}

func ErrParamsWithErr(c *gin.Context, data ...any) {
	rb := ResultBody{
		TraceId: c.GetHeader("trace_id"),
		Ret:     InvalidParameter,
		Msg:     "InvalidParameter",
		Data:    util.First(nil, data),
	}
	rb.result(c)
}

func ErrParamsMissing(c *gin.Context, msg ...string) {
	rb := ResultBody{
		TraceId: c.GetHeader("trace_id"),
		Ret:     MissingParameter,
		Msg:     util.First("MissingParameter", msg),
	}
	rb.result(c)
}

func ErrToken(c *gin.Context, msg ...string) {
	rb := ResultBody{
		TraceId: c.GetHeader("trace_id"),
		Ret:     InvalidToken,
		Msg:     util.First("InvalidToken", msg),
	}
	rb.result(c)
}

func ErrBadRequest(c *gin.Context, msg ...string) {
	rb := ResultBody{
		TraceId: c.GetHeader("trace_id"),
		Ret:     BadRequest,
		Msg:     util.First("BadRequest", msg),
	}
	rb.result(c)
}

func ErrAuth(c *gin.Context, msg ...string) {
	rb := ResultBody{
		TraceId: c.GetHeader("trace_id"),
		Ret:     Unauthorized,
		Msg:     util.First("Unauthorized", msg),
	}
	rb.result(c)
}

func ErrNotFound(c *gin.Context, msg ...string) {
	rb := ResultBody{
		TraceId: c.GetHeader("trace_id"),
		Ret:     NotFound,
		Msg:     util.First("NotFound", msg),
	}
	rb.result(c)
}
