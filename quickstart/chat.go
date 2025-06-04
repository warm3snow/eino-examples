/**
 * @Author: xueyanghan
 * @File: chat.go
 * @Version: 1.0.0
 * @Description: desc.
 * @Date: 2025/6/4 13:51
 */

package quickstart

import (
	"context"
	"log"

	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

var (
	template = prompt.FromMessages(schema.GoTemplate,
		schema.SystemMessage("你是一个{{.Role}}。你需要用{{.Style}}的语气回答问题。你的目标是帮助程序员保持积极乐观的心态，提供技术建议的同时也要关注他们的心理健康。"),
		schema.MessagesPlaceholder("chat_history", true),
		schema.UserMessage("问题: {{.Question}}"),
	)
)

func Chat(role, style, question string, chat_history []string) {
	// 1. create message, use GoTemplate
	messages := CreateMessages(role, style, question, chat_history)

	// 2. create chatModel of Ollama
	chatModel := CreateOllamaChatModel(context.Background())

	// 3. run chatModel
	response, err := chatModel.Generate(context.Background(), messages)
	if err != nil {
		panic(err)
	}
	// 4. print response
	log.Printf("Chat response: %v", response)
}

func CreateMessages(role, style, question string, chat_history []string) []*schema.Message {
	chatHistoryMessages := make([]*schema.Message, 0, len(chat_history)/2)
	for i := 0; i < len(chat_history); i += 2 {
		chatHistoryMessages = append(chatHistoryMessages, schema.UserMessage(chat_history[i]))
		chatHistoryMessages = append(chatHistoryMessages, schema.AssistantMessage(chat_history[i+1], nil))
	}
	var vs = map[string]any{
		"Role":         role,
		"Style":        style,
		"Question":     question,
		"chat_history": chatHistoryMessages,
	}
	messages, err := template.Format(context.Background(), vs)
	if err != nil {
		panic(err)
	}
	return messages
}

func CreateOllamaChatModel(ctx context.Context) *ollama.ChatModel {
	chatModel, err := ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
		BaseURL: "http://localhost:11434",
		Model:   "llama2",
	})
	if err != nil {
		panic(err)
	}
	return chatModel
}
