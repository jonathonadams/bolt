package router

import (
	"regexp"
	"strings"

	"github.com/jonathonadams/bolt/bolt"
)

func createPathSegments(path string) *[]string {
	var validSegments []string
	segments := strings.Split(normalizeURL(path), "/")

	for _, slug := range segments {
		if slug != "" {
			if ok := isSlugValid(slug); ok == false {
				panic("slug " + slug + " is not a valid url")
			}
			validSegments = append(validSegments, slug)
		}
	}
	return &validSegments
}

func joinURLSegments(segments *[]string) string {
	return strings.Join((*segments), "/")
}

// the token symbol is ':id'
var validToken = regexp.MustCompile("^:[a-z]{1}[a-zA-Z]*$")

func isSlugAToken(slug string) (bool, string) {
	match := validToken.MatchString(slug)
	if match == true {
		return true, slug[1:]
	}
	return false, ""
}

var validSlug = regexp.MustCompile("^[a-zA-Z0-9_-]+$")

// first check if its a token, if not check its valid
func isSlugValid(slug string) bool {
	token, _ := isSlugAToken(slug)
	if token == true {
		return true
	}
	return validSlug.MatchString(slug)
}

var multiSlashes = regexp.MustCompile("//+")
var trailingSlashes = regexp.MustCompile("/$")
var leadingSlashes = regexp.MustCompile("^/")
var single = []byte("/")
var empty = []byte("")

func normalizeURL(path string) string {
	normal := multiSlashes.ReplaceAll([]byte(path), single)
	normal = trailingSlashes.ReplaceAll(normal, empty)
	normal = leadingSlashes.ReplaceAll(normal, empty)
	return string(normal)
}

func createRouterNode() *RouterNode {
	m := make([]bolt.Middleware, 0)
	return &RouterNode{
		middleware: &m,
		children:   make(map[string]*RouterNode),
		methods:    make(map[string][]bolt.Middleware),
	}
}
