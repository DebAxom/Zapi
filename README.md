# Zapi
Zapi is a lightweight backend web framework built on top of Go’s native net/http package. I created it as a personal passion project to have a simple, minimal, and expressive Go framework for my future work. While I designed the overall structure and features, I collaborated with ChatGPT to fine-tune certain parts especially the router, since I prefer to avoid the pain of working directly with regular expressions. Thanks to this, Zapi’s router supports clean, expressive, and fully parameterized routes right out of the box.

## Docs

### Initializing a new App
```go
package main

import (
    "fmt"
    "os"
	"path"
    "github.com/DebAxom/Zapi"
)

func main() {
    App := zapi.New() // Initializing a new app

    App.Get("/", func(req *zapi.Req, res zapi.Res) {
		res.WriteHeader(200)
		res.Write("Hello World !")
	})

    fmt.Println("Server running on http://localhost:3000")
    App.run(":3000") // Running the server on port 3000
}

```

### Methods
```go
App.Get(route, handlerFunc) // Handle GET Request
App.Post(route, handlerFunc) // Handle POST Request
App.Put(route, handlerFunc) // Handle PUT Request
App.Delete(route, handlerFunc) // Handle DELETE Request
```

### Routing
```go
    App := zapi.New()
    App.Get("/", func(req *zapi.Req, res zapi.Res){ ... })
    App.Get("/blog", func(req *zapi.Req, res zapi.Res){ ... })
    
    // Parameterized Routes
    App.Get("/blog/[id]", func(req *zapi.Req, res zapi.Res){ 
        id := req.Params["id"]
        ....
    })

    App.Get("/@[username]/p-[id]", func(req *zapi.Req, res zapi.Res){ 
        username := req.Params["username"]
        id := req.Params["id"]
        ....
    })

    // Handling 404 ; This route should be at the end of all routes
    App.Get("*", func(req *zapi.Req, res zapi.Res){
        res.WriteHeader(404)
        res.Write("404 Not Found !")
    })

```

### Req struct
Built on top of `http.Request`. Contains additional map `Params` and an additional function `BindJSON`.

### Res struct
It contains all the methods of `http.ResponseWriter` along with a few additional features.

### CORS
```go
App := zapi.New()
APP.CORS.AllowedMethods = []string{"GET", "POST"}
APP.CORS.AllowedOrigins = []string{"http://localhost:8080"}
APP.CORS.AllowCredentials = true
```

### Handling Cookies
```go
// Get cookie using the GetCookie method on Req struct
c, err := req.GetCookie(cookieName)

// Deleting a cookie
req.DeleteCookie(cookieName)

// Set cookie using the SetCookie method on Req struct with default options
req.SetCookie(cookieName, cookieValue, nil)

// Set cookie using the SetCookie method on Req struct with your own options
req.SetCookie(cookieName, cookieValue, &zapi.CookieOptions{...})

// All options available
type CookieOptions struct {
    Path     string
	Domain   string
	Expires  time.Time
	MaxAge   int
	Secure   bool
	HttpOnly bool
	SameSite http.SameSite
}
```

### Serving Static Files
``` go
App.Public(route, dirPath)

/*Example :
    App.Public("/assets", path.Join(os.Getwd(), "assets"))
*/
```

### Sending File as response
```go
App.Get("/", func(req *zapi.Req, res zapi.Res) {
	res.WriteHeader(200)
	res.SendFile(filepath)
})
```

### Sending JSON as response
```go
// Sending a map as response

App.Get("/", func(req *zapi.Req, res zapi.Res) {
	res.WriteHeader(200)
	res.JSON(map[string]string{"msg" : "Hello World"})
})

// Sending a struct as respone

type User struct {
	Name    string   `json:"name"`
	Age     int      `json:"age"`
	Hobbies []string `json:"hobbies"`
}

App.Get("/user", func(req *zapi.Req, res zapi.Res) {
	res.WriteHeader(200)
	res.JSON(User{
		Name:    "DebAxom",
		Age:     19,
		Hobbies: []string{"Coding", "Watching Youtube"},
	})
})
```

### Getting JSON data from Body
```go

type Auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

App.Get("/auth", func(req *zapi.Req, res zapi.Res) {
	
    var data Auth
    err := req.BindJSON(&data)

    if err!=nil{
        // Handle Error
    }

    username := Auth.Username
    pwd := Auth.Password
	
    // Do whatever you want with this data
})
```

### Redirecting user to a different route
```go
App.Get("/index", func(req *zapi.Req, res zapi.Res) {
	res.Redirect("/")
})
```