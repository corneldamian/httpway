package httpway

import (
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"fmt"
)

//get the context associated with the request
func GetContext(r *http.Request) *Context {
	crc, ok := r.Body.(contextReadCloser)

	if !ok {
		return nil
	}

	return crc.ctx()
}

//this is the context that is created for each request
type Context struct {
	data    map[string]interface{}
	logger  Logger
	session Session
	statusCode int
	transferedBytes uint64

	handlers          *[]httprouter.Handle
	runNextHandlerIdx int
}

//execute the next middleware
func (c *Context) Next(w http.ResponseWriter, r *http.Request, pr httprouter.Params) {
	c.runNextHandlerIdx--

	if c.runNextHandlerIdx < 0 {
		panic("No next middleware, don't call it in final handler")
	}

	(*c.handlers)[c.runNextHandlerIdx](w, r, pr)
}

//set a key on context
func (c *Context) Set(key string, value interface{}) {
	c.data[key] = value
}

//get a key from context and if was set
func (c *Context) GetOk(key string) (value interface{}, found bool) {
	value, found = c.data[key]
	return
}

//get a a key from the context
func (c *Context) Get(key string) interface{} {
	return c.data[key]
}

//check if a key was set on the context
func (c *Context) Has(key string) bool {
	_, has := c.data[key]

	return has
}

func (c *Context) StatusCode() int {
	return c.statusCode
}

func (c *Context) TransferedBytes() uint64 {
	return c.transferedBytes
}

//returns the logger associated with the request
func (c *Context) Log() Logger {
	if c.logger == nil {
		panic("No logger set")
	}

	return c.logger
}

//check if log is enabled
func (c *Context) HasLog() bool {
	if c.logger == nil {
		return false
	}

	return true
}

//returns the session associated with the request
func (c *Context) Session() Session {
	if c.session == nil {
		panic("No session set")
	}

	return c.session
}

//check if session is enabled
func (c *Context) HasSession() bool {
	if c.session == nil {
		return false
	}

	return true
}

func createContext(router *Router, w http.ResponseWriter, r *http.Request, handlers *[]httprouter.Handle, handlersLen *int) http.ResponseWriter {
	crc := &contextReadClose{
		ReadCloser: r.Body,
		ctxObj: &Context{
			data:              make(map[string]interface{}),
			handlers:          handlers,
			runNextHandlerIdx: *handlersLen,
		},
	}

	if router.SessionManager != nil {
		crc.ctxObj.session = router.SessionManager.Get(w, r)
	}

	if router.Logger != nil {
		crc.ctxObj.logger = router.Logger
	}

	r.Body = crc

	return &internalResponseWriter{w, crc.ctxObj}
}

type contextReadCloser interface {
	io.ReadCloser
	ctx() *Context
}

type contextReadClose struct {
	io.ReadCloser
	ctxObj *Context
}

func (crc *contextReadClose) ctx() *Context {
	return crc.ctxObj
}

type internalResponseWriter struct {
	rw http.ResponseWriter
	ctx *Context
}

func (irw *internalResponseWriter) Header() http.Header {
	return irw.rw.Header()
}

func (irw *internalResponseWriter)  Write(b []byte) (n int, err error) {
	if irw.ctx.statusCode == 0 {
		irw.ctx.statusCode = 200
	}
	n, err = irw.rw.Write(b)

	irw.ctx.transferedBytes += uint64(n)
	return
}

func (irw *internalResponseWriter) WriteHeader(status int) {
	if irw.ctx.statusCode == 0 {
		irw.ctx.statusCode = status
	}
	irw.rw.WriteHeader(status)
}