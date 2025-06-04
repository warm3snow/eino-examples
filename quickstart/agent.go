package quickstart

import (
	"context"
	"log"
	"sync"

	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

var agent compose.Runnable[[]*schema.Message, []*schema.Message]
var once sync.Once

func componentCallbacks() callbacks.Handler {
	handler := callbacks.NewHandlerBuilder().
		OnStartFn(
			func(ctx context.Context, info *callbacks.RunInfo, input callbacks.CallbackInput) context.Context {
				log.Printf("onStart, runInfo: %v, input: %v", info, input)
				return ctx
			}).
		OnEndFn(
			func(ctx context.Context, info *callbacks.RunInfo, output callbacks.CallbackOutput) context.Context {
				log.Printf("onEnd, runInfo: %v, out: %v", info, output)
				return ctx
			}).
		Build()
	return handler
}

func initAgent(tools []tool.BaseTool) compose.Runnable[[]*schema.Message, []*schema.Message] {
	once.Do(func() {
		ctx := context.Background()

		// 获取todoTools
		todoToolInfos := getToolsInfo(tools)

		// 初始化ChatModel, 并绑定tools
		chatModel := CreateOllamaChatModel(ctx)
		err := chatModel.BindTools(todoToolInfos)
		if err != nil {
			panic(err)
		}

		// create tools node
		todoToolsNode, err := compose.NewToolNode(ctx, &compose.ToolsNodeConfig{
			Tools: tools,
		})
		if err != nil {
			panic(err)
		}
		// build chain
		chain := compose.NewChain[[]*schema.Message, []*schema.Message]()
		chain.
			AppendChatModel(chatModel, compose.WithNodeName("chat_model")).
			AppendToolsNode(todoToolsNode, compose.WithNodeName("tools"))

		// run chain
		a, err := chain.Compile(ctx)
		if err != nil {
			panic(err)
		}
		agent = a
	})
	return agent
}

func InvokeAgent(ctx context.Context, question string) {
	if agent == nil {
		agent = initAgent(getTodoTools())
	}

	messages := []*schema.Message{
		{
			Role:    schema.User,
			Content: question,
		},
	}
	resp, err := agent.Invoke(ctx, messages, compose.WithCallbacks(componentCallbacks()))
	if err != nil {
		panic(err)
	}
	for _, msg := range resp {
		log.Printf("\n")
		log.Printf("message: %s: %s", msg.Role, msg.Content)
	}
}
