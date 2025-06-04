package quickstart

import (
	"context"
	"log"
	"testing"
)

func TestAgent_TodoTool(t *testing.T) {
	ctx := context.Background()
	InvokeAgent(ctx, "添加一个学习Go语言的待办事项，内容是：学习Go语言的并发编程，开始时间是明天早上9点，结束时间是明天下午5点。")
	InvokeAgent(ctx, "列出所有待办事项")
	InvokeAgent(ctx, "更新待办事项，ID为1，内容是：学习Go语言的并发编程，开始时间是明天早上9点，结束时间是明天下午5点。")
	InvokeAgent(ctx, "列出所有待办事项")
	log.Println("Agent invoked successfully")
}

func TestAgent_SearchTool(t *testing.T) {
	ctx := context.Background()
	InvokeAgent(ctx, "搜索关于Go语言并发编程的文章, 并添加到待办事项中用于之后的学习，并设置结束时间为后天下午5点。")
	log.Println("Agent invoked successfully")
}
