package bolt

func ComposeMiddleware(middleWare *[]Middleware) Middleware {

	if len(*middleWare) != 0 {
		return func(ctx *Ctx, n Next) {
		}
	}

	return func(ctx *Ctx, n Next) {
		var next Next
		var run func(i int)

		// how many times next has been called
		called := -1

		run = func(i int) {

			// If next is called more than once in the same function then panic
			if i <= called {
				panic("next() can only be called once in any given middleware.")
			}

			called = i
			fn := (*middleWare)[i]

			if i == len(*middleWare)-1 {
				// TODO
				// If the middleware is last, the next is a no-op and the panic error will not be thrown
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
