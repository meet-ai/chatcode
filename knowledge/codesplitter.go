package knowledge

import (
	"strings"
)

type CodeSplitter struct {
}

// 应该拆分成多大的单位
func (csp CodeSplitter) SplitText(page string) ([]string, error) {
	var blocks []string
	buff := ""
	for _, line := range strings.Split(page, "\n") {
		if len(buff)+len(line)+1 > 4096 {
			blocks = append(blocks, buff)
			buff = line
		} else {
			buff += "\n" + line
		}
	}
	if buff != "" {
		blocks = append(blocks, buff)
	}
	return blocks, nil
}

func NewCodeSplitter() *CodeSplitter {
	return &CodeSplitter{}
}
