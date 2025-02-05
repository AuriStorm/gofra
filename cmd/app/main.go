package main

import (
	"gofra/internal/app"
	"gofra/internal/utils"
)

func main() {
	// it's kinda we no need to set GOMAXPROCS manually anymore

	// I would prefer use https://github.com/urfave/cli/ to wrap cli commands & flags
	cliArgs := utils.MustReadArgs()

	// lets save some space to main here for pre run things like cache warms, daemon runs, etc

	app := app.New()
	app.Configure(cliArgs)
	app.RunServer()

	// and teardown clean-ups, offline processing hook invokes

}
