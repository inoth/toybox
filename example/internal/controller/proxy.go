package controller

import (
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/inoth/toybox/config"
	"github.com/inoth/toybox/ginsvr"
)

type ProxyController struct {
	Host string `toml:"host"`

	reverseProxy *httputil.ReverseProxy
}

func NewProxyController(conf config.ConfigMate) *ProxyController {
	p := ProxyController{}
	err := conf.PrimitiveDecode(&p)
	if err != nil {
		panic(err)
	}
	url, err := url.Parse(p.Host)
	if err != nil {
		panic(err)
	}
	p.reverseProxy = httputil.NewSingleHostReverseProxy(url)
	return &p
}

func (p *ProxyController) Name() string { return "proxy" }

func (p *ProxyController) Prefix() string { return "" }

func (p *ProxyController) Middlewares() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		p.FlowCount,
		p.Reject,
	}
}

func (p *ProxyController) Routers() []ginsvr.Router { return nil }

func (p *ProxyController) FlowCount(c *gin.Context) {
	c.Request.Header.Add("SELF_PROXY", "SELF_PROXY")
	c.Next()
}

func (p *ProxyController) Reject(c *gin.Context) {
	p.reverseProxy.ServeHTTP(c.Writer, c.Request)
	c.Abort()
}
