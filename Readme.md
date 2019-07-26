# Flamingo Example Helloworld

## Quick reference

### Precondition
* Go >= 1.11.4 installed

### Run the example

```bash
git clone git@github.com:i-love-flamingo/example-helloworld.git
cd example-helloworld
go run main.go serve
```
Open http://localhost:3322

## Tutorial Steps

### Start Over

When you checkout the example you already see a very basic flamingo module.

```bash
git checkout start-over
```

Here are some details for the files present in your project ( `flamingo.me/example-helloworld` )

#### Folder structure
First, we start with some basic information on the folder structure:

* config/
    * Here we ﬁnd our configuration and routing files
* src/
    * Source code for the project
* templates/
    * Templates for Go Templates
* go.mod
    * the projects dependencies
* main.go
    * Entry point for our project


The `main.go` file looks like this

```go
package main

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3"
	"flamingo.me/flamingo/v3/core/gotemplate"
	"flamingo.me/flamingo/v3/core/requestlogger"
	"flamingo.me/flamingo/v3/core/zap"
)

// main is our entry point
func main() {
	flamingo.App([]dingo.Module{
		new(zap.Module),           // log formatter
		new(requestlogger.Module), // requestlogger show request logs
		new(gotemplate.Module),    // gotemplate installs a go template engine (in debug mode, todo fix this)
	})
}
```

Also we need a few default templates, they live in `templates/error`:

`403.html`
```html
<h1>Forbidden</h1>
<hr/>
<pre>{{index . "error"}}</pre>
```

`404.html`
```html
<h1>File Not Found</h1>
<hr/>
<pre>{{index . "error"}}</pre>
```

`503.html`
```html
<h1>Unavailable</h1>
<hr/>
<pre>{{index . "error"}}</pre>
```

`withCode.html`
```html
<h1>Error {{index . "code"}}</h1>
<hr/>
<pre>{{index . "error"}}</pre>
```

We need this setup to make the `gotemplate` Module not panic due to missing error-templates at all. 
(See Default configuration of the InitModule in Flamingo's `framework/module.go`)

The config file `config/config.yml` is nearly empty for now. We just enable Flamingo's debug mode.

Now we are ready and can already start with `go run main.go`!

You can start the server with `go run main.go serve` but you'll be stuck with 404 errors for now. (Obviously, since we do not have any routes registered.)

The flamingo default app runs on port 3322, so go and visit http://localhost:3322/

You'll see log-output like
```
2019-07-26T13:42:10.073+0200    INFO    requestlogger/logger.go:131     GET / 404: 42b in 5.788304ms (Error: action for method "GET" not found and no any fallback)     {"area": "root", "traceID": "7158d0b6f0214ac018ccd0b6241da2f3", "spanID": "f16b1510928ff8df", "method": "GET", "path": "/", "client_ip": "[::1]:52572", "businessId": "", "accesslog": 1, "response_code": 404, "response_time": 0.005788304, "referer": ""}
```

In Step 1, we will make sure that the 404 error won't stay for long.

### Step 1

You can either start with your current code or you can just checkout the branch "start-over" with

```bash
git checkout start-over
```

At first, we create a template to be shown when visiting the page.
Create a template file called `index.html` in the `templates` folder with basic content:
```html
<html>
    <body>
        <h1>Hello World!</h1>
        <h2>This is the index page</h2>
    </body>
</html>
```

But to see this page in the browser, we need some routing.

Create the file `config/routes.yml` and add the route:
```yaml
- path: /
  name: index
  controller: flamingo.render(tpl="index")
```
This is the basic route on path "/" and it renders the "index" template via the builtin flamingo.render controller.

#### Builtin controllers
Flamingo comes with a couple of builtin controllers, that we can use in `routes.yml`. E.g.:

* flamingo.render(tpl="...")
    * Default template renderer
* flamingo.redirect(to="...") and flamingo.redirectPermanent(to="...")
    * Default "redirect to other route" controller
* flamingo.redirectUrl(url="...") and flamingo.redirectPermanentUrl(url="...")
    * Default "redirect to a URL" controller
* flamingo.error and flamingo.notfound
    * Default Flamingo error and notfound controller

Now, you can start the server again with `go run main.go serve` and see the result in the browser when you visit http://localhost:3322/.

If something doesn't work, you can always compare your code with the master branch.

### Step 2

#### Custom Controller

In this step we will create our first own controller and route.

* In Flamingo this is done with a custom module.
* We place our custom module in `src/helloworld/module.go`

We will now learn step by step how to:

1. Kickstart a new local Flamingo module `helloworld`
2. Write a simple Controller
3. Register that Controller and a new route in the module
4. Add another template

Let us start with the `helloworld` module:

Create a new file `module.go` inside  `src/helloworld` with:

```go
package helloworld

import (
	"flamingo.me/dingo"
)

// Module is our helloWorld Module
type Module struct {}

// Configure is the default Method a Module needs to implement
func (m *Module) Configure(injector *dingo.Injector) {
	// ...
}
```

Our Module is a `dingo.Module`, which can be used in our Flamingo project. To load our own Module we need to add a line
in the main.go file. So open your project `main.go` and add the new Module to the Bootstrap, it should look like this:

```go
package main

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3"
	"flamingo.me/flamingo/v3/core/gotemplate"
	"flamingo.me/flamingo/v3/core/requestlogger"
	"flamingo.me/flamingo/v3/core/zap"
	"flamingo.me/example-helloworld/src/helloworld"
)

// main is our entry point
func main() {
	flamingo.App([]dingo.Module{
		new(zap.Module),           // log formatter
		new(requestlogger.Module), // requestlogger show request logs
		new(gotemplate.Module),    // gotemplate installs a go template engine
		new(helloworld.Module),
	})
}
```

Now we create an own controller, first we start with some information about controllers:

* A Controller in Flamingo is a struct with a couple of methods, used as „Actions“
* Controllers can use the `*web.Responder` to create responses
* We can also simply use standard http.Handler from Go's stdlib

Create a folder `interfaces` in your `helloworld` Module, then create a go file called `hello_controller.go` in the 
`interfaces` folder with the following content:

```go
package interfaces

import (
	"context"

	"flamingo.me/flamingo/v3/framework/web"
)

type HelloController struct {
	responder *web.Responder
}

func (controller *HelloController) Inject(responder *web.Responder) *HelloController {
	controller.responder = responder

	return controller
}

func (controller *HelloController) Get(ctx context.Context, r *web.Request) web.Result {
	// Calling the Render method from the response helper and render the template "hello"
	return controller.responder.Render("hello", nil)
}
```

Our Controller reacts on GET requests and renders templates.

Beside `Render()` there are more functions such as `Data()`, `RouteRedirect()` and `NotFound()`.

Now, we have to register our controller and a route in our `module.go` file:

```go
package helloworld

import (
	"flamingo.me/dingo"
	"flamingo.me/example-helloworld/src/helloWorld/interfaces"
	"flamingo.me/flamingo/v3/framework/web"
)

// Module is our helloWorld Module
type Module struct{}

// Configure is the default Method a Module need to implement
func (m *Module) Configure(injector *dingo.Injector) {
	//Call Bind helper of router Module
	// It is a shortcut for: injector.BindMulti((*router.Module)(nil)).To(new(routes))
	// So what it does is register our routes struct as a router Module - so that it is "known" to the router module
	web.BindRoutes(injector, new(routes))
}

// routes struct that gets the interface methods for router.Module
type routes struct {
	// helloController - we will defined routes that are handled by our HelloController - so we need this as a dependency
	helloController *interfaces.HelloController
}

// Inject method is called by Dingo and gets an initialized instance of the HelloController passed automatically
func (r *routes) Inject(controller *interfaces.HelloController) *routes {
	r.helloController = controller

	return r
}

// Routes method which defines all routes handlers in module
func (r *routes) Routes(registry *web.RouterRegistry) {
	// Bind the path /hello to a handle with the name "hello"
	registry.Route("/hello", "hello")

	// Bind the controller.Action to the handle "hello":
	registry.HandleGet("hello", r.helloController.Get)
}
```

We also need a matching template, so we create a template for our own controller.
Create a template file called `hello.html` in the `templates` folder with some simple content:

```html
<h1>Hello Controller</h1>
```

Now, you can start the server again with `go run main.go serve` and on http://localhost:3322/hello we will now 
have our custom controller, that renders the new template.

If something doesn't work, you can always compare your code with the master branch.

### Step 3

#### Passing data to templates

* It is important for controllers to be able to pass data into templates
* This is done via the render call, as the second parameter (which is currently nil)

First we create a new struct, helloViewData, which we use as the data transfer object for the data that should be available in the template.

We will then pass this data to our template.

For all that, we need to navigate to the file src/helloworld/interfaces/helloController.go and change it as follows:

```go
package interfaces

import (
	"context"

	"flamingo.me/flamingo/v3/framework/web"
)

type (
	HelloController struct {
		responder *web.Responder
	}

	helloViewData struct {
		Name string
	}
)

func (controller *HelloController) Inject(responder *web.Responder) *HelloController {
	controller.responder = responder

	return controller
}

func (controller *HelloController) Get(ctx context.Context, r *web.Request) web.Result {
	// Calling the Render method from the response helper and render the template "hello"
	return controller.responder.Render("hello", helloViewData{
		Name: "World",
	})
}
```

As you can see in the `Get` function, we pass the variable `Name` with the value `World` to the ´hello.html` template.
To use that variable now in our `hello.html` template, we need to change the template accordingly:
```gotemplate
<h1>Hello {{ .Name }}</h1>
```
Now, you can start the server again with `go run main.go serve` and see the result in the browser when you visit http://localhost:3322/hello.

Now http://localhost:3322/hello will show "Hello World", where the "World" part is the string passed in the controller.

If something doesn't work, you can always compare your code with the master branch.

### Step 4

#### Handle URL parameters

Now it's time to handle URL Parameters

There are 3 ways of getting data into the controller:
* POST/Form data
* GET Parameter
* URL Path parameter

Create a new Action `GreetMe` in your `src/helloworld/interfaces/helloController.go` that will use URL parameters.

The Controller uses c.Query1 to get the first GET Query parameter with the key `name`.

It will return an error if not set, so we can default it to something we want if it's not set.

We will let it default to "World (default)" for now:

```go
func (controller *HelloController) GreetMe(ctx context.Context, r *web.Request) web.Result {
	name, err := r.Query1("name")
	if err != nil {
		name = "World (default)"
	}
	return controller.responder.Render("hello", helloViewData{
		Name: name,
	})
}
```

Now we need to add new `greetme` routes in the `src/helloworld/module.go`:

```go
// Routes method which defines all routes handlers in module
func (r *routes) Routes(registry *router.Registry) {
	//Bind the controller.Action to the handle "hello":
	registry.HandleGet("helloWorld.hello", r.helloController.Get)
	//Bind the path /hello to a handle with the name "hello"
	registry.Route("/hello", "helloWorld.hello")

	registry.HandleGet("helloWorld.greetme", r.helloController.GreetMe)
	registry.Route("/greetme", "helloWorld.greetme")
}
```

Now, you can start the server again with `go run main.go serve` and see the result in the browser when you visit http://localhost:3322/greetme?name=Flamingo.

Now http://localhost:3322/greetme?name=Flamingo will show "Hello Flamingo", where the "Flamingo" part is the "name" parameter passed from the URL via the controller to the template.

#### Path Parameters

Beside "GET" parameters we can also add "Path" parameters
* Extend the module.go, and add another route:
```go
registry.Route("/greetme/:nickname", "helloWorld.greetme")
```
* We can also set default parameters for routes like this:
```go
registry.Route("/greetflamingo", `helloWorld.greetme(nickname="Flamingo")`)
```

Add both routes to the `src/helloworld/module.go` file.

So the routes in the `src/helloworld/module.go` now should look like this:
```go
func (r *routes) Routes(registry *router.Registry) {
	//Bind the controller.Action to the handle "hello":
	registry.HandleGet("helloWorld.hello", r.helloController.Get)
	//Bind the path /hello to a handle with the name "hello"
	registry.Route("/hello", "helloWorld.hello")

	registry.HandleGet("helloWorld.greetme", r.helloController.GreetMe)
	registry.Route("/greetme", "helloWorld.greetme")
	registry.Route("/greetme/:nickname", "helloWorld.greetme")
	registry.Route("/greetflamingo", `helloWorld.greetme(nickname="Flamingo")`)
}
```

Additionally, we need to adjust our controller so we can get the `nickname` path variable via `r.Params`.

Go and change it accordingly:
```go
func (controller *HelloController) GreetMe(ctx context.Context, r *web.Request) web.Result {
	name, err := r.Query1("name")
	if err != nil {
		name = "World (default)"
	}

	nick, _ := r.Params["nickname"]

	return controller.responder.Render("hello", helloViewData{
		Name: name,
		Nickname: nick,
	})
}
```

We also need to extend the helloViewData struct, so add the "Nickname":

```go
helloViewData struct {
    Name     string
    Nickname string
}
```

Change the "hello.html" template to include the nickname:
```gotemplate
<h1>Hello {{ .Name }}</h1>
{{ if .Nickname }}
    <h2>Your Nickname is {{ .Nickname }} </h2>
{{end}}
```

Now, you can start the server again with `go run main.go serve` and see the result in the browser when you visit:

* http://localhost:3322/greetme/awesome?name=Flamingo
    * This should show "Flamingo" as "Name" and "awesome" as "Nickname".
* http://localhost:3322/greetflamingo
    * This should show "World (default)" as "Name" and "Flamingo" as "Nickname".

If something doesn't work, you can always compare your code with the master branch.

### Step 5

#### Requesting Data from DataActions

One of the important things for templates is to be able to "ask for data". There are cases where a controller can not simply create all data for a view, instead the view/template should be able to request data.
DataActions work similar to normal Actions, but are called from within templates. So let us check how we can request data.

First in our templates/index.html ﬁle we call the templatefunc „data“ to access a Datacontroller:

```gotemplate
<html>
    <body>
        <h1>Hello World!</h1>
        <h2>This is the index page</h2>
        <p>Currenttime: {{ data "currenttime" }}</p>
    </body>
</html>
```

Running `go run main.go serve` and opening  http://localhost:3322/  will now show an exception, because we do not 
actually have a data controller with the name `currenttime`. We need to add a data action to our helloController:

```go
// CurrentTime is a DataAction that handles data calls from templates
func (controller *HelloController) CurrentTime(ctx context.Context, r *web.Request, callParams web.RequestParams) interface{} {
	return time.Now().Format(time.RFC822)
}
```

After that, we register it as Data Handler in our `module.go`, but without a route (though we could also route it if necessary):

```go
registry.HandleData("currenttime", r.helloController.CurrentTime)
```

So the routes part in the module.go should look like this:

```go
// Routes method which defines all routes handlers in module
func (r *routes) Routes(registry *router.Registry) {
	//Bind the controller.Action to the handle "hello":
	registry.HandleGet("helloWorld.hello", r.helloController.Get)
	//Bind the path /hello to a handle with the name "hello"
	registry.Route("/hello", "helloWorld.hello")

	registry.HandleGet("helloWorld.greetme", r.helloController.GreetMe)
	registry.Route("/greetme", "helloWorld.greetme")
	registry.Route("/greetme/:nickname", "helloWorld.greetme")
	registry.Route("/greetflamingo", `helloWorld.greetme(nickname="Flamingo")`)

	registry.HandleData("currenttime", r.helloController.CurrentTime)
}
```

Now we can run flamingo again and we have the data from the data action available in our template. The DataAction has 
access to the request and session data, so it can return everything necessary.

If something doesn't work, you can always compare your code with the master branch.

Congratulations! You completed all steps in the Flamingo Helloworld Example and learned some basic Flamingo features.
