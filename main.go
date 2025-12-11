package structcopygen

import (
	"flag"
	"fmt"
	"go/ast"
	"go/token"
	"log"
	"log/slog"
	"os"
	"path"
	"runtime"

	"github.com/spf13/pflag"
	"github.com/structcopy/structcopy-gen/config"
	"github.com/structcopy/structcopy-gen/internal/gen"
	"github.com/structcopy/structcopy-gen/internal/load"
	"golang.org/x/tools/go/packages"
)

func Run() error {
	cfg, err := config.LoadAppConfig("", "")
	if err != nil {
		log.Panic(err)
	}

	flagSet := pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)

	flagSet.BoolVarP(&cfg.CliFlags.Version, "version", "v", false, "Version")
	flagSet.BoolVarP(&cfg.CliFlags.Standalone, "standalone", "s", false, "Standalone mode")
	output := flagSet.StringP("out", "o", "", "Set the output file path")
	logs := flagSet.BoolP("log", "l", false, "Write log messages to <output path>.log.")
	dryRun := flagSet.BoolP("dry", "d", false, "Perform a dry run without writing files.")
	prints := flagSet.BoolP("print", "p", false, "Print the resulting code to STDOUT as well.")

	if err := flagSet.Parse(os.Args[1:]); err != nil {
		return err
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	if cfg.CliFlags.Version {
		fmt.Println(config.Version)
		fmt.Printf("%s__%s__%s__%s\n", config.Version, config.CommitHash, config.BuildTime, runtime.Version())
	} else if cfg.CliFlags.Standalone {
		inp := "examples/internal/standalone/structcopy-gen.go"
		ext := path.Ext(inp)
		out := inp[0:len(inp)-len(ext)] + ".gen" + ext

		pluginFunc := func(pkg *packages.Package, fset *token.FileSet, file *ast.File) error {
			g, err := gen.NewGenerator(
				pkg,
				fset,
				file,
				gen.WithInputPath(inp),
				gen.WithOutputPath(out),
				gen.WithLogger(logger),
				gen.WithDryRun(*dryRun),
				gen.WithPrints(*prints),
			)
			if err != nil {
				logger.Error("generate failed", slog.Any("error", err))
				return err
			}

			_, err = g.Generate(out, *prints, *dryRun)
			if err != nil {
				return err
			}

			return nil
		}

		err := load.LoadPackage(inp, out, pluginFunc)
		if err != nil {
			return err
		}

		logger.Info("Running")
		select {}
	} else {
		inputPath := flagSet.Arg(0)
		if inputPath == "" {
			// get Go file path which go:generate comment resides
			inputPath = os.Getenv("GOFILE")
		}
		if inputPath == "" {
			flag.Usage()
			os.Exit(1)
		}
		inp := inputPath
		out := ""
		log := ""

		if *output != "" {
			out = *output
		} else {
			ext := path.Ext(inputPath)
			out = inputPath[0:len(inputPath)-len(ext)] + ".gen" + ext
		}

		if *logs {
			ext := path.Ext(out)
			log = out[0:len(out)-len(ext)] + ".log"
		}

		pluginFunc := func(pkg *packages.Package, fset *token.FileSet, file *ast.File) error {
			g, err := gen.NewGenerator(
				pkg,
				fset,
				file,
				gen.WithInputPath(inp),
				gen.WithOutputPath(out),
				gen.WithLogPath(log),
				gen.WithLogEnabled(*logs),
				gen.WithDryRun(*dryRun),
				gen.WithPrints(*prints),
			)
			if err != nil {
				return err
			}

			_, err = g.Generate(out, *prints, *dryRun)
			if err != nil {
				return err
			}

			return nil
		}

		err := load.LoadPackage(inp, out, pluginFunc)
		if err != nil {
			return err
		}

		return nil
	}

	return nil
}
