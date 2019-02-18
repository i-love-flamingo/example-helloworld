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

// routes struct - our internal struct that gets the interface methods for router.Module
type routes struct {
	// helloController - we will defined routes that are handled by our HelloController - so we need this as a dependency
	helloController *interfaces.HelloController
}

// Inject method - this is called by Dingo and gets an initializes instance of the HelloController passed automatically
func (r *routes) Inject(controller *interfaces.HelloController) {
	r.helloController = controller
}

// Routes method which defines all routes handlers in module
func (r *routes) Routes(registry *web.RouterRegistry) {
	// Bind the path /hello to a handle with the name "hello"
	registry.Route("/hello", "hello")

	// Bind the controller.Action to the handle "hello":
	registry.HandleGet("hello", r.helloController.Get)

	registry.HandleGet("helloWorld.greetme", r.helloController.GreetMe)
	registry.Route("/greetme", "helloWorld.greetme")
	registry.Route("/greetme/:nickname", "helloWorld.greetme")
	registry.Route("/greetflamingo", `helloWorld.greetme(nickname="Flamingo")`)

	registry.HandleData("currenttime", r.helloController.CurrentTime)
}
