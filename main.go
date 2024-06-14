package main

import (
	"chatcode/knowledge"
	"chatcode/llmchat"

	"github.com/alecthomas/kingpin/v2"
	"github.com/gookit/color"
	"github.com/gookit/gcli/v3/interact"
)

var (
	dir   = kingpin.Flag("dir", "dir of code.").Short('d').Required().String()
	build = kingpin.Flag("build", "build knowledge from dir").Short('b').Default("false").Bool()
)

// 了解指定项目代码的这些信息
// 1. 了解这个项目的背景和目的，能提供哪些功能，不能提供哪些功能
// 2. 为了实现这些功能，都通过哪些方法来做
// 3. 这些方法都依赖那些技术
// 交互流程
// 语义级别数据录入
// AI 半自动化流程
// AI 自动化梳理流程

// 对整体目录进行学习,列出整体的目录
func main() {
	kingpin.Parse()

	if *build {
		result := knowledge.CreateKnowledge(*dir)
		if result.IsError() {
			println(result.Error())
			return
		}
	}

	for {
		question, _ := interact.ReadLine("Your question? ")
		if question != "" {
			color.Println("Your input: ", question)
			color.Println(llmchat.ChatWithKnowledge(*dir, question).MustGet())
		} else {
			color.Cyan.Println("No input!")
		}
	}

}
