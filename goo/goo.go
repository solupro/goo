package goo

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
)

type HandlerFunc func(*Context)

type Engine struct {
	*RouterGroup
	groups []*RouterGroup

	// template
	htmlTemplates *template.Template
	funcMap       template.FuncMap
}

type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc
	engine      *Engine
	router      *router
}

func (g *RouterGroup) Group(prefix string) *RouterGroup {
	newGroup := &RouterGroup{
		prefix: g.prefix + prefix,
		engine: g.engine,
		router: g.router,
	}
	g.engine.groups = append(g.engine.groups, newGroup)
	return newGroup
}

func (g *RouterGroup) addRouter(method, path string, handler HandlerFunc) {
	pattern := g.prefix + path
	log.Printf("Route %4s - %s", method, pattern)
	g.router.addRoute(method, pattern, handler)
}

func (g *RouterGroup) GET(path string, handler HandlerFunc) {
	g.addRouter("GET", path, handler)
}

func (g *RouterGroup) POST(path string, handler HandlerFunc) {
	g.addRouter("POST", path, handler)
}

func (g *RouterGroup) Use(middlewares ...HandlerFunc) {
	g.middlewares = append(g.middlewares, middlewares...)
}

func (g *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(g.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))

	return func(c *Context) {
		file := c.Param("filepath")
		if _, err := fs.Open(file); nil != err {
			c.Status(http.StatusNotFound)
			return
		}

		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

func (g *RouterGroup) Static(relativePath, root string) {
	handler := g.createStaticHandler(relativePath, http.Dir(root))
	pattern := path.Join(relativePath, "/*filepath")
	g.GET(pattern, handler)
}

func New() *Engine {
	engine := &Engine{}
	routerGroup := &RouterGroup{
		engine: engine,
		router: newRouter(),
	}
	engine.RouterGroup = routerGroup
	engine.groups = append(engine.groups, routerGroup)
	return engine
}

func Default() *Engine {
	engine := New()
	engine.Use(ApiCostTime(), Recovery())

	return engine
}

func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}

func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	c.engine = engine
	engine.router.handle(c)
}

func (engine *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, engine)
}
