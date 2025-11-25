package httprouter

import (
	"net/http"
	"strings"
)

// =============== INTERFACE ===============

type HttpRouter interface {
	http.Handler

	Use(mw func(http.Handler) http.Handler)

	Handle(pattern string, h http.Handler, mws ...func(http.Handler) http.Handler)
	HandleFunc(pattern string, h http.HandlerFunc, mws ...func(http.Handler) http.Handler)

	GET(path string, h http.HandlerFunc, mws ...func(http.Handler) http.Handler)
	POST(path string, h http.HandlerFunc, mws ...func(http.Handler) http.Handler)
	PUT(path string, h http.HandlerFunc, mws ...func(http.Handler) http.Handler)
	PATCH(path string, h http.HandlerFunc, mws ...func(http.Handler) http.Handler)
	DELETE(path string, h http.HandlerFunc, mws ...func(http.Handler) http.Handler)

	Group(prefix string, fn func(g HttpRouter))
	Static(prefix string, dir string, mws ...func(http.Handler) http.Handler)
}

// =============== IMPLEMENTASI ===============

type Router struct {
	mux         *http.ServeMux
	middlewares []func(http.Handler) http.Handler
}

type Group struct {
	router      *Router
	prefix      string
	middlewares []func(http.Handler) http.Handler
}

func New() *Router {
	return &Router{
		mux:         http.NewServeMux(),
		middlewares: nil,
	}
}

// Middleware global (untuk semua route di router ini)
func (r *Router) Use(mw func(http.Handler) http.Handler) {
	r.middlewares = append(r.middlewares, mw)
}

// chain: apply mwN(mwN-1(...(h))) urutan benar
func chain(h http.Handler, mws []func(http.Handler) http.Handler) http.Handler {
	final := h
	for i := len(mws) - 1; i >= 0; i-- {
		final = mws[i](final)
	}
	return final
}

// Handle: pattern full, contoh: "GET /users/{id}"
func (r *Router) Handle(pattern string, h http.Handler, mws ...func(http.Handler) http.Handler) {
	all := append(append([]func(http.Handler) http.Handler{}, r.middlewares...), mws...)
	final := chain(h, all)

	r.mux.Handle(pattern, final)
}

// HandleFunc: helper untuk http.HandlerFunc
func (r *Router) HandleFunc(pattern string, h http.HandlerFunc, mws ...func(http.Handler) http.Handler) {
	r.Handle(pattern, h, mws...)
}

// Helper: method + path (Go 1.22+ pattern)
func (r *Router) GET(path string, h http.HandlerFunc, mws ...func(http.Handler) http.Handler) {
	r.Handle("GET "+clean(path), h, mws...)
}

func (r *Router) POST(path string, h http.HandlerFunc, mws ...func(http.Handler) http.Handler) {
	r.Handle("POST "+clean(path), h, mws...)
}

func (r *Router) PUT(path string, h http.HandlerFunc, mws ...func(http.Handler) http.Handler) {
	r.Handle("PUT "+clean(path), h, mws...)
}

func (r *Router) PATCH(path string, h http.HandlerFunc, mws ...func(http.Handler) http.Handler) {
	r.Handle("PATCH "+clean(path), h, mws...)
}

func (r *Router) DELETE(path string, h http.HandlerFunc, mws ...func(http.Handler) http.Handler) {
	r.Handle("DELETE "+clean(path), h, mws...)
}

// Group: prefix + middleware khusus group
func (r *Router) Group(prefix string, fn func(g HttpRouter)) {
	g := &Group{
		router:      r,
		prefix:      clean(prefix),
		middlewares: nil,
	}
	fn(g)
}

func (r *Router) Static(prefix string, dir string, mws ...func(http.Handler) http.Handler) {
	// normalisasi prefix, harus berakhir dengan "/"
	prefix = clean(prefix)
	if !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}

	// FileServer standar
	fs := http.FileServer(http.Dir(dir))
	base := http.StripPrefix(prefix, fs)

	// Guard method: hanya GET & HEAD
	staticHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet && req.Method != http.MethodHead {
			w.Header().Set("Allow", "GET, HEAD")
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		base.ServeHTTP(w, req)
	})

	// global mw + mw khusus static ini
	all := append(append([]func(http.Handler) http.Handler{}, r.middlewares...), mws...)
	final := chain(staticHandler, all)

	// Pattern dengan wildcard untuk match semua file: "/static/{path...}"
	pattern := prefix + "{path...}"
	r.mux.Handle(pattern, final)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

// =============== GROUP ===============

func (g *Group) Use(mw func(http.Handler) http.Handler) {
	g.middlewares = append(g.middlewares, mw)
}

func (g *Group) Handle(pattern string, h http.Handler, mws ...func(http.Handler) http.Handler) {
	// Parse pattern untuk extract method dan path
	// Pattern bisa: "GET /users" atau "/users" (tanpa method)
	var fullPattern string
	if idx := strings.Index(pattern, " "); idx > 0 {
		// Ada method: "GET /users"
		method := pattern[:idx]
		path := pattern[idx+1:]
		fullPath := join(g.prefix, path)
		fullPattern = method + " " + fullPath
	} else {
		// Tanpa method: "/users" atau "/static/{path...}"
		fullPattern = join(g.prefix, pattern)
	}

	// global router → group → route
	all := append([]func(http.Handler) http.Handler{}, g.router.middlewares...)
	all = append(all, g.middlewares...)
	all = append(all, mws...)

	final := chain(h, all)
	g.router.mux.Handle(fullPattern, final)
}

func (g *Group) HandleFunc(pattern string, h http.HandlerFunc, mws ...func(http.Handler) http.Handler) {
	g.Handle(pattern, h, mws...)
}

func (g *Group) GET(path string, h http.HandlerFunc, mws ...func(http.Handler) http.Handler) {
	g.Handle("GET "+clean(path), h, mws...)
}

func (g *Group) POST(path string, h http.HandlerFunc, mws ...func(http.Handler) http.Handler) {
	g.Handle("POST "+clean(path), h, mws...)
}

func (g *Group) PUT(path string, h http.HandlerFunc, mws ...func(http.Handler) http.Handler) {
	g.Handle("PUT "+clean(path), h, mws...)
}

func (g *Group) PATCH(path string, h http.HandlerFunc, mws ...func(http.Handler) http.Handler) {
	g.Handle("PATCH "+clean(path), h, mws...)
}

func (g *Group) DELETE(path string, h http.HandlerFunc, mws ...func(http.Handler) http.Handler) {
	g.Handle("DELETE "+clean(path), h, mws...)
}

func (g *Group) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	g.router.ServeHTTP(w, req)
}

func (g *Group) Group(prefix string, fn func(g HttpRouter)) {
	newGroup := &Group{
		router: g.router,
		prefix: join(g.prefix, prefix),

		// Mewarisi middleware parent group
		middlewares: append([]func(http.Handler) http.Handler{}, g.middlewares...),
	}

	fn(newGroup)
}

func (g *Group) Static(prefix string, dir string, mws ...func(http.Handler) http.Handler) {
	// prefix group + prefix static
	fullPrefix := join(g.prefix, prefix)
	if !strings.HasSuffix(fullPrefix, "/") {
		fullPrefix += "/"
	}

	fs := http.FileServer(http.Dir(dir))
	base := http.StripPrefix(fullPrefix, fs)

	staticHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			w.Header().Set("Allow", "GET, HEAD")
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		base.ServeHTTP(w, r)
	})

	// global router → group mw → mw tambahan
	all := append([]func(http.Handler) http.Handler{}, g.router.middlewares...)
	all = append(all, g.middlewares...)
	all = append(all, mws...)

	final := chain(staticHandler, all)

	// Pattern dengan wildcard untuk match semua file
	pattern := fullPrefix + "{path...}"
	g.router.mux.Handle(pattern, final)
}

// =============== UTIL PATH ===============

func clean(p string) string {
	p = strings.TrimSpace(p)
	if p == "" || p == "/" {
		return "/"
	}
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	if len(p) > 1 && strings.HasSuffix(p, "/") {
		p = strings.TrimSuffix(p, "/")
	}
	return p
}

func join(prefix, path string) string {
	prefix = clean(prefix)
	path = clean(path)

	if prefix == "/" {
		return path
	}
	if path == "/" {
		return prefix
	}
	return prefix + path
}
