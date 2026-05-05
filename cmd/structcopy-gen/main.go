package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/pflag"
	structcopygen "github.com/structcopy/structcopy-gen"
	"github.com/structcopy/structcopy-gen/config"
)

func main() {
	cfg, err := config.LoadAppConfig("", "")
	if err != nil {
		log.Panic(err)
	}

	flagSet := pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)

	flagSet.BoolVarP(&cfg.CliFlags.Version, "version", "v", false, "Version")
	flagSet.BoolVarP(&cfg.CliFlags.Standalone, "standalone", "s", false, "Standalone mode")
	flagSet.StringVarP(&cfg.CliFlags.OutputPath, "output", "o", "", "Set the output file path")
	flagSet.BoolVarP(&cfg.CliFlags.LogEnabled, "log", "l", false, "Write log messages to <output path>.log.")
	flagSet.BoolVarP(&cfg.CliFlags.DebugEnabled, "debug", "p", false, "Print the resulting code to STDOUT as well.")
	flagSet.BoolVarP(&cfg.CliFlags.DryRun, "dry-run", "d", false, "Perform a dry run without writing files.")

	if err := flagSet.Parse(os.Args[1:]); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}

	if len(flagSet.Args()) > 0 {
		cfg.CliFlags.InputPath = flagSet.Arg(0)
	}

	if cfg.CliFlags.LogEnabled {
		cfg.LogEnabled = cfg.CliFlags.LogEnabled
	}

	app, err := structcopygen.NewApp(cfg)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}

	if err := app.Run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}
