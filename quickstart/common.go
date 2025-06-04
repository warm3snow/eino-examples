package quickstart

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino/schema"
)

func CreateOllamaChatModel(ctx context.Context) *ollama.ChatModel {
	chatModel, err := ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
		BaseURL: "http://localhost:11434",
		Model:   "qwen3:1.7b", // this model supports tools
	})
	if err != nil {
		panic(err)
	}
	return chatModel
}

func ReportStream(resp *schema.StreamReader[*schema.Message]) {
	defer resp.Close()

	for {
		msg, err := resp.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error reading stream: %v", err)
			break
		}
		fmt.Printf("%s", msg.Content)
	}
	fmt.Println()
}
