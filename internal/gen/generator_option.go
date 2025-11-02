package gen

import "log/slog"

// GeneratorOption is option for Generator.
type GeneratorOption func(g *Generator)

// WithLogger sets logger.
func WithLogger(logger *slog.Logger) GeneratorOption {
	return func(g *Generator) {
		g.logger = logger
	}
}

// WithInputPath sets input path.
func WithInputPath(input string) GeneratorOption {
	return func(g *Generator) {
		g.input = input
	}
}

// WithOutputPath sets output path.
func WithOutputPath(output string) GeneratorOption {
	return func(g *Generator) {
		g.output = output
	}
}

// WithLogPath sets output path.
func WithLogPath(log string) GeneratorOption {
	return func(g *Generator) {
		g.log = log
	}
}

// WithLogEnabled sets output path.
func WithLogEnabled(logs bool) GeneratorOption {
	return func(g *Generator) {
		g.logs = logs
	}
}

// WithPrints sets output path.
func WithPrints(prints bool) GeneratorOption {
	return func(g *Generator) {
		g.prints = prints
	}
}

// WithDryRun sets output path.
func WithDryRun(dryRun bool) GeneratorOption {
	return func(g *Generator) {
		g.dryRun = dryRun
	}
}
