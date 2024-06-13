package codeanalyze

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/samber/mo"
)

type PythonParser struct {
}

func (pp *PythonParser) ParseDirectory(path string) mo.Result[Module] {
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
		} else if filepath.Ext(entry.Name()) == ".py" {
			srcFile := pp.ParseFile(filepath.Join(path, entry.Name()))
			if srcFile.IsError() {
				continue
			}
			module.SrcFiles = append(module.SrcFiles, srcFile.MustGet())
		}
	}

	return mo.Ok[Module](module)
}

func (pp *PythonParser) ParseFile(path string) mo.Result[SrcFile] {
	// parseFile parses a Python file and extracts classes, functions, and methods.
	return parsePythonFile(path)
}

func parsePythonFile(filePath string) mo.Result[SrcFile] {
	file, err := os.Open(filePath)
	if err != nil {
		mo.Err[SrcFile](err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var srcFile SrcFile

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "class ") {
			var class Class
			class.ClassName = line
			srcFile.Classes = append(srcFile.Classes, class)
		} else if strings.HasPrefix(line, "def ") {
			var func_ Func
			func_.FuncName = line
			srcFile.Funcs = append(srcFile.Funcs, func_)
		}
	}

	if err := scanner.Err(); err != nil {
		return mo.Err[SrcFile](err)
	}

	return mo.Ok[SrcFile](srcFile)
}
