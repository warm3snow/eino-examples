package quickstart

import (
	"context"
	"log"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
)

func ReactAgent(ctx context.Context, systemMessage string, userMessages []string) (string, error) {
	// ollama chatModel
	chatModel := CreateOllamaChatModel(ctx)

	// toolsNodeConfig
	tools := compose.ToolsNodeConfig{
		Tools: getUserInfoTool(),
	}

	// create agent
	agent, err := react.NewAgent(ctx, &react.AgentConfig{
		ToolCallingModel: chatModel,
		ToolsConfig:      tools,
		MessageModifier: func(ctx context.Context, input []*schema.Message) []*schema.Message {
			res := make([]*schema.Message, 0, len(input)+1)
			res = append(res, schema.SystemMessage(systemMessage))
			res = append(res, input...)
			return res
		},
		MaxStep: 12, // max tool call is 5. 12/2 - 1 = 5
	})
	if err != nil {
		return "", err
	}

	// run agent
	var messages []*schema.Message
	for _, userMessage := range userMessages {
		messages = append(messages, schema.UserMessage(userMessage))
	}
	// response, err := agent.Generate(context.Background(), messages)
	// if err != nil {
	// 	return "", err
	// }

	modelCallbackHandler := componentCallbacks()

	chain := compose.NewChain[[]*schema.Message, string]()
	agentLambda, _ := compose.AnyLambda(agent.Generate, agent.Stream, nil, nil)

	chain.
		AppendLambda(agentLambda).
		AppendLambda(compose.InvokableLambda(func(ctx context.Context, input *schema.Message) (string, error) {
			log.Printf("got agent response: %s\n", input.Content)
			return input.Content, nil
		}))

	r, _ := chain.Compile(ctx)
	res, _ := r.Invoke(ctx, messages,
		compose.WithCallbacks(modelCallbackHandler),
	)

	return res, nil
}

func getUserInfoTool() []tool.BaseTool {
	type userInfoRequest struct {
		Name  string `json:"name" jsonschema:"description=用户姓名,required=true"`
		Email string `json:"email" jsonschema:"description=用户邮箱,required=true"`
	}

	type userInfoResponse struct {
		Name     string `json:"name" jsonschema:"description=用户姓名,required=true"`
		Email    string `json:"email" jsonschema:"description=用户邮箱,required=true"`
		Company  string `json:"company,omitempty" jsonschema:"description=用户公司,required=false"`
		Position string `json:"position,omitempty" jsonschema:"description=用户职位,required=false"`
		Salary   string `json:"salary,omitempty" jsonschema:"description=用户薪资,required=false"`
	}

	userInfoTool := utils.NewTool(
		&schema.ToolInfo{
			Name: "user_info",
			Desc: "根据用户的姓名和邮箱，查询用户的公司、职位、薪酬信息",
			ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
				"name": {
					Type: "string",
					Desc: "用户的姓名",
				},
				"email": {
					Type: "string",
					Desc: "用户的邮箱",
				},
			}),
		},
		func(ctx context.Context, input *userInfoRequest) (output *userInfoResponse, err error) {
			return &userInfoResponse{
				Name:     input.Name,
				Email:    input.Email,
				Company:  "Cool Company LLC.",
				Position: "CEO",
				Salary:   "9999",
			}, nil
		})

	return []tool.BaseTool{
		userInfoTool,
		// Add other tools here if needed
	}
}
