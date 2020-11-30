package router

import (
	"github.com/jonathonadams/bolt/bolt"
)	

// TODO -> Need to incorporate the Router root into the path to validate
func addHandler(
	router *Router,
	method string,
	path string,
	middleware []bolt.Middleware,
) bool {

	if len(middleware) == 0 {
		panic("No middleware supplied for " + method + ": /" + path)
	}

	segments := createPathSegments(path)

	node := router.trie.getRouteNode(segments)

	if node == nil {
		// register directly on the node
		node = createRouterNode()
		node.methods[method] = middleware

		_, err := router.trie.setRouterNode(segments, node)

		if err != nil {
			// TODO -> Log error/panic
			return false
		}
		return true
	}

	// TODO check if overriding?
	node.methods[method] = middleware
	return true

}
