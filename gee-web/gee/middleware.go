package gee

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

type BaseMiddleWare struct{}

func (m *BaseMiddleWare) processRequest(c *Context) error {
	return nil
}

func (m *BaseMiddleWare) processView(c *Context) error {
	return nil
}

func (m *BaseMiddleWare) processResponse(c *Context) {}

func (m *BaseMiddleWare) processException(c *Context) error {return nil}

type IMiddleWare interface {
	processRequest(c *Context) error
	processView(c *Context) error
	processException(c *Context) error
	processResponse(c *Context)
}

type Logger struct {
	BaseMiddleWare
}

func (l *Logger) processRequest(c *Context) error {
	log.Printf("[%s] - path: %s", c.Method, c.Req.RequestURI)
	c.meta["startTime"] = time.Now()
	return nil
}

func (l *Logger) processResponse(c *Context) {
	log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(c.meta["startTime"].(time.Time)))
}

func (l *Logger) New() IMiddleWare {
	return &Logger{}
}

func NewLogger() IMiddleWare {
	return &Logger{}
}

type V2Logger struct {
	BaseMiddleWare
	prefix string
}

func (l *V2Logger) processRequest(c *Context) error {
	if !strings.HasPrefix(c.Req.RequestURI, l.prefix) {
		return errors.New(fmt.Sprintf("path： %s does not start with v2", c.Req.RequestURI))
	}
	log.Printf("path: %s is valid", c.Req.RequestURI)
	return nil
}

func NewV2Logger() *V2Logger {
	return &V2Logger{prefix: "/v2"}
}

type Recover struct {
	BaseMiddleWare
}

func (r *Recover) processException(c *Context) error{
	err, ok := c.meta["error"]
	if ok {
		err := errors.New(err.(string))
		c.Fail(500, err.Error())
		return err
	}
	return nil
}

func NewRecover() *Recover{
	return &Recover{}
}
