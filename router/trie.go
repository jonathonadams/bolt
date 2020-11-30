package router

import (
	"errors"
)

// TODO Document token replacement
func (m *RouterNode) setRouterNode(
	segments *[]string,
	ma *RouterNode) (*RouterNode, error) {

	// segments will ALWAYS have a length greater than 0
	// because when retrieving the node, it will always retrive the root
	// router and hence the `registerHandler` will never return nil

	slug := (*segments)[0]
	remaining := (*segments)[1:]

	if len(*segments) == 1 {
		token, name := isSlugAToken(slug)

		if token == true {
			if m.tokenName == "" {
				m.tokenName = name
				m.tokenNode = ma
				return ma, nil
			}
			// TODO Give full path error message
			return nil, errors.New("A path paramater has already been registered")
		}
		// TODO: Currently overriding the node if there?
		m.children[slug] = ma
		return ma, nil
	}

	// there are more that one segment left
	token, name := isSlugAToken(slug)
	if token == true {
		if m.tokenName == "" {
			tokenNode := createRouterNode()
			m.tokenName = name
			m.tokenNode = tokenNode
			return tokenNode.setRouterNode(&remaining, ma)
		}
		// TODO Give full path error message
		return nil, errors.New("A path paramater has already been registered")
	}

	node, ok := m.children[slug]
	if ok == false {
		node = createRouterNode()
		m.children[slug] = node
		return node.setRouterNode(&remaining, ma)
	}

	return node.setRouterNode(&remaining, ma)
}

func (m *RouterNode) getRouteNode(segments *[]string) *RouterNode {
	if len(*segments) == 0 {
		return m
	}

	path := (*segments)[0]
	remaining := (*segments)[1:]

	if node, ok := m.children[path]; ok {
		return node.getRouteNode(&remaining)
	}

	if m.tokenName != "" {
		return m.tokenNode.getRouteNode(&remaining)
	}

	return nil
}
