package utils

import (
	"context"
	"github.com/sashabaranov/go-openai"
	"os"
)

// Reference to models which will be available in our platform
// const (
// 	O1Mini                  = "o1-mini"
// 	O1Preview               = "o1-preview"
// 	O1                      = "o1"
// 	O3                      = "o3"
// 	O320250416              = "o3-2025-04-16"
// 	O3Mini                  = "o3-mini"
// 	O4Mini                  = "o4-mini"
// 	GPT4o                   = "gpt-4o"
// 	GPT4oLatest             = "chatgpt-4o-latest"
// 	GPT4oMini               = "gpt-4o-mini"
// 	GPT4                    = "gpt-4"
// 	GPT4Dot1                = "gpt-4.1"
// )

func GetGptResponse(prompt string, model string) (string, error) {
	openaiClient := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	response, err := openaiClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: model,
			Messages: []openai.ChatCompletionMessage{
				{Role: "user", Content: prompt},
			},
		},
	)

	if err != nil {
		return "", err
	}

	if len(response.Choices) == 0 {
		return "", nil
	}

	return response.Choices[0].Message.Content, nil
}
