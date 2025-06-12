package quickstart

import (
	"context"
	"testing"

	"github.com/test-go/testify/assert"
)

func TestReact_Agent(t *testing.T) {
	ctx := context.Background()

	systemMessage := "你是一个用户信息查询助手，用户会提供姓名和邮箱，你需要根据这些信息查询用户的详细信息。"
	userMessages := []string{
		"请帮我查询用户信息，用户姓名是张三，邮箱是zhangsan@example.com",
	}
	// Invoke the React agent
	response, err := ReactAgent(ctx, systemMessage, userMessages)
	assert.NoError(t, err, "React agent should not return an error")
	assert.NotEmpty(t, response, "React agent response should not be empty")

	t.Logf("React agent response: %s", response)
}
