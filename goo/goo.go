package goo

import (
	"log"
	"net/http"
	"strings"
)

type HandlerFunc func(*Context)

type Engine struct {
	*RouterGroup
	groups []*RouterGroup
}

type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc
	engine      *Engine
	router      *router
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

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	engine.router.handle(c)
}

func (engine *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, engine)
}
