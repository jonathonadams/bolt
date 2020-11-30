package router

import (
	"net/http"

	"github.com/jonathonadams/bolt/bolt"
)

type Router struct {
	root string
	trie *RouterNode
}

type RouterNode struct {
	middleware *[]bolt.Middleware
	path       string
	methods    MethodMiddleware
	children   map[string]*RouterNode
	tokenName  string
	tokenNode  *RouterNode
}

type MethodMiddleware = map[string][]bolt.Middleware

func NewRouter() *Router {
	router := Router{
		root: "/",
		trie: createRouterNode(),
	}

	return &router

}

func (r *Router) Routes() bolt.Middleware {
	return func(ctx *bolt.Ctx, next bolt.Next) {

		segments, err := splitIncomingRoute(ctx)
		if err != nil {
			// TODO
			return
		}

		middleware := r.trie.getRequestHandlers(ctx.Method, segments, &ctx.PathParams)

		if middleware != nil {
			noOpNext := func() {}
			composedMiddleware := bolt.ComposeMiddleware(middleware)
			composedMiddleware(ctx, noOpNext)

		} else {
			ctx.Status = 404
		}

		// Once the router is finished call next to run downstream middleware
		next()
	}
}

func (r *Router) Mount(path string, subRouter *Router) error {
	segments := createPathSegments(r.root + "/" + path)

	subRouter.root = joinURLSegments(segments)

	node := r.trie.getRouteNode(segments)

	if node != nil {
		node = subRouter.trie
	} else {
		// register
		r.trie.setRouterNode(segments, subRouter.trie)
	}

	return nil
}

func (r *Router) Use(middleWare ...bolt.Middleware) {
	*r.trie.middleware = append(*r.trie.middleware, middleWare...)
}

func (r *Router) GET(path string, middleware ...bolt.Middleware) {
	addHandler(r, http.MethodGet, path, middleware)
}

func (r *Router) PUT(path string, middleware ...bolt.Middleware) {
	addHandler(r, http.MethodPut, path, middleware)
}

func (r *Router) POST(path string, middleware ...bolt.Middleware) {
	addHandler(r, http.MethodPost, path, middleware)
}

func (r *Router) DELETE(path string, middleware ...bolt.Middleware) {
	addHandler(r, http.MethodDelete, path, middleware)
}
