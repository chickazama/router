package router

import (
	"net/http"
	"regexp"
)

type Router struct {
	Handlers        map[*regexp.Regexp]http.Handler
	NotFoundHandler http.Handler
}

func New() *Router {
	return &Router{Handlers: make(map[*regexp.Regexp]http.Handler)}
}

func (rt *Router) AddHandler(exp *regexp.Regexp, h http.Handler) {
	rt.Handlers[exp] = h
}

func (rt *Router) AddNotFoundHandler(h http.Handler) {
	rt.NotFoundHandler = h
}

func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for exp, h := range rt.Handlers {
		if exp.MatchString(r.URL.Path) {
			h.ServeHTTP(w, r)
			return
		}
	}
	if rt.NotFoundHandler != nil {
		rt.NotFoundHandler.ServeHTTP(w, r)
		return
	}
	http.Error(w, "404 - Not Found", http.StatusNotFound)
}

func (rt *Router) ServeStaticFiles() {
	fsRoot := http.Dir("./static/")
	fs := http.FileServer(fsRoot)
	rt.AddHandler(regexp.MustCompile(`^/static/`), http.StripPrefix("/static/", fs))
}
