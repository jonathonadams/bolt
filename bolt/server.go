package bolt

import (
	"net/http"
)

type Next = func()
type Middleware = func(ctx *Ctx, next Next)

type Bolt struct {
	middleware []Middleware
	run        *Next
	handler    Handler
}

type Handler struct {
	stack Middleware
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

func NewBolt() *Bolt {
	app := Bolt{
		middleware: make([]Middleware, 0),
		handler:    Handler{},
	}
	return &app
}

func (z *Bolt) Use(m Middleware) {
	z.middleware = append(z.middleware, m)
}

func (b *Bolt) HttpHandler() Handler {
	middlewareStack := ComposeMiddleware(&b.middleware)
	b.handler.stack = middlewareStack

	return b.handler
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	ctx := createContext(w, r)

	defer func() {
		if r := recover(); r != nil {
			http.Error(ctx.W, "Something went wrong. Error:", http.StatusInternalServerError)
		}

		h.respond(w, ctx)
	}()

	noOpNext := func() {}
	h.stack(ctx, noOpNext)
}

func (h *Handler) respond(w http.ResponseWriter, ctx *Ctx) (int, error) {

	for k, v := range ctx.Headers {
		w.Header().Set(k, v)
	}

	w.WriteHeader(ctx.Status)

	if ctx.response != nil {
		return w.Write(*ctx.response)
	}

	return 0, nil

}
