package gee

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
	"sync"
)

// HandlerFunc defines the request handler used by gee
type HandlerFunc func(*Context)

// Engine implement the interface of ServeHTTP
type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup
	pool sync.Pool
	htmlTemplates *template.Template // for html render
	funcMap       template.FuncMap
}

func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}


type RouterGroup struct {
	prefix      string
	middlewares []IMiddleWare
	parent      *RouterGroup
	engine      *Engine
}

// New is the constructor of gee.Engine
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.pool = sync.Pool{New: func() interface{} {
		return &Context{}
	}}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) addRoute(method, pattern string, handler HandlerFunc) {
	pattern = group.prefix + pattern
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

func (group *RouterGroup) Use(middlewares ...IMiddleWare) {
	group.middlewares = append(group.middlewares, middlewares...)
}

func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

// create static handler
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		// Check if file exists and/or if we have permission to access it
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

// serve static files
func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	// Register GET handlers
	group.GET(urlPattern, handler)
}


func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	engine.router.addRoute(method, pattern, handler)
}

// GET defines the method to add GE及T request
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

// Run defines the method to start a http server
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []IMiddleWare
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}

	//c := newContext(w, req, engine)
	c := engine.pool.Get().(*Context)
	c.set(w, req, engine)
	defer func() {
		if err := recover(); err != nil {
			c.meta["error"] = err
			processException(middlewares, c)
		}
	}()
	err := processRequest(middlewares, c)
	if err != nil {
		c.Fail(500, err.Error())
		return
	}
	n, params := engine.router.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		err = processView(middlewares, c)
		if err != nil {
			c.Fail(500, err.Error())
			return
		}
		engine.router.handlers[key](c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
	processResponce(middlewares, c)
	c.reset()
	engine.pool.Put(c)
}

func processRequest(middlewares []IMiddleWare, c *Context) error{
	for _, middleware := range middlewares {
		err := middleware.processRequest(c)
		if err != nil {
			return err
		}
	}
	return nil
}

func processResponce(middlewares []IMiddleWare, c *Context) {
	for i := len(middlewares) - 1; i >= 0; i-- {
		middlewares[i].processResponse(c)
	}
}

func processView(middlewares []IMiddleWare, c *Context) error {
	for _, middleware := range middlewares {
		err := middleware.processView(c)
		if err != nil {
			return err
		}
	}
	return nil
}


func processException(middlewares []IMiddleWare, c *Context) {
	for i := len(middlewares) - 1; i >= 0; i-- {
		err := middlewares[i].processException(c)
		if err != nil {
			return
		}
	}
}