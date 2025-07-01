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

func (ai *AIService) ChatWithAI(userID, clerkID, sessionID, message, modelID, provider, apiKey string) (*models.ChatResponse, error) {
	memories, err := ai.memoryService.SearchUserMemories(clerkID, message, 5)
	if err != nil {
		fmt.Printf("Warning: Failed to search memories: %v\n", err)
		memories = []models.Memory{}
	}

	context := ai.buildContextFromMemories(memories)

	systemPrompt := `You are a direct, professional AI assistant. Follow these rules:
- NEVER use conversational phrases like "Here's", "Okay", "I understand", "Sure", etc.
- NEVER include meta-commentary about what you're doing
- Provide direct, actionable responses
- Be concise, factual, and focused
- If providing information, present it clearly without unnecessary preamble
- If answering questions, give direct answers without conversational padding`

	if context != "" {
		context = systemPrompt + "\n\nPrevious conversation context:\n" + context
	} else {
		context = systemPrompt
	}

	var response string
	switch provider {
	case "openai":
		response, err = ai.callOpenAI(message, context, modelID, apiKey)
	case "gemini":
		response, err = ai.callGemini(message, context, modelID, apiKey)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
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
		Model:     provider,
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
		updatedContent, err = ai.callOpenAI(prompt, "", "gpt-3.5-turbo", apiKey)
	case "gemini":
		updatedContent, err = ai.callGemini(prompt, "", "gemini-1.5-flash", apiKey)
	default:
		return "", fmt.Errorf("unsupported model: %s", model)
	}

	if err != nil {
		return "", fmt.Errorf("AI API call failed: %w", err)
	}

	return updatedContent, nil
}

func (ai *AIService) callOpenAI(message, context, modelID, apiKey string) (string, error) {
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
		Model:       modelID,
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

func (ai *AIService) callGemini(message, context, modelID, apiKey string) (string, error) {
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", modelID, apiKey)

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
	basePrompt := `You are a precise note-updating assistant. Your ONLY job is to enhance the existing note content.

CRITICAL INSTRUCTIONS:
- NEVER include conversational text like "Here's the updated note", "I understand", "Okay", etc.
- NEVER replace existing content unless explicitly contradicted by new information
- ALWAYS preserve the existing structure, formatting, and all current content
- ADD new information by extending lists, adding sections, or appending relevant details
- If the user requests a 5th item and there are already 4 items, ADD the 5th item to the existing 4 items
- If adding to a list, continue the same numbering/bullet format
- If adding new sections, maintain consistent formatting
- Return ONLY the enhanced note content with no explanations or commentary

CONTENT PRESERVATION RULES:
1. Keep ALL existing content intact
2. Add new information in the appropriate location
3. Maintain original formatting (bullets, numbers, headers, etc.)
4. If user asks for "add X", append X to existing content
5. If user asks for "5th item", add it after existing items 1-4
6. Only remove content if new information directly contradicts old information

CURRENT NOTE CONTENT:
%s

CONVERSATION CONTEXT:
%s

ENHANCEMENT REQUEST: Based on the conversation, enhance the note by adding relevant information while preserving all existing content. Return only the enhanced note content with no additional text.`

	if customPrompt != "" {
		basePrompt = fmt.Sprintf(`CUSTOM INSTRUCTION: %s

%s`, customPrompt, basePrompt)
	}

	conversationHistory := ai.buildContextFromMemories(memories)

	return fmt.Sprintf(basePrompt, currentNote, conversationHistory)
}

func (ai *AIService) ValidateAPIKey(model, apiKey string) error {
	testMessage := "Hello, this is a test message."

	switch model {
	case "openai":
		_, err := ai.callOpenAI(testMessage, "", model, apiKey)
		return err
	default:
		_, err := ai.callGemini(testMessage, "", model, apiKey)
		return err
	}
}
