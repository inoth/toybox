package res

import (
	"github.com/inoth/toybox/util"

	"github.com/gin-gonic/gin"
)

func Failed(c *gin.Context, msg ...string) {
	r := ResultBody{
		TraceId: c.GetHeader("trace_id"),
		Ret:     InternalServerError,
		Msg:     util.First("InternalServerError", msg),
	}
	r.result(c)
}

func ErrParams(c *gin.Context, msg ...string) {
	r := ResultBody{
		TraceId: c.GetHeader("trace_id"),
		Ret:     InvalidParameter,
		Msg:     util.First("InvalidParameter", msg),
	}
	r.result(c)
}

func ErrParamsWithErr(c *gin.Context, data ...any) {
	r := ResultBody{
		TraceId: c.GetHeader("trace_id"),
		Ret:     InvalidParameter,
		Msg:     "InvalidParameter",
		Data:    util.First(nil, data),
	}
	r.result(c)
}

func ErrParamsMissing(c *gin.Context, msg ...string) {
	r := ResultBody{
		TraceId: c.GetHeader("trace_id"),
		Ret:     MissingParameter,
		Msg:     util.First("MissingParameter", msg),
	}
	r.result(c)
}

func ErrToken(c *gin.Context, msg ...string) {
	r := ResultBody{
		TraceId: c.GetHeader("trace_id"),
		Ret:     InvalidToken,
		Msg:     util.First("InvalidToken", msg),
	}
	r.result(c)
}

func ErrBadRequest(c *gin.Context, msg ...string) {
	r := ResultBody{
		TraceId: c.GetHeader("trace_id"),
		Ret:     BadRequest,
		Msg:     util.First("BadRequest", msg),
	}
	r.result(c)
}

func ErrAuth(c *gin.Context, msg ...string) {
	r := ResultBody{
		TraceId: c.GetHeader("trace_id"),
		Ret:     Unauthorized,
		Msg:     util.First("Unauthorized", msg),
	}
	r.result(c)
}

func ErrNotFound(c *gin.Context, msg ...string) {
	r := ResultBody{
		TraceId: c.GetHeader("trace_id"),
		Ret:     NotFound,
		Msg:     util.First("NotFound", msg),
	}
	r.result(c)
}
