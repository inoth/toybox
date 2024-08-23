package wssvr

import (
	"encoding/json"
	"sync"
	"time"
)

type HandlerFunc func(*Context)

type Context struct {
	body  []byte
	state bool

	Keys map[string]any
	m    sync.RWMutex
	ws   *WebsocketServer
}

func (c *Context) reset() {
	c.Keys = nil
	c.body = nil
	c.state = true
}

func (c *Context) send(msg []byte) {
	c.body = msg
}

func (c *Context) Body() []byte {
	return c.body
}

func (c *Context) Json(id string, data any) {
	buf, err := json.Marshal(&data)
	if err != nil {
		return
	}
	c.Render(id, buf)
}

func (c *Context) String(id, body string) {
	c.Render(id, []byte(body))
}

func (c *Context) Render(id string, msg []byte) {
	c.Abort()
	c.ws.output <- Message{
		ID:   id,
		Body: msg,
	}
}

func (c *Context) Abort() {
	c.state = false
}

func (c *Context) BindJson(obj any) error {
	if err := json.Unmarshal(c.body, obj); err != nil {
		return err
	}
	return nil
}

func (c *Context) Set(key string, val any) {
	c.m.Lock()
	defer c.m.Unlock()
	if c.Keys == nil {
		c.Keys = make(map[string]any)
	}
	c.Keys[key] = val
}

func (c *Context) Get(key string) (value any, exists bool) {
	c.m.RLock()
	defer c.m.RUnlock()
	value, exists = c.Keys[key]
	return
}

func (c *Context) GetString(key string) (s string) {
	if val, ok := c.Get(key); ok && val != nil {
		s, _ = val.(string)
	}
	return
}

func (c *Context) GetBool(key string) (b bool) {
	if val, ok := c.Get(key); ok && val != nil {
		b, _ = val.(bool)
	}
	return
}

func (c *Context) GetInt(key string) (i int) {
	if val, ok := c.Get(key); ok && val != nil {
		i, _ = val.(int)
	}
	return
}

func (c *Context) GetInt64(key string) (i64 int64) {
	if val, ok := c.Get(key); ok && val != nil {
		i64, _ = val.(int64)
	}
	return
}

func (c *Context) GetUint(key string) (ui uint) {
	if val, ok := c.Get(key); ok && val != nil {
		ui, _ = val.(uint)
	}
	return
}

func (c *Context) GetUint64(key string) (ui64 uint64) {
	if val, ok := c.Get(key); ok && val != nil {
		ui64, _ = val.(uint64)
	}
	return
}

func (c *Context) GetFloat64(key string) (f64 float64) {
	if val, ok := c.Get(key); ok && val != nil {
		f64, _ = val.(float64)
	}
	return
}

func (c *Context) GetTime(key string) (t time.Time) {
	if val, ok := c.Get(key); ok && val != nil {
		t, _ = val.(time.Time)
	}
	return
}

func (c *Context) GetDuration(key string) (d time.Duration) {
	if val, ok := c.Get(key); ok && val != nil {
		d, _ = val.(time.Duration)
	}
	return
}

func (c *Context) GetStringSlice(key string) (ss []string) {
	if val, ok := c.Get(key); ok && val != nil {
		ss, _ = val.([]string)
	}
	return
}

func (c *Context) GetStringMap(key string) (sm map[string]any) {
	if val, ok := c.Get(key); ok && val != nil {
		sm, _ = val.(map[string]any)
	}
	return
}

func (c *Context) GetStringMapString(key string) (sms map[string]string) {
	if val, ok := c.Get(key); ok && val != nil {
		sms, _ = val.(map[string]string)
	}
	return
}

func (c *Context) GetStringMapStringSlice(key string) (smss map[string][]string) {
	if val, ok := c.Get(key); ok && val != nil {
		smss, _ = val.(map[string][]string)
	}
	return
}
