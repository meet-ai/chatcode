package codeanalyze

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
)

const UiRef = "├───"
const UiEndRef = "└───"
const UiVerticalLine = "│  "
const UiIndent = "  "

type NodeDependencies struct {
	name      string
	fullPath  string
	dependent *[]NodeDependencies
	isDir     bool
	size      int64
}

type TreeBuilder struct {
	isPrintFile bool
	depthLevel  int
}

func (tb *TreeBuilder) BuildTree(path string) (*NodeDependencies, error) {
	rootFile := &NodeDependencies{fullPath: path, isDir: true}
	if errorHappen := tb.buildTreeRecursive(rootFile, 0); errorHappen != nil {
		return nil, errorHappen
	}
	return rootFile, nil
}

func (tb *TreeBuilder) buildTreeRecursive(parent *NodeDependencies, depthLevel int) error {
	if tb.depthLevel > 0 && depthLevel >= tb.depthLevel {
		return nil
	}
	files, err := ioutil.ReadDir(parent.fullPath)
	if err != nil {
		parent.name += "(?error)"
		//noinspection GoUnhandledErrorResult
		fmt.Fprintln(os.Stderr, err)
		return err
	}

	els := make([]NodeDependencies, 0, len(files))
	index := 0
	for _, file := range files {
		if !tb.isPrintFile && !file.IsDir() {
			continue
		}
		els = append(els, NodeDependencies{
			name:     file.Name(),
			fullPath: parent.fullPath + string(os.PathSeparator) + file.Name(),
			isDir:    file.IsDir(),
			size:     file.Size(),
		})
		index += 1
	}
	(*parent).dependent = &els

	for idx := range els {
		var node = &els[idx]
		if !node.isDir {
			continue
		}
		if err := tb.buildTreeRecursive(node, depthLevel+1); err != nil {
			//return err
		}
	}
	return nil
}

func (tb *TreeBuilder) PrintTree(out io.Writer, nodes *[]NodeDependencies, prefix string) error {
	if nodes == nil {
		return nil
	}
	sort.Slice(*nodes, func(i, j int) bool {
		return (*nodes)[i].name < (*nodes)[j].name
	})
	size := len(*nodes)
	for idx := range *nodes {
		node := (*nodes)[idx]

		nodeSize := tb.getFileSizeSuffix(node)
		//is last
		if idx == size-1 {
			line := fmt.Sprintf("%s%s %s%s\n", prefix, UiEndRef, node.name, nodeSize)
			_, errWrite := fmt.Fprint(out, line)
			if errWrite != nil {
				return errWrite
			}

			errPrint := tb.PrintTree(out, node.dependent, prefix+UiIndent)
			if errPrint != nil {
				return errPrint
			}
		} else {
			line := fmt.Sprintf("%s%s %s%s\n", prefix, UiRef, node.name, nodeSize)
			_, errWrite := fmt.Fprint(out, line)
			if errWrite != nil {
				return errWrite
			}
			errPrint := tb.PrintTree(out, node.dependent, prefix+UiVerticalLine)
			if errPrint != nil {
				return errPrint
			}
		}
	}
	return nil
}

func (tb *TreeBuilder) getFileSizeSuffix(node NodeDependencies) string {
	var nodeSize = ""
	if !node.isDir {
		nodeSize = " (empty)"
		if node.size > 0 {
			nodeSize = fmt.Sprintf(" (%db)", node.size)
		}
	}
	return nodeSize
}

func dirTree(out io.Writer, path string, printFiles bool, depth int) error {
	treeBuilder := &TreeBuilder{printFiles, depth}
	rootNode, err := treeBuilder.BuildTree(path)
	if err != nil {
		return err
	}
	return treeBuilder.PrintTree(out, rootNode.dependent, "")
}
