package interfaces

import (
	"context"
	"time"

	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// Hello represents our first simple controller
	HelloController struct {
		responder *web.Responder
	}

	helloViewData struct {
		Name     string
		Nickname string
	}
)

// Inject dependencies
func (controller *HelloController) Inject(responder *web.Responder) *HelloController {
	controller.responder = responder

	return controller
}

// Get is a controller action that renders the `hello.html` template
func (controller *HelloController) Hello(_ context.Context, r *web.Request) web.Result {
	// Calling the Render method from the response helper and render the template "hello"
	return controller.responder.Render("hello", helloViewData{
		Name: "World",
	})
}

// Greet is a controller action that renders the `hello.html` template ands prints a provided URL param
func (controller *HelloController) Greet(_ context.Context, r *web.Request) web.Result {
	name, err := r.Query1("name")
	if err != nil {
		name = "World (default)"
	}

	nick, _ := r.Params["nickname"]

	return controller.responder.Render("hello", helloViewData{
		Name:     name,
		Nickname: nick,
	})
}

// HelloJSON is a controller action that renders Data
func (controller *HelloController) HelloJSON(_ context.Context, r *web.Request) web.Result {
	// Calling the Render method from the response helper and render the template "hello"
	return controller.responder.Data(helloViewData{
		Name: "World",
	})
}

// CurrentTime is a DataAction that handles data calls from templates
func (controller *HelloController) CurrentTime(ctx context.Context, r *web.Request, params web.RequestParams) interface{} {
	return time.Now().Format(time.RFC822)
}
