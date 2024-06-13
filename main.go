package main

import (
	"github.com/alecthomas/kingpin/v2"
)

var (
	dir = kingpin.Flag("dir", "dir of code.").Short('d').Required().String()
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
	//llmchat.ChatCodeDir(*dir)

}
