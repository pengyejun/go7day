package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	// origin objects
	Writer http.ResponseWriter
	Req    *http.Request
	// request info
	Path   string
	Method string
	// http req params
	Params map[string]string
	// response info
	StatusCode int
	// middleware params
	meta map[string]interface{}
	engine *Engine
}

func (c *Context) reset() {
	c.Writer = nil
	c.Params = nil
	c.Req = nil
	c.Path = ""
	c.StatusCode = 0
	c.meta = nil
}

func (c *Context) set(w http.ResponseWriter, req *http.Request, engine *Engine) {
	c.Writer = w
	c.Req = req
	c.Path = req.URL.Path
	c.Method = req.Method
	c.meta = make(map[string]interface{})
	c.engine = engine
}

func newContext(w http.ResponseWriter, req *http.Request, engine *Engine) *Context {
	return &Context{
		Writer:    w,
		Req:       req,
		Path:      req.URL.Path,
		Method:    req.Method,
		meta: make(map[string]interface{}),
		engine:engine,
	}
}

func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) Fail(code int, msg string) {
	c.Status(code)
	c.Writer.Write([]byte(msg))
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

func (c *Context) HTML(code int, name string, data interface{}) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	if err := c.engine.htmlTemplates.ExecuteTemplate(c.Writer, name, data); err != nil {
		c.Fail(500, err.Error())
	}
}
