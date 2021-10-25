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
		new(requestlogger.Module), // request logger show request logs
		new(gotemplate.Module),    // enables the gotemplate template engine module
		new(helloworld.Module),
	})
}
