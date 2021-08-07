package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	kongyaml "github.com/alecthomas/kong-yaml"

	"github.com/florianloch/gobak/internal/config"
)

type cli struct {
	Config config.Config `embed:""`
}

func main() {
	cli := &cli{}

	kong.ConfigureHelp(kong.HelpOptions{Compact: false, Summary: true})

	kong.Name("gobak")
	kong.Description("Small toll watching an input directory for changes and copying files over to a backup destination.\n" +
		"Not synchronizing, only mirroring.")

	ctx := kong.Parse(cli, kong.Configuration(kongyaml.Loader, "gobak.config.yaml"))
	if ctx.Error != nil {
		fmt.Println("Failed to parse input parameters/commands: " + ctx.Error.Error())

		os.Exit(1)
	}
}
