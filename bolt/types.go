package bolt

import "net/http"

type Next = func()
type Middleware = func(ctx *Ctx, next Next)

type Bolt struct {
	middleware []Middleware
	run        *Next
}

type Ctx struct {
	Status     int
	Path       string
	Method     string
	Body       interface{}
	JSON       *[]byte
	PathParams map[string]string
	response   *[]byte
	Headers    map[string]string
	R          *http.Request
	W          http.ResponseWriter
}
