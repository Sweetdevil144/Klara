package utils

import (
	"context"
	"os"

	"github.com/openai/openai-go"
)

func GetGptResponse(prompt string) (string, error) {
	client, err := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	if err != nil {
		return "", err
	}

	response, err := client.CreateChatCompletion(
		context.Background(),
		openai.CreateChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{Role: openai.ChatMessageRoleUser, Content: prompt},
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