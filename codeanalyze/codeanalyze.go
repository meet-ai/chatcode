package codeanalyze

import (
	"bytes"
	"errors"

	"github.com/samber/mo"
)

type CodeParser interface {
	ParseDirectory(string) mo.Result[Module]
	ParseFile(string) mo.Result[SrcFile]
}

func LoadCode(filepath string, lang string) mo.Result[Module] {
	switch lang {
	case "python":
		pyCodeParser := PythonParser{}
		return pyCodeParser.ParseDirectory(filepath)

	case "go":
		goCodeParser := GoParser{}
		return goCodeParser.ParseDirectory(filepath)
	}
	return mo.Err[Module](errors.New("failed"))
}

type Module struct {
	Name       string
	Submodules []Module
	SrcFiles   []SrcFile
}
type SrcFile struct {
	SrcName string
	Classes []Class
	Funcs   []Func
}

type Class struct {
	ClassName string
	Funcs     []Func
	Members   []string
}

type Func struct {
	FuncName    string
	FuncContent string
}

//	func ParseDirectory(path string) (Module, error) {
//		var module Module
//
//		entries, err := os.ReadDir(path)
//		if err != nil {
//			return module, err
//		}
//
//		for _, entry := range entries {
//			if entry.IsDir() {
//				subModule, err := ParseDirectory(filepath.Join(path, entry.Name()))
//				if err != nil {
//					return module, err
//				}
//				module.Submodules = append(module.Submodules, subModule)
//			} else if filepath.Ext(entry.Name()) == ".go" {
//				srcFile, err := parseFile(filepath.Join(path, entry.Name()))
//				if err != nil {
//					return module, err
//				}
//				module.SrcFiles = append(module.SrcFiles, srcFile)
//			}
//		}
//
//		return module, nil
//	}
//
//	func parseFile(path string) (SrcFile, error) {
//		var srcFile SrcFile
//
//		//	fset := token.NewFileSet()
//		//	node, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
//		//	if err != nil {
//		//		return srcFile, err
//		//	}
//
//		//for _, decl := range node.Decls {
//		//		switch decl := decl.(type) {
//		//		case *ast.FuncDecl:
//		//			funcName := decl.Name.Name
//		//			funcContent := decl.Body.String()
//		//			srcFile.Funcs = append(srcFile.Funcs, Func{FuncName: funcName, FuncContent: funcContent})
//		//		case *ast.GenDecl:
//		//			for _, spec := range decl.Specs {
//		//				switch spec := spec.(type) {
//		//				case *ast.TypeSpec:
//		//					if class, ok := spec.Type.(*ast.StructType); ok {
//		//						className := spec.Name.Name
//		//						var funcs []Func
//		//						var members []Member
//		//						for _, field := range class.Fields.List {
//		//							if field.Names != nil {
//		//								memberName := field.Names[0].Name
//		//								members = append(members, Member{MemberName: memberName})
//		//							}
//		//						}
//		//						srcFile.Classes = append(srcFile.Classes, Class{Funcs: funcs, Members: members})
//		//					}
//		//				}
//		//			}
//		//		}
//		//}
//
//		return srcFile, nil
//	}
func DirTree(path string, depth int) string {
	var treeBuff = bytes.NewBuffer([]byte(""))
	dirTree(treeBuff, path, true, depth)
	return treeBuff.String()
}
