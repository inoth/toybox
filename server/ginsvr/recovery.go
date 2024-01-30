package ginsvr

import (
	"bytes"
	"fmt"
	"os"
	"runtime"

	"github.com/inoth/toybox/component/logger"
	"github.com/inoth/toybox/server/ginsvr/res"

	"github.com/gin-gonic/gin"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

func stack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := os.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	if lastSlash := bytes.LastIndex(name, slash); lastSlash >= 0 {
		name = name[lastSlash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.ReplaceAll(name, centerDot, dot)
	return name
}

func source(lines [][]byte, n int) []byte {
	n--
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log := logger.New("Recovery", c.GetHeader("TraceId"))
				stack := stack(3)
				switch err := err.(type) {
				case error:
					log.Error(fmt.Sprintf("[Recovery:error] panic recovered:\n%s\n%s", err.Error(), string(stack)))
					res.Failed(c, err.Error())
				case string:
					log.Error(fmt.Sprintf("[Recovery:string] panic recovered:\n%s\n%s", err, string(stack)))
					res.Failed(c, err)
				default:
					log.Error(fmt.Sprintf("[Recovery:invalid] panic recovered:\n%s", string(stack)))
					res.Failed(c, string(stack))
				}
				c.Abort()
			}
		}()
		c.Next()
	}
}
