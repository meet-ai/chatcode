package codeanalyze

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/samber/mo"
)

type GoParser struct {
}

func (pp *GoParser) ParseDirectory(path string) mo.Result[Module] {

	entries, err := os.ReadDir(path)
	if err != nil {
		return mo.Err[Module](err)
	}

	module := Module{Name: path}
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		if entry.IsDir() {
			subModule := pp.ParseDirectory(filepath.Join(path, entry.Name()))
			if subModule.IsError() {
				continue
			}
			module.Submodules = append(module.Submodules, subModule.MustGet())
		} else if filepath.Ext(entry.Name()) == ".go" {
			srcFile := pp.ParseFile(filepath.Join(path, entry.Name()))
			if srcFile.IsError() {
				continue
			}
			module.SrcFiles = append(module.SrcFiles, srcFile.MustGet())
		}
	}

	return mo.Ok[Module](module)
}

func (pp *GoParser) ParseFile(path string) mo.Result[SrcFile] {
	srcFile := SrcFile{SrcName: path}

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
	if err != nil {
		return mo.Err[SrcFile](err)
	}

	for _, decl := range node.Decls {
		switch decl := decl.(type) {
		case *ast.FuncDecl:
			funcName := decl.Name.Name
			// 不携带 func 内容，否则数据量过大
			var buf bytes.Buffer
			//			err := printer.Fprint(&buf, fset, decl.Body)
			//			if err != nil {
			//				println(err.Error())
			//			}

			srcFile.Funcs = append(srcFile.Funcs, Func{FuncName: funcName, FuncContent: buf.String()})
		case *ast.GenDecl:
			for _, spec := range decl.Specs {
				switch spec := spec.(type) {
				case *ast.TypeSpec:
					if class, ok := spec.Type.(*ast.StructType); ok {
						//className := spec.Name.Name
						var funcs []Func
						var members []string
						for _, field := range class.Fields.List {
							if field.Names != nil {
								memberName := field.Names[0].Name
								members = append(members, memberName)
							}
						}
						srcFile.Classes = append(srcFile.Classes, Class{Funcs: funcs, Members: members})
					}
				}
			}
		}
	}

	return mo.Ok[SrcFile](srcFile)
}
