package bolt

// TODO -> Credit KoA Compose
func ComposeMiddleware(middleWare *[]Middleware) Middleware {

	return func(ctx *Ctx, n Next) {
		if len(*middleWare) != 0 {

			// how many times next has been called
			nextCalled := -1
			var next func()
			var run func(i int)

			run = func(i int) {
				if i <= nextCalled {
					// this happend if next is called more than once in the same function
					// throe error
					panic("next() can only be called once in any given middleware.")
				}
				nextCalled = i
				fn := (*middleWare)[i]

				if i == len(*middleWare)-1 {
					// TODO
					// If the middleware is last, the next is a no-op and the panic error will not be shown
					// Is this ok?
					next = n
				} else {
					next = func() {
						run(i + 1)
					}
				}

				fn(ctx, next)
			}

			run(0)
		}
	}
}
