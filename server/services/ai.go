package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"server/models"
	"strings"
	"time"
)

type AIService struct {
	memoryService *MemoryService
	client        *http.Client
}

type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIRequest struct {
	Model       string          `json:"model"`
	Messages    []OpenAIMessage `json:"messages"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
	Temperature float64         `json:"temperature,omitempty"`
}

type OpenAIResponse struct {
	Choices []struct {
		Message OpenAIMessage `json:"message"`
	} `json:"choices"`
	Usage struct {
		TotalTokens int `json:"total_tokens"`
	} `json:"usage"`
}

type GeminiContent struct {
	Parts []struct {
		Text string `json:"text"`
	} `json:"parts"`
	Role string `json:"role,omitempty"`
}

type GeminiRequest struct {
	Contents []GeminiContent `json:"contents"`
}

type GeminiResponse struct {
	Candidates []struct {
		Content GeminiContent `json:"content"`
	} `json:"candidates"`
}

func NewAIService() *AIService {
	return &AIService{
		memoryService: NewMemoryService(),
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (ai *AIService) ChatWithAI(userID, clerkID, sessionID, message, model, apiKey string) (*models.ChatResponse, error) {
	memories, err := ai.memoryService.SearchUserMemories(clerkID, message, 5)
	if err != nil {
		fmt.Printf("Warning: Failed to search memories: %v\n", err)
		memories = []models.Memory{}
	}

	context := ai.buildContextFromMemories(memories)

	var response string
	switch model {
	case "openai":
		response, err = ai.callOpenAI(message, context, apiKey)
	case "gemini":
		response, err = ai.callGemini(message, context, apiKey)
	default:
		return nil, fmt.Errorf("unsupported model: %s", model)
	}

	if err != nil {
		return nil, fmt.Errorf("AI API call failed: %w", err)
	}

	go func() {
		ai.memoryService.AddChatMemory(clerkID, sessionID, message, "user")
		ai.memoryService.AddChatMemory(clerkID, sessionID, response, "assistant")
	}()

	return &models.ChatResponse{
		SessionID: sessionID,
		Message:   response,
		Role:      "assistant",
		Model:     model,
		Memories:  memories,
		CreatedAt: time.Now(),
	}, nil
}

func (ai *AIService) UpdateNoteWithAI(userID, clerkID, sessionID, noteID, noteContent, model, apiKey, customPrompt string) (string, error) {
	memories, err := ai.memoryService.SearchUserMemories(clerkID, sessionID, 10)
	if err != nil {
		return "", fmt.Errorf("failed to get chat history: %w", err)
	}

	prompt := ai.buildNoteUpdatePrompt(noteContent, memories, customPrompt)

	var updatedContent string
	switch model {
	case "openai":
		updatedContent, err = ai.callOpenAI(prompt, "", apiKey)
	case "gemini":
		updatedContent, err = ai.callGemini(prompt, "", apiKey)
	default:
		return "", fmt.Errorf("unsupported model: %s", model)
	}

	if err != nil {
		return "", fmt.Errorf("AI API call failed: %w", err)
	}

	return updatedContent, nil
}

func (ai *AIService) callOpenAI(message, context, apiKey string) (string, error) {
	url := "https://api.openai.com/v1/chat/completions"

	messages := []OpenAIMessage{}

	if context != "" {
		messages = append(messages, OpenAIMessage{
			Role:    "system",
			Content: fmt.Sprintf("Context from previous conversations:\n%s", context),
		})
	}

	messages = append(messages, OpenAIMessage{
		Role:    "user",
		Content: message,
	})

	request := OpenAIRequest{
		Model:       "gpt-3.5-turbo",
		Messages:    messages,
		MaxTokens:   1000,
		Temperature: 0.7,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Add("Content-Type", "application/json")

	resp, err := ai.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("OpenAI API error (%d): %s", resp.StatusCode, string(body))
	}

	var openAIResp OpenAIResponse
	if err := json.Unmarshal(body, &openAIResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(openAIResp.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI")
	}

	return openAIResp.Choices[0].Message.Content, nil
}

func (ai *AIService) callGemini(message, context, apiKey string) (string, error) {
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent?key=%s", apiKey)

	contents := []GeminiContent{}

	if context != "" {
		contents = append(contents, GeminiContent{
			Parts: []struct {
				Text string `json:"text"`
			}{
				{Text: fmt.Sprintf("Context from previous conversations:\n%s", context)},
			},
			Role: "user",
		})
	}

	contents = append(contents, GeminiContent{
		Parts: []struct {
			Text string `json:"text"`
		}{
			{Text: message},
		},
		Role: "user",
	})

	request := GeminiRequest{
		Contents: contents,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := ai.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("gemini api error (%d): %s", resp.StatusCode, string(body))
	}

	var geminiResp GeminiResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response from Gemini")
	}

	return geminiResp.Candidates[0].Content.Parts[0].Text, nil
}

func (ai *AIService) buildContextFromMemories(memories []models.Memory) string {
	if len(memories) == 0 {
		return ""
	}

	var contextParts []string
	for _, memory := range memories {
		contextParts = append(contextParts, memory.Memory)
	}

	return strings.Join(contextParts, "\n")
}

func (ai *AIService) buildNoteUpdatePrompt(currentNote string, memories []models.Memory, customPrompt string) string {
	basePrompt := `You are a focused note-taking assistant that updates notes based on conversation context.

STRICT GUIDELINES:
- ONLY update the note with information directly relevant to the note's topic
- DO NOT add tangential information, personal opinions, or unrelated content
- DO NOT include conversational elements, greetings, or meta-commentary
- DO NOT divert from the note's purpose unless explicitly requested by the user
- Preserve existing note structure and formatting exactly
- Maintain factual accuracy and professional tone

UPDATE RULES:
1. Review current note and conversation history
2. Identify ONLY facts, insights, or updates that directly relate to the note's topic
3. Integrate relevant information while preserving existing content structure
4. Only remove content if it's factually contradicted by new information
5. Use clear formatting (bullet points, headers) for readability
6. Keep content concise and focused on the note's purpose

CURRENT NOTE:
%s

CONVERSATION CONTEXT:
%s

TASK: Update the note by incorporating ONLY relevant information from the conversation. Return only the updated note content. Do not add any explanations, commentary, or off-topic content.`

	if customPrompt != "" {
		basePrompt = customPrompt + "\n\n" + basePrompt
	}

	conversationHistory := ai.buildContextFromMemories(memories)

	return fmt.Sprintf(basePrompt, currentNote, conversationHistory)
}

func (ai *AIService) ValidateAPIKey(model, apiKey string) error {
	testMessage := "Hello, this is a test message."

	switch model {
	case "openai":
		_, err := ai.callOpenAI(testMessage, "", apiKey)
		return err
	default:
		_, err := ai.callGemini(testMessage, "", apiKey)
		return err
	}
}
