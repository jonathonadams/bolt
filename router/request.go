package router

import (
	"errors"
	"strings"

	"github.com/jonathonadams/bolt/bolt"
)

/**
 * If a trailing slash, 301 redirect to the one that does not contain a forward slash.
 * if double slashes next to each other than 301 redirect
 * https://webmasters.stackexchange.com/questions/8354/what-does-the-double-slash-mean-in-urls/8381#8381
 *
 * TODO -> Review this, because the 301 might not be correct and trailing slash might not be corrected
 * Also, regext might be better than to build a new slice as 'string.Split()' return array anyway
 *
 * TODO: Should this be permenant redirect
 */
func splitIncomingRoute(ctx *bolt.Ctx) (*[]string, error) {

	panic("I PANICED")

	var validSegments []string

	// the path will always have a leading / even on an empty request
	// so strip the first value
	path := ctx.Path[1:]

	// if the remaining path is empty, then it is the root route
	if path == "" {
		return &validSegments, nil
	}

	segments := strings.Split(path, "/")

	for i, slug := range segments {

		// if the slug is empty and it is not the last iteration
		// then two `//` were next to each other. This corrects to a single "/"
		if slug == "" && i != len(segments)-1 {
			ctx.Status = 301
			ctx.Headers["Location"] = strings.ReplaceAll(ctx.Path, "//", "/")

			return nil, errors.New("301")
		}

		// If the path ends with a "/", send a permenant redirect to
		// the same path without the trailing "/"
		if slug == "" && i == len(segments)-1 {
			ctx.Status = 301
			ctx.Headers["Location"] = ctx.Path[:len(ctx.Path)-1]

			return nil, errors.New("301")
		}

		validSegments = append(validSegments, slug)
	}
	return &validSegments, nil
}

func (m *RouterNode) getRequestHandlers(
	method string,
	segments *[]string,
	tokens *map[string]string,
) *[]bolt.Middleware {

	requestMiddleware := make([]bolt.Middleware, len(*m.middleware))
	copy(requestMiddleware, *m.middleware)

	node := m.getRequestNode(segments, &requestMiddleware, tokens)

	if node != nil {
		methodHandlers, ok := node.methods[method]
		if ok {

			middlewareStack := append(requestMiddleware, methodHandlers...)
			return &middlewareStack
		}
	}
	return nil
}

func (m *RouterNode) getRequestNode(segments *[]string, middleware *[]bolt.Middleware, pathParams *map[string]string) *RouterNode {
	if len(*segments) == 0 {
		*middleware = append(*middleware, *m.middleware...)
		return m
	}

	path := (*segments)[0]
	remaining := (*segments)[1:]

	if node, ok := m.children[path]; ok {
		*middleware = append(*middleware, *node.middleware...)
		return node.getRequestNode(&remaining, middleware, pathParams)
	}

	if m.tokenName != "" {
		(*pathParams)[m.tokenName] = path
		*middleware = append(*middleware, *m.middleware...)
		return m.tokenNode.getRequestNode(&remaining, middleware, pathParams)
	}

	return nil
}
