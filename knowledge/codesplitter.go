package knowledge

import (
	"fmt"
	"strings"
)

type CodeSplitter struct {
}

// 应该拆分成多大的单位
func (csp CodeSplitter) SplitText(page string) ([]string, error) {
	lines := strings.Split(page, "\n")
	var each4096Block []string
	var buff string
	for _, line := range lines {
		if len(buff) < 4096 {
			buff = fmt.Sprintf("%s\n%s", buff, line)
		} else {
			each4096Block = append(each4096Block, buff)
			buff = ""
		}
	}
	return each4096Block, nil
}

func NewCodeSplitter() *CodeSplitter {
	return &CodeSplitter{}
}
