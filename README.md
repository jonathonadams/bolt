Early development ALPHA stage

# Bolt 
Koa like middleware in Go

## Usage
```go
package main

import (
    "fmt"
	"log"
	"net/http"
    
    "github.com/jonathonadams/bolt/bolt"
)


func runFirst(ctx *bolt.Ctx, next bolt.Next) {
	fmt.Println("I run first")
	next()
	fmt.Println("I run fourth")
}

func runSecond(ctx *bolt.Ctx, next bolt.Next) {
	fmt.Println("I run second")
	next()
	fmt.Println("I run third")
}

func main() {

	port := "9000"

    app := bolt.NewBolt()
    
    app.Use(runFirst)
    app.Use(runSecond)

	fmt.Printf("Server is running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, app))
}
```

Output:
```bash
I run first
I run second
I run third
I run fourth
```

## Router

```go
package main

import (
    "fmt"
	"log"
	"net/http"
    
    "github.com/jonathonadams/bolt/bolt"
    "github.com/jonathonadams/bolt/router"
)

func getRoot(ctx *bolt.Ctx, next bolt.Next) {
	ctx.Body = "{ path: '/'}"
}

func pointlessMiddleware(ctx *bolt.Ctx, next bolt.Next) {
    next()
}

func getTodo(ctx *bolt.Ctx, next bolt.Next) {
    todoId := ctx.PathParams["todoId"]
	ctx.Body = "{path: '/todos/" + todoId + "' }"
}

func secure(ctx *bolt.Ctx, next bolt.Next) {
	ctx.Body = "{path: 'secure'}"
}

func auth(ctx *bolt.Ctx, next bolt.Next) {
    panic("You shall not pass!!!")
}

func main() {

	port := "9000"

    app := bolt.NewBolt()

	app.Use(bolt.BodyParser())
    
    router := router.NewRouter()

    // Add specific middleware 
    router.GET("/", getRoot) // GET/PUT/POST/DELETE etc....
	router.GET("todos/:todoId", pointlessMiddleware, getTodo)

    secureRouter := router.NewRouter()
    
    secureRouter.Use(auth)
	secureRouter.POST("/try-me", secure)
	
    router.Mount("secure", secureRouter)

    app.Use(router.Routes())

	fmt.Printf("Server is running on port %s\n", port)
    log.Fatal(http.ListenAndServe(":"+port, app))
}
```

GET :: "/" -> "{ path: '/'}"  
GET :: "/todos/23" -> "{ path: '/todos/23'}"  
GET :: "/secure/try-me" -> 500 Internal server error. 



### TODO
- Document opinionated path correction
- Ctx API