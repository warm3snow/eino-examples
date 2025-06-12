package quickstart

import (
	"context"
	"log"

	"github.com/cloudwego/eino-ext/components/tool/duckduckgo"
	"github.com/cloudwego/eino-ext/components/tool/duckduckgo/ddgsearch"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
)

func getTodoTools() []tool.BaseTool {
	return []tool.BaseTool{
		getAddTodoTool(),
		getUpdateTodoTool(),
		&ListTodoTool{},
		// getSearchTool(),
	}
}

func getToolsInfo(tools []tool.BaseTool) []*schema.ToolInfo {
	toolsInfo := make([]*schema.ToolInfo, 0, len(tools))
	for _, t := range tools {
		info, err := t.Info(context.Background())
		if err != nil {
			panic(err)
		}
		toolsInfo = append(toolsInfo, info)
	}
	return toolsInfo
}

// 参数结构体
type TodoAddParams struct {
	Content   string `json:"content" jsonschema:"description=content of the todo,required=true"`
	StartTime *int64 `json:"start_time,omitempty" jsonschema:"description=start time in unix timestamp"`
	EndTime   *int64 `json:"end_time,omitempty" jsonschema:"description=end time in unix timestamp"`
}
type TodoUpdateParams struct {
	ID        string  `json:"id" jsonschema:"description=id of the todo"`
	Content   *string `json:"content,omitempty" jsonschema:"description=content of the todo"`
	StartTime *int64  `json:"start_time,omitempty" jsonschema:"description=start time in unix timestamp"`
	EndTime   *int64  `json:"end_time,omitempty" jsonschema:"description=deadline of the todo in unix timestamp"`
	Done      *bool   `json:"done,omitempty" jsonschema:"description=done status"`
}

// getAddTodoTool using NewTool to create a tool for adding todo items.
func getAddTodoTool() tool.InvokableTool {
	info := &schema.ToolInfo{
		Name: "add_todo",
		Desc: "添加待办事项",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"content": {
				Type:     schema.String,
				Desc:     "待办事项内容",
				Required: true,
			},
			"start_time": {
				Type: schema.Integer,
				Desc: "待办事项开始时间, Unix时间戳",
			},
			"end_time": {
				Type: schema.Integer,
				Desc: "待办事项结束时间",
			},
		}),
	}
	return utils.NewTool(info, func(_ context.Context, params *TodoAddParams) (any, error) {
		log.Printf("invoke tool add_todo: %+v", params)
		// 这里可以实现添加待办事项的逻辑
		return `{"msg": "待办事项已添加"}`, nil
	})
}

// getUpdateTodoTool using InferTool to create a tool for updating todo items.
func getUpdateTodoTool() tool.InvokableTool {
	updateTool, err := utils.InferTool(
		"update_todo",
		"更新待办事项",
		func(_ context.Context, params *TodoUpdateParams) (any, error) {
			log.Printf("invoke tool update_todo: %+v", params)
			// 这里可以实现更新待办事项的逻辑
			return `{"msg": "待办事项已更新"}`, nil
		},
	)
	if err != nil {
		panic(err)
	}
	return updateTool
}

// getListTodoTool
type ListTodoTool struct{}

func (l *ListTodoTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "list_todo",
		Desc: "列出待办事项",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"finished": {
				Type:     schema.Boolean,
				Desc:     "是否只列出已完成的待办事项",
				Required: false,
			},
		}),
	}, nil
}
func (l *ListTodoTool) InvokableRun(ctx context.Context, argumentsInJson string, opts ...tool.Option) (string, error) {
	log.Printf("invoke tool list_todo with arguments: %s", argumentsInJson)
	// 这里可以实现列出待办事项的逻辑
	return `{"todos": [{"id": "1", "content": "学习Go语言", "start_time": 1700000000, "end_time": 1700003600, "done": false}]}`, nil
}

// getSearchTool using official tool to create a tool for searching web content.
func getSearchTool() tool.InvokableTool {
	// 	accept: text/event-stream
	// accept-language: en-US,en;q=0.9
	// cache-control: no-cache
	// content-type: application/json
	// pragma: no-cache
	// User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36
	// origin: https://duckduckgo.com
	// referer: https://duckduckgo.com/
	// x-vqd-accept: 1
	config := &duckduckgo.Config{
		MaxResults: 1,
		DDGConfig: &ddgsearch.Config{
			Headers: map[string]string{
				"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36",
				"Accept":          "text/event-stream",
				"Accept-Language": "en-US,en;q=0.9",
				"Cache-Control":   "no-cache",
				"Content-Type":    "application/json",
				"Pragma":          "no-cache",
				"Origin":          "https://duckduckgo.com",
				"Referer":         "https://duckduckgo.com/",
				"x-vqd-accept":    "1",
			},
		},
	}
	searchTool, err := duckduckgo.NewTool(context.Background(), config)
	if err != nil {
		panic(err)
	}
	return searchTool
}
