package res

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Ok(c *gin.Context, msg string, data ...interface{}) {
	if len(data) > 0 {
		c.JSON(http.StatusOK, ok(msg, data[0]))
	} else {
		c.JSON(http.StatusOK, resultOK(msg))
	}
}

func Err(c *gin.Context, msg string) {
	c.JSON(FAILED, err(msg))
}
func Errf(c *gin.Context, msg string, args ...interface{}) {
	c.JSON(FAILED, err(fmt.Sprintf(msg, args...)))
}

func NotFound(c *gin.Context, msg string) {
	c.JSON(NOTFOUND, notFound(msg))
}
func NotFoundf(c *gin.Context, msg string, args ...interface{}) {
	c.JSON(NOTFOUND, notFound(fmt.Sprintf(msg, args...)))
}

func ParamErr(c *gin.Context, msg string) {
	c.JSON(PARAMETERERR, paramErr(msg))
}

func ParamErrf(c *gin.Context, msg string, args ...interface{}) {
	c.JSON(PARAMETERERR, paramErr(fmt.Sprintf(msg, args...)))
}

func Unauthrized(c *gin.Context, msg string) {
	c.JSON(UNAUTHORIZATION, unauthrized(msg))
}

func Unauthrizedf(c *gin.Context, msg string, args ...interface{}) {
	c.JSON(UNAUTHORIZATION, unauthrized(fmt.Sprintf(msg, args...)))
}
