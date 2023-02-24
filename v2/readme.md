# Golang HTTP Router

## Installation
```bash
go get github.com/Nigel2392/router/v2@latest
```

## Usage

### Registering routes
```go
func indexFunc(req *request.Request) {
	// Sessions are easily accessible from the request when utilizing
	// the right middlewares!
	req.Session.Set("Hello", "World")
    req.WriteString("Index!")
}

func variableFunc(req *request.Request) {
    // Get from the session
    fmt.Println(rq.Session.Get("Hello"))
    // Get from the URL parameters
	for k, v := range rq.URLParams {
		rq.WriteString(fmt.Sprintf("%s: %s ", k, v))
	}
}

func defaultJSONFunc(req *request.Request) {
    // Easily render JSON responses!
    req.JSON.Encode(map[string]interface{}{
        "Hello": "World",
    })
}

var r = router.NewRouter(nil)
r.Use(middleware.Printer)

// Register URLs
r.Get("/", indexFunc, "index")
r.Get("/user/<<name:string>>/<<id:int>>", variableFunc, "user")

// Register groups of URLs
var group = r.Group("/api", "api")
group.Get("/json", defaultJSONFunc, "json")
group.Post("/json2", defaultJSONFunc, "json2")
```

### Getting urls, formatting them
```go
// Find routess by name with the following syntax:
var index = r.URL(router.ALL, "index")
var variableRoute = r.URL(router.ALL, "variable")

// Format the route urls.
fmt.Println(index.Format())
fmt.Println(variableRoute.Format("John-Doe", 123))

// Extra parameters are ignored!
fmt.Println(variableRoute.Format("John-Doe", 123, 456, 789))

// Getting child urls
var json = r.URL(router.ALL, "api:json")
var json2 = r.URL(router.POST, "api:json2")
fmt.Println(json.Format())
fmt.Println(json2.Format())
```

### Rendering templates
Firstly, we need to define some variables in the `router/templates` package like so:
```go
// Configure default template settings.
templates.TEMPLATEFS = os.DirFS("templates/")
templates.BASE_TEMPLATE_SUFFIXES = []string{".tmpl"}
templates.BASE_TEMPLATE_DIRS = []string{"base"}
templates.TEMPLATE_DIRS = []string{"templates"}
templates.USE_TEMPLATE_CACHE = false
```
As you might see from the above code, this follows your file structure.
We do not have to define the regular template directories, but we do have to define the base template directories.
We define the regular directories when rendering them.
```bash
templates/
├── base
│   └── base.tmpl
└── app
    ├── index.tmpl
    └── user.tmpl
```
Then, we can render templates like so:
```go
func indexFunc(req *request.Request) {
    // Render the template with the given data.
    var err = req.Render("app/index.tmpl")
	if err != nil {
		req.WriteString(err.Error())
	}
}