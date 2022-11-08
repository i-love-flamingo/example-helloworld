package helloworld

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework/web"

	"flamingo.me/example-helloworld/src/helloworld/interfaces"
)

// Module is our helloWorld Module
type Module struct{}

// Configure is the default Method a Module needs to implement
func (m *Module) Configure(injector *dingo.Injector) {
	// register our routes struct as a router Module - so that it is "known" to the router module
	web.BindRoutes(injector, new(routes))
}

// routes struct - our internal struct that gets the interface methods for router.Module
type routes struct {
	// helloController - we will defined routes that are handled by our HelloController - so we need this as a dependency
	helloController *interfaces.HelloController
}

// Inject dependencies - this is called by Dingo and gets an initializes instance of the HelloController passed automatically
func (r *routes) Inject(controller *interfaces.HelloController) *routes {
	r.helloController = controller

	return r
}

// Routes method which defines all routes handlers in module
func (r *routes) Routes(registry *web.RouterRegistry) {
	// Bind the path /hello to a handle with the name "hello"
	registry.MustRoute("/hello", "helloWorld.hello")
	// Bind the controller.Action to the handle "hello":
	registry.HandleGet("helloWorld.hello", r.helloController.Hello)

	registry.MustRoute("/greet", "helloWorld.greet")
	registry.HandleGet("helloWorld.greet", r.helloController.Greet)

	registry.MustRoute("/api", "helloWorld.api")
	registry.HandleGet("helloWorld.api", r.helloController.ApiHello)

	// Bind a route with a path parameter
	registry.MustRoute("/greet/:nickname", "helloWorld.greet")

	// Bind a route with a default value for a param
	registry.MustRoute("/greetflamingo", `helloWorld.greet(nickname="Flamingo")`)

	registry.HandleData("currenttime", r.helloController.CurrentTime)

}
