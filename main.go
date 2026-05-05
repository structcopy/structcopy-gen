package structcopygen

import (
	"flag"
	"fmt"
	"go/ast"
	"go/token"
	"log/slog"
	"os"
	"path"

	"github.com/structcopy/structcopy-gen/config"
	"github.com/structcopy/structcopy-gen/internal/gen"
	"github.com/structcopy/structcopy-gen/internal/load"
	"golang.org/x/tools/go/packages"
)

type App struct {
	cfg *config.AppConfig

	logEnabled   bool
	debugEnabled bool
}

func NewApp(cfg *config.AppConfig) (*App, error) {
	return &App{
		cfg:          cfg,
		logEnabled:   cfg.CliFlags.LogEnabled,
		debugEnabled: cfg.CliFlags.DebugEnabled,
	}, nil
}

func (a *App) Run() error {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	if a.cfg.CliFlags.Version {
		fmt.Println(config.Version)
		// fmt.Printf("%s__%s__%s__%s\n", config.Version, config.CommitHash, config.BuildTime, runtime.Version())
		// fmt.Printf("%s.%s\n", config.GetBuildInfoVersion(), config.GetBuildInfoRevision())
	} else if a.cfg.CliFlags.Standalone {
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
			)
			if err != nil {
				logger.Error("generate failed", slog.Any("error", err))
				return err
			}

			_, err = g.Generate(out, a.cfg.CliFlags.DebugEnabled, a.cfg.CliFlags.DryRun)
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
		inputPath := a.cfg.CliFlags.InputPath
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

		if a.cfg.CliFlags.OutputPath != "" {
			out = a.cfg.CliFlags.OutputPath
		} else {
			ext := path.Ext(inputPath)
			out = inputPath[0:len(inputPath)-len(ext)] + ".gen" + ext
		}

		if a.cfg.LogEnabled {
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
				gen.WithLogEnabled(a.cfg.LogEnabled),
			)
			if err != nil {
				return err
			}

			_, err = g.Generate(out, a.cfg.CliFlags.DebugEnabled, a.cfg.CliFlags.DryRun)
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
