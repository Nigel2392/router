# Router

Easy to use router for your Golang HTTP server.

## Installation
```bash
go get github.com/Nigel2392/router
```

## Usage
```go
	// Create router
	r := router.NewRouter()

	// Handle static files
	// Matches anything after /static/
	r.Handle("GET", "/static/<<any>>", http.StripPrefix("/static/", http.FileServer(staticFiles)))

	// /home/joe
	var home = r.Get("/home/<<name:string>>", func(v router.Vars, w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %v", v)
	})

	// /home/Joe/page/123
	home.Get("/page/<<id:int>>", func(v router.Vars, w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %v", v)
	})

	// Middleware for a single route
	home.Use(func(next router.Handler) router.Handler {
		return router.HandleFuncWrapper{F: func(v router.Vars, w http.ResponseWriter, r *http.Request) {
			fmt.Println("Home middleware", r.URL.Path)
			next.ServeHTTP(v, w, r)
		}}
	})

	// Global middleware
	r.Use(middleware.Printer)

	r.HandleFunc("GET", "/<<any>>", indexFunc)
	r.HandleFunc("GET", "/", indexFunc)

	// Start server
	http.ListenAndServe(addr, r)

```