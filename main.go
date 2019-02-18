package main

import (
	"flamingo.me/dingo"
	"flamingo.me/example-helloworld/src/helloworld"
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
		new(helloworld.Module),
	})
}
