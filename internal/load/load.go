package load

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"

	"golang.org/x/tools/go/packages"
)

const buildTag = "structcopygen"

// parserLoadMode is a packages.Load mode that loads types and syntax trees.
const parserLoadMode = packages.NeedName | packages.NeedImports | packages.NeedDeps |
	packages.NeedTypes | packages.NeedSyntax | packages.NeedTypesInfo

func LoadPackage(srcPath, dstPath string, fsetFunc func(pkg *packages.Package, fset *token.FileSet, file *ast.File) error) error {
	fileSet := token.NewFileSet()

	srcStat, err := os.Stat(srcPath)
	if err != nil {
		return err
	}

	dstStat, _ := os.Stat(dstPath)

	var parseErr error
	var fileSrc *ast.File
	cfg := &packages.Config{
		Mode:       parserLoadMode,
		BuildFlags: []string{"-tags", buildTag},
		Fset:       fileSet,
		ParseFile: func(fset *token.FileSet, filename string, src []byte) (*ast.File, error) {
			stat, err := os.Stat(filename)
			if err != nil {
				return nil, err
			}

			// If previously generation target file exists, skip reading it.
			if os.SameFile(stat, dstStat) {
				return nil, nil
			}

			if !os.SameFile(stat, srcStat) {
				return parser.ParseFile(fset, filename, src, 0)
			}

			file, err := parser.ParseFile(fset, filename, src, parser.ParseComments)
			if err != nil {
				parseErr = err
				return nil, err
			}

			fileSrc = file

			return file, nil
		},
	}
	pkgs, err := packages.Load(cfg, "file="+srcPath)
	if err != nil {
		// return nil, logger.Errorf("%v: failed to load type information: \n%w", srcPath, err)
		return err
	}
	if len(pkgs) == 0 {
		// return nil, logger.Errorf("%v: failed to load package information", srcPath)
		return err
	}

	if fileSrc == nil && parseErr != nil {
		// return nil, logger.Errorf("%v: %v", srcPath, parseErr)
		return parseErr
	}

	err = fsetFunc(pkgs[0], fileSet, fileSrc)
	if err != nil {
		return err
	}

	return nil
}
