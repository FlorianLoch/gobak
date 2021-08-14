package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	kongyaml "github.com/alecthomas/kong-yaml"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/florianloch/gobak/internal"
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

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	cfg := &cli.Config

	watcher, err := internal.NewPollWatcher(cfg)
	if err != nil {
		log.Error().Err(err).Msg("Failed to initialize PollWatcher")
	}

	if err := watcher.AddRecursively(cfg.Source); err != nil {
		log.Error().Err(err).Str("source_directory", cfg.Source).Msg("Failed to add source directory to watcher.")
	}

	copier := internal.NewCopier(cfg.Source, cfg.Target)

	gobak := internal.NewGobak(watcher, copier)

	gobak.Run()
}
