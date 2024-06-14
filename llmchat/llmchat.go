package llmchat

// 根据生成的代码 与大模型沟通设计架构
import (
	"chatcode/codeanalyze"
	"chatcode/knowledge"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/samber/mo"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

func ChatCode(filepath string) mo.Result[string] {
	ctx := context.Background()
	llm, err := openai.New()
	if err != nil {
		panic(err)
	}
	prompt := `你是一个高级代码架构师, 请按照指令理解下面的代码, 
	1. 理解模块结构,分析模块之间的依赖关系
	2. 理解代码要实现的目标
	3. 理解实现目标的实现路径
	<code>%s</code>

	下面是一个回复的 demo 
	<demo> 根据代码结构和内容,可以推断该项目的目标是实现一个聊天应用程序,其中包含了代码分析功能。实现路径是通过codeanalyze子模块中的函数来解析不同类型的代码文件,然后在llmchat模块中实现聊天功能。</demo>
	`

	var codeParser codeanalyze.GoParser
	module := codeParser.ParseDirectory(filepath)
	moduleStr, _ := json.Marshal(module)
	completion, err := llms.GenerateFromSinglePrompt(ctx, llm, fmt.Sprintf(prompt, moduleStr))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(completion)
	return mo.Ok[string](completion)
}

func ChatCodeDir(filepath string) mo.Result[string] {
	ctx := context.Background()
	llm, err := openai.New()
	if err != nil {
		panic(err)
	}
	prompt := `你是一个高级代码架构师, 请按照指令理解下面的代码文件夹, 
	1. 理解模块结构,分析模块之间的依赖关系
	2. 理解代码要实现的目标
	3. 理解实现目标的实现路径
	<code>%s</code>

	下面是一个回复的 demo 
	<demo> 根据代码结构和内容,可以推断该项目的目标是实现一个聊天应用程序,其中包含了代码分析功能。实现路径是通过codeanalyze子模块中的函数来解析不同类型的代码文件,然后在llmchat模块中实现聊天功能。</demo>
	`

	//	var codeParser codeanalyze.GoParser
	//	module := codeParser.ParseDirectory(filepath)
	//	moduleStr, _ := json.Marshal(module)
	dirStr := codeanalyze.DirTree(filepath, 4)
	completion, err := llms.GenerateFromSinglePrompt(ctx, llm, fmt.Sprintf(prompt, dirStr))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(completion)
	return mo.Ok[string](completion)
}

func ChatWithKnowledge(filepath string, question string) mo.Result[string] {
	ctx := context.Background()
	llm, err := openai.New()
	if err != nil {
		panic(err)
	}
	prompt := `你是一个高级有影响力的代码架构师, 请根据下面的步骤,来理解新手程序员提出的问题,并给出专业的答复 
	1. 根据 <code> 标记理解模块结构,分析程序模块之间的依赖关系
	2. 根据 <context> 标记中的代码上下文,
	3. 回答 <question> 中用户提出的问题
	<context> %s </context>
	<code> %s </code>
	<question> %s </question>`

	//	var codeParser codeanalyze.GoParser
	//	module := codeParser.ParseDirectory(filepath)
	//	moduleStr, _ := json.Marshal(module)
	dirStr := codeanalyze.DirTree(filepath, 4)

	contextStr := knowledge.BuildContext(knowledge.QueryEmbedding(question, 5, 0.6))
	prompt = fmt.Sprintf(prompt, contextStr, dirStr, question)
	fmt.Println(prompt)
	completion, err := llms.GenerateFromSinglePrompt(ctx, llm, prompt)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(completion)
	return mo.Ok[string](completion)
}
