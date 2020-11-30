package bolt

import (
	"net/http"
)

func NewBolt() *Bolt {
	app := Bolt{
		middleware: make([]Middleware, 0),
	}
	return &app
}

func (z *Bolt) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	ctx := createContext(w, r)

	defer func() {
		if r := recover(); r != nil {
			http.Error(ctx.W, "Something went wrong. Error:", http.StatusInternalServerError)
		}

		z.respondWith(w, ctx)
	}()

	middlewareStack := ComposeMiddleware(&z.middleware)

	noOpNext := func() {}
	middlewareStack(ctx, noOpNext)

}

func (z *Bolt) respondWith(w http.ResponseWriter, ctx *Ctx) (int, error) {

	for k, v := range ctx.Headers {
		w.Header().Set(k, v)
	}

	w.WriteHeader(ctx.Status)

	if ctx.response != nil {
		return w.Write(*ctx.response)
	}

	return 0, nil

}

func (z *Bolt) Use(m Middleware) {
	z.middleware = append(z.middleware, m)
}
