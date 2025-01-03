package generatecmd

import (
	"bytes"
	"context"
	"go/ast"
	"go/parser"
	"go/token"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/garrettladley/templ-static/internal/meta"
	"github.com/garrettladley/templ-static/internal/xast"
)

func NewGenerate(log *slog.Logger, args Arguments) (g *Generate, err error) {
	g = &Generate{
		Log:  log,
		Args: &args,
	}
	return g, nil
}

type Generate struct {
	Log  *slog.Logger
	Args *Arguments
}

func (cmd Generate) Run(ctx context.Context) (err error) {
	components, err := findStaticComponents(cmd.Args.Path)
	if err != nil {
		panic(err)
	}
	for _, comp := range components {
		println(comp.FilePath + ": " + comp.FilePath + " -> " + comp.Path)
	}
	return
}

func findStaticComponents(root string) ([]meta.Meta, error) {
	var components []meta.Meta
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// only process templ-generated files
		if !strings.HasSuffix(path, "_templ.go") {
			return nil
		}

		// Read the file
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// check if file has static directive and get info
		meta, found := meta.Extract(bytes.NewReader(content), path)
		if !found {
			return nil
		}

		// parse the file
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, path, content, parser.ParseComments)
		if err != nil {
			return err
		}

		// find all functions that return templ.Component
		ast.Inspect(f, func(n ast.Node) bool {
			fn, ok := n.(*ast.FuncDecl)
			if !ok {
				return true
			}

			if xast.IsTemplComponent(fn.Type.Results) {
				components = append(components, meta)
			}
			return true
		})

		return nil
	})

	return components, err
}
