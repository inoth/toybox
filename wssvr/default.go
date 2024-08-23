package wssvr

import "fmt"

func defaultHandle() HandlerFunc {
	return func(c *Context) {
		fmt.Println(string(c.body))
		c.Abort()
	}
}
